import structlog
import httpx
from celery import shared_task
from app.config import settings
from app.tasks.celery_app import celery_app

logger = structlog.get_logger()

@celery_app.task(bind=True, max_retries=3, default_retry_delay=60)
def weather_sync(self):
    import asyncio
    loop = asyncio.new_event_loop()
    asyncio.set_event_loop(loop)
    try:
        return loop.run_until_complete(_fetch_and_cache_weather())
    finally:
        loop.close()

@celery_app.task(bind=True, max_retries=3, default_retry_delay=60)
def pre_touring_weather_notify(self):
    import asyncio
    loop = asyncio.new_event_loop()
    asyncio.set_event_loop(loop)
    try:
        return loop.run_until_complete(_notify_upcoming_touring())
    finally:
        loop.close()

async def _fetch_and_cache_weather():
    import redis.asyncio as aioredis
    r = aioredis.from_url(settings.redis_url, decode_responses=True)
    try:
        active_rooms_key = "cache:weather:active_regions"
        regions = await r.smembers(active_rooms_key)
        if not regions:
            logger.info("no_active_regions_for_weather")
            return {"synced": 0}

        synced = 0
        async with httpx.AsyncClient(timeout=15) as client:
            for region in regions:
                lock_key = f"lock:weather:sync:{region}"
                locked = await r.setnx(lock_key, "1")
                if not locked:
                    continue
                await r.expire(lock_key, 300)
                try:
                    lat, lon = region.split(",")
                    url = (
                        f"https://api.openweathermap.org/data/2.5/weather"
                        f"?lat={lat}&lon={lon}&appid={settings.openweather_api_key}&units=metric"
                    )
                    resp = await client.get(url)
                    if resp.status_code == 200:
                        data = resp.json()
                        cache_key = f"cache:weather:{region}"
                        await r.setex(cache_key, 1800, data)
                        synced += 1
                        if _is_extreme_weather(data):
                            await r.publish(
                                "weather.alert",
                                f'{{"region":"{region}","alert":{_extract_alert(data)}}}'
                            )
                    else:
                        logger.warning("weather_api_error", region=region, status=resp.status_code)
                except Exception as e:
                    logger.error("weather_sync_error", region=region, error=str(e))
                finally:
                    await r.delete(lock_key)
        return {"synced": synced}
    except Exception as e:
        logger.error("weather_sync_failed", error=str(e))
        raise
    finally:
        await r.aclose()

async def _notify_upcoming_touring():
    import redis.asyncio as aioredis
    r = aioredis.from_url(settings.redis_url, decode_responses=True)
    try:
        upcoming_key = "cache:touring:upcoming_24h"
        active_rooms = await r.smembers(upcoming_key)
        if not active_rooms:
            return {"notified": 0}

        notified = 0
        async with httpx.AsyncClient(timeout=15) as client:
            for room_id in active_rooms:
                room_data_raw = await r.get(f"cache:room:{room_id}:route")
                if not room_data_raw:
                    continue
                import json
                room_data = json.loads(room_data_raw)
                lat, lon = room_data.get("lat", 0), room_data.get("lon", 0)
                url = (
                    f"https://api.openweathermap.org/data/2.5/forecast"
                    f"?lat={lat}&lon={lon}&appid={settings.openweather_api_key}&units=metric"
                )
                resp = await client.get(url)
                if resp.status_code == 200:
                    forecast = resp.json()
                    await r.setex(f"cache:weather:forecast:{room_id}", 3600, json.dumps(forecast))
                    await r.publish(
                        "weather.forecast_ready",
                        json.dumps({"room_id": room_id})
                    )
                    notified += 1
        return {"notified": notified}
    except Exception as e:
        logger.error("pre_touring_weather_failed", error=str(e))
        raise
    finally:
        await r.aclose()

def _is_extreme_weather(data: dict) -> bool:
    main = data.get("main", {})
    wind = data.get("wind", {})
    weather = data.get("weather", [{}])[0]
    condition = weather.get("main", "").lower()
    temp = main.get("temp", 25)
    wind_speed = wind.get("speed", 0)
    return any([
        condition in ("thunderstorm", "tornado", "hurricane"),
        wind_speed > 20,
        temp > 40,
        temp < 5,
    ])

def _extract_alert(data: dict) -> str:
    import json
    weather = data.get("weather", [{}])[0]
    return json.dumps({
        "condition": weather.get("main", ""),
        "description": weather.get("description", ""),
        "temp": data.get("main", {}).get("temp"),
        "wind_speed": data.get("wind", {}).get("speed"),
    })
