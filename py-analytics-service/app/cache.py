import json
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

    async def zadd(self, key: str, members: dict[str, float]):
        await self.client.zadd(key, members)

    async def zrevrange(self, key: str, start: int = 0, stop: int = -1, withscores: bool = False) -> list:
        return await self.client.zrevrange(key, start, stop, withscores=withscores)

    async def zrevrank(self, key: str, member: str) -> Optional[int]:
        return await self.client.zrevrank(key, member)

    async def zscore(self, key: str, member: str) -> Optional[float]:
        return await self.client.zscore(key, member)

redis_client = RedisClient()
