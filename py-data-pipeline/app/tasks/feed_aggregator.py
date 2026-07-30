import structlog
import json
from app.config import settings
from app.tasks.celery_app import celery_app

logger = structlog.get_logger()

@celery_app.task(bind=True, max_retries=3, default_retry_delay=30)
def aggregate_touring_feed(self, room_id: str, touring_data: dict):
    import asyncio
    loop = asyncio.new_event_loop()
    asyncio.set_event_loop(loop)
    try:
        return loop.run_until_complete(_aggregate(room_id, touring_data))
    finally:
        loop.close()

async def _aggregate(room_id: str, touring_data: dict):
    import asyncpg
    import redis.asyncio as aioredis

    conn = await asyncpg.connect(
        host=settings.db_host,
        port=settings.db_port,
        user=settings.db_user,
        password=settings.db_password,
        database=settings.db_name,
    )
    r = aioredis.from_url(settings.redis_url, decode_responses=True)
    try:
        route_data = touring_data.get("route", {})
        participants = touring_data.get("participants", [])

        stats_snapshot = {
            "distance_km": route_data.get("distance_km", 0),
            "duration_hours": route_data.get("duration_hours", 0),
            "avg_speed_kmh": route_data.get("avg_speed_kmh", 0),
            "elevation_gain_m": route_data.get("elevation_gain_m", 0),
            "participant_count": len(participants),
            "start_time": touring_data.get("start_time"),
            "end_time": touring_data.get("end_time"),
        }

        route_rows = await conn.fetch("""
            SELECT id, name, distance_km, duration_hours, polyline
            FROM routes WHERE id = $1
        """, route_data.get("route_id"))

        if route_rows:
            rw = route_rows[0]
            route_snapshot = {
                "route_id": str(rw["id"]),
                "name": rw["name"],
                "distance_km": rw["distance_km"],
                "duration_hours": rw["duration_hours"],
                "simplified_polyline": _simplify_polyline(rw.get("polyline", "")),
            }

            feed_post = {
                "room_id": room_id,
                "type": "touring_completed",
                "stats_snapshot": stats_snapshot,
                "route_snapshot": route_snapshot,
                "participants": [str(p.get("user_id", "")) for p in participants],
                "timestamp": touring_data.get("end_time"),
            }

            feed_key = f"cache:feed:room:{room_id}"
            await r.lpush(feed_key, json.dumps(feed_post))
            await r.ltrim(feed_key, 0, 49)
            await r.expire(feed_key, 86400)

            await r.publish("feed.new_post", json.dumps(feed_post))

        logger.info("touring_feed_aggregated", room_id=room_id, participants=len(participants))
        return {"aggregated": True, "participants": len(participants)}
    except Exception as e:
        logger.error("feed_aggregation_failed", room_id=room_id, error=str(e))
        raise
    finally:
        await conn.close()
        await r.aclose()

def _simplify_polyline(polyline: str, max_points: int = 50) -> str:
    if not polyline:
        return ""
    try:
        import polyline as pl
        coords = pl.decode(polyline)
        if len(coords) <= max_points:
            return polyline
        step = max(1, len(coords) // max_points)
        simplified = [coords[i] for i in range(0, len(coords), step)]
        return pl.encode(simplified)
    except ImportError:
        return polyline
    except Exception:
        return polyline
