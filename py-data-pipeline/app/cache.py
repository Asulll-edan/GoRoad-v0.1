import json
import pickle
from typing import Any, Optional
import redis.asyncio as aioredis
from app.config import settings

redis_client: Optional[aioredis.Redis] = None


async def init_cache():
    global redis_client
    redis_client = aioredis.Redis(
        host=settings.REDIS_HOST,
        port=settings.REDIS_PORT,
        password=settings.REDIS_PASSWORD or None,
        db=settings.REDIS_DB or 0,
        decode_responses=False,
    )
    await redis_client.ping()


async def close_cache():
    if redis_client:
        await redis_client.close()


async def get_cache(key: str) -> Optional[Any]:
    if not redis_client:
        return None
    data = await redis_client.get(key)
    if data is None:
        return None
    try:
        return pickle.loads(data)
    except (pickle.UnpicklingError, TypeError):
        return data


async def set_cache(key: str, value: Any, ttl: int = 300):
    if not redis_client:
        return
    data = pickle.dumps(value)
    await redis_client.setex(key, ttl, data)


async def delete_cache(key: str):
    if not redis_client:
        return
    await redis_client.delete(key)


async def cache_set_json(key: str, value: Any, ttl: int = 300):
    if not redis_client:
        return
    data = json.dumps(value, default=str)
    await redis_client.setex(key, ttl, data)


async def cache_get_json(key: str) -> Optional[Any]:
    if not redis_client:
        return None
    data = await redis_client.get(key)
    if data is None:
        return None
    try:
        return json.loads(data)
    except (json.JSONDecodeError, TypeError):
        return None


async def acquire_lock(lock_key: str, ttl: int = 10) -> bool:
    return bool(await redis_client.setnx(lock_key, "1")) if redis_client else False


async def release_lock(lock_key: str):
    if redis_client:
        await redis_client.delete(lock_key)


async def get_connection() -> aioredis.Redis:
    return redis_client
