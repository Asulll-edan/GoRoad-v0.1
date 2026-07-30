import json
import hashlib
from typing import Optional, Callable, Any, Awaitable
from functools import wraps

import redis.asyncio as redis
from pydantic import BaseModel


class RedisClient:
    def __init__(self, redis_url: str = "redis://localhost:6379/0"):
        self.redis_url = redis_url
        self._pool: Optional[redis.Redis] = None

    async def connect(self):
        if self._pool is None:
            self._pool = redis.from_url(
                self.redis_url,
                encoding="utf-8",
                decode_responses=True,
                max_connections=20,
            )

    async def disconnect(self):
        if self._pool:
            await self._pool.close()
            self._pool = None

    @property
    def client(self) -> redis.Redis:
        if self._pool is None:
            raise RuntimeError("Redis not connected. Call connect() first.")
        return self._pool

    # Basic operations
    async def get(self, key: str) -> Optional[str]:
        return await self.client.get(key)

    async def set(self, key: str, value: str, ttl: int = 300):
        await self.client.set(key, value, ex=ttl)

    async def delete(self, *keys: str):
        await self.client.delete(*keys)

    async def exists(self, key: str) -> bool:
        return await self.client.exists(key) > 0

    async def ttl(self, key: str) -> int:
        return await self.client.ttl(key)

    # JSON helper
    async def get_json(self, key: str, model: Optional[type[BaseModel]] = None) -> Optional[Any]:
        data = await self.get(key)
        if data is None:
            return None
        parsed = json.loads(data)
        if model:
            return model.model_validate(parsed)
        return parsed

    async def set_json(self, key: str, value: Any, ttl: int = 300):
        if isinstance(value, BaseModel):
            value = value.model_dump(mode="json")
        await self.set(key, json.dumps(value, default=str), ttl)

    # Hash operations (for rider positions)
    async def hset(self, key: str, field: str, value: Any):
        if isinstance(value, (dict, list)):
            value = json.dumps(value, default=str)
        await self.client.hset(key, field, value)

    async def hgetall(self, key: str) -> dict:
        return await self.client.hgetall(key)

    async def hdel(self, key: str, *fields: str):
        await self.client.hdel(key, *fields)

    # Set operations
    async def sadd(self, key: str, *members: str):
        await self.client.sadd(key, *members)

    async def srem(self, key: str, *members: str):
        await self.client.srem(key, *members)

    async def smembers(self, key: str) -> set:
        return await self.client.smembers(key)

    # Sorted Set (for leaderboard)
    async def zadd(self, key: str, members: dict[str, float]):
        await self.client.zadd(key, members)

    async def zrevrange(self, key: str, start: int = 0, stop: int = -1, withscores: bool = False) -> list:
        return await self.client.zrevrange(key, start, stop, withscores=withscores)

    async def zrevrank(self, key: str, member: str) -> Optional[int]:
        return await self.client.zrevrank(key, member)

    async def zscore(self, key: str, member: str) -> Optional[float]:
        return await self.client.zscore(key, member)

    # Rate limiting
    async def incr_with_expiry(self, key: str, ttl: int) -> int:
        count = await self.client.incr(key)
        if count == 1:
            await self.client.expire(key, ttl)
        return count

    # Distributed lock
    async def acquire_lock(self, key: str, ttl: int = 300) -> bool:
        return await self.client.set(f"lock:{key}", "1", nx=True, ex=ttl)

    async def release_lock(self, key: str):
        await self.client.delete(f"lock:{key}")

    # Pub/Sub
    async def publish(self, channel: str, message: Any):
        if isinstance(message, (dict, list, BaseModel)):
            message = json.dumps(message, default=str)
        await self.client.publish(channel, message)

    async def subscribe(self, channel: str):
        pubsub = self.client.pubsub()
        await pubsub.subscribe(channel)
        return pubsub


class cached:
    def __init__(self, key_pattern: str, ttl: int = 300, model: Optional[type[BaseModel]] = None):
        self.key_pattern = key_pattern
        self.ttl = ttl
        self.model = model

    def __call__(self, func: Callable[..., Awaitable[Any]]):
        @wraps(func)
        async def wrapper(*args, **kwargs):
            redis_client: RedisClient = kwargs.get("redis") or args[0].redis if hasattr(args[0], "redis") else None
            if redis_client is None:
                return await func(*args, **kwargs)

            key = self._build_key(*args, **kwargs)
            cached_result = await redis_client.get_json(key, self.model)
            if cached_result is not None:
                return cached_result

            result = await func(*args, **kwargs)
            await redis_client.set_json(key, result, self.ttl)
            return result

        return wrapper

    def _build_key(self, *args, **kwargs) -> str:
        import inspect
        sig = inspect.signature(self._original_func if hasattr(self, '_original_func') else (lambda: None))
        bound = sig.bind(*args, **kwargs)
        bound.apply_defaults()

        params = {}
        for name, value in bound.arguments.items():
            if name in ('self', 'cls', 'redis'):
                continue
            if hasattr(value, 'model_dump'):
                value = value.model_dump(mode="json")
            params[name] = value

        key = self.key_pattern.format(**params)
        return f"cache:{key}"


def make_hash_key(*args, **kwargs) -> str:
    raw = json.dumps({"args": args, "kwargs": kwargs}, sort_keys=True, default=str)
    return hashlib.sha256(raw.encode()).hexdigest()[:16]
