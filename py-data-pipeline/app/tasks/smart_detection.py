import structlog
import numpy as np
from app.config import settings
from app.tasks.celery_app import celery_app

logger = structlog.get_logger()

STOPPED_THRESHOLD_SEC = 300
OFF_ROUTE_THRESHOLD_M = 500
SPEED_LIMIT_KMH = 130
BATTERY_LOW_PCT = 15
OFFLINE_THRESHOLD_SEC = 120
STRAGGLER_DISTANCE_KM = 1.5

@celery_app.task(bind=True, max_retries=2, default_retry_delay=30)
def detect_anomalies(self, room_id: str, rider_positions: list):
    import asyncio
    loop = asyncio.new_event_loop()
    asyncio.set_event_loop(loop)
    try:
        return loop.run_until_complete(_run_detection(room_id, rider_positions))
    finally:
        loop.close()

async def _run_detection(room_id: str, rider_positions: list):
    import redis.asyncio as aioredis
    r = aioredis.from_url(settings.redis_url, decode_responses=True)
    try:
        now = __import__("time").time()
        positions = np.array([
            [p["lat"], p["lon"], p.get("speed", 0), p.get("battery", 100), p.get("ts", now)]
            for p in rider_positions
        ])
        if len(positions) < 2:
            return {"anomalies": []}

        centroid_lat = np.mean(positions[:, 0])
        centroid_lon = np.mean(positions[:, 1])
        alerts = []

        for i, p in enumerate(rider_positions):
            user_id = p.get("user_id")
            if not user_id:
                continue

            cooldown_key = f"cache:detect:cooldown:{room_id}:{user_id}"
            cooldown_val = await r.get(cooldown_key)
            if cooldown_val:
                continue

            lat, lon = p["lat"], p["lon"]
            speed = p.get("speed", 0)
            battery = p.get("battery", 100)
            ts = p.get("ts", now)
            lat_diff = lat - centroid_lat
            lon_diff = lon - centroid_lon
            dist_km = np.sqrt(lat_diff**2 + lon_diff**2) * 111

            # straggler detection
            if dist_km > STRAGGLER_DISTANCE_KM:
                alerts.append({
                    "type": "straggler",
                    "user_id": user_id,
                    "room_id": room_id,
                    "distance_km": round(dist_km, 2),
                })
                await r.setex(cooldown_key, 120, "1")

            # speed limit
            if speed > SPEED_LIMIT_KMH:
                alerts.append({
                    "type": "speed_limit",
                    "user_id": user_id,
                    "room_id": room_id,
                    "speed": speed,
                })
                await r.setex(cooldown_key, 60, "1")

            # battery low
            if battery < BATTERY_LOW_PCT:
                alerts.append({
                    "type": "battery_low",
                    "user_id": user_id,
                    "room_id": room_id,
                    "battery": battery,
                })
                await r.setex(cooldown_key, 300, "1")

            # stopped too long
            age = now - ts
            if age > STOPPED_THRESHOLD_SEC and speed < 1:
                alerts.append({
                    "type": "stopped_long",
                    "user_id": user_id,
                    "room_id": room_id,
                    "duration_sec": int(age),
                })
                await r.setex(cooldown_key, 300, "1")

        if alerts:
            import json
            for alert in alerts:
                await r.publish("detection.alert", json.dumps(alert))

        return {"alerts": alerts}
    except Exception as e:
        logger.error("anomaly_detection_failed", room_id=room_id, error=str(e))
        return {"alerts": []}
    finally:
        await r.aclose()
