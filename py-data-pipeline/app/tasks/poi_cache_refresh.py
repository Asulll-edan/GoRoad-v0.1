import structlog
import httpx
from app.config import settings
from app.tasks.celery_app import celery_app

logger = structlog.get_logger()

OVERPASS_URL = "https://overpass-api.de/api/interpreter"

@celery_app.task(bind=True, max_retries=3, default_retry_delay=120)
def poi_cache_refresh(self):
    import asyncio
    loop = asyncio.new_event_loop()
    asyncio.set_event_loop(loop)
    try:
        return loop.run_until_complete(_refresh_all_regions())
    finally:
        loop.close()

async def _refresh_all_regions():
    import redis.asyncio as aioredis
    r = aioredis.from_url(settings.redis_url, decode_responses=True)
    try:
        active_key = "cache:poi:active_geohashes"
        geohashes = await r.smembers(active_key)
        if not geohashes:
            logger.info("no_active_geohashes_for_poi")
            return {"refreshed": 0}

        refreshed = 0
        async with httpx.AsyncClient(timeout=30) as client:
            for geohash in geohashes:
                lock_key = f"lock:poi:refresh:{geohash}"
                locked = await r.setnx(lock_key, "1")
                if not locked:
                    continue
                await r.expire(lock_key, 900)
                try:
                    bbox = _geohash_to_bbox(geohash)
                    overpass_query = f"""
                    [out:json];
                    (
                        node["amenity"]({bbox[0]},{bbox[1]},{bbox[2]},{bbox[3]});
                        node["tourism"]({bbox[0]},{bbox[1]},{bbox[2]},{bbox[3]});
                        node["shop"]({bbox[0]},{bbox[1]},{bbox[2]},{bbox[3]});
                        node["leisure"]({bbox[0]},{bbox[1]},{bbox[2]},{bbox[3]});
                        way["highway"]({bbox[0]},{bbox[1]},{bbox[2]},{bbox[3]});
                    );
                    out center;
                    """
                    resp = await client.post(OVERPASS_URL, data={"data": overpass_query})
                    if resp.status_code == 200:
                        data = resp.json()
                        elements = data.get("elements", [])
                        pois = []
                        for elem in elements:
                            poi = {
                                "id": elem.get("id"),
                                "type": elem.get("type"),
                                "lat": elem.get("lat") or elem.get("center", {}).get("lat"),
                                "lon": elem.get("lon") or elem.get("center", {}).get("lon"),
                                "tags": elem.get("tags", {}),
                            }
                            pois.append(poi)
                        cache_key = f"cache:poi:region:{geohash}"
                        await r.setex(cache_key, 86400, pois)
                        await r.setex(f"{cache_key}:updated", 86400, "1")
                        refreshed += 1
                    else:
                        logger.warning("overpass_api_error", geohash=geohash, status=resp.status_code)
                except Exception as e:
                    logger.error("poi_refresh_error", geohash=geohash, error=str(e))
                finally:
                    await r.delete(lock_key)
        return {"refreshed": refreshed}
    except Exception as e:
        logger.error("poi_cache_refresh_failed", error=str(e))
        raise
    finally:
        await r.aclose()

def _geohash_to_bbox(geohash: str, precision: int = 6):
    import geohash2
    lat, lon = geohash2.decode(geohash)
    lat_err, lon_err = geohash2.decode_exactly(geohash)[2:]
    return (lat - lat_err, lon - lon_err, lat + lat_err, lon + lon_err)
