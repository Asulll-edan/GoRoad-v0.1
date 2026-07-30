import json
import hashlib
from typing import Optional, Type, Any
import redis.asyncio as redis
from pydantic import BaseModel
from app.config import settings

class RedisClient:
    def __init__(self):
        self.url = settings.redis_url
        self._pool: Optional[redis.Redis] = None

    async def connect(self):
        self._pool = redis.from_url(self.url, encoding="utf-8", decode_responses=True, max_connections=20)

    async def disconnect(self):
        if self._pool:
            await self._pool.close()

    @property
    def client(self) -> redis.Redis:
        if self._pool is None:
            raise RuntimeError("Redis not connected")
        return self._pool

    async def get(self, key: str) -> Optional[str]:
        return await self.client.get(key)

    async def set(self, key: str, value: str, ttl: int = 300):
        await self.client.set(key, value, ex=ttl)

    async def get_json(self, key: str, model: Optional[Type[BaseModel]] = None) -> Optional[Any]:
        data = await self.get(key)
        if data is None:
            return None
        parsed = json.loads(data)
        return model.model_validate(parsed) if model else parsed

    async def set_json(self, key: str, value: Any, ttl: int = 300):
        if isinstance(value, BaseModel):
            value = value.model_dump(mode="json")
        await self.set(key, json.dumps(value, default=str), ttl)

    async def delete(self, *keys: str):
        await self.client.delete(*keys)

    async def incr_with_expiry(self, key: str, ttl: int) -> int:
        count = await self.client.incr(key)
        if count == 1:
            await self.client.expire(key, ttl)
        return count

redis_client = RedisClient()

def make_hash_key(*args, **kwargs) -> str:
    raw = json.dumps({"args": args, "kwargs": kwargs}, sort_keys=True, default=str)
    return hashlib.sha256(raw.encode()).hexdigest()[:16]
