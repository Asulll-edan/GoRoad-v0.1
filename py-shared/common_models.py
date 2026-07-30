from typing import Any, Dict, List, Optional
from pydantic import BaseModel
from datetime import datetime


class PaginatedRequest(BaseModel):
    cursor: Optional[str] = None
    limit: int = 20


class PaginatedResponse(BaseModel):
    data: List[Dict[str, Any]]
    cursor: Optional[str] = None
    has_more: bool = False


class CacheConfig(BaseModel):
    host: str = "localhost"
    port: int = 6379
    db: int = 0
    password: Optional[str] = None
    default_ttl: int = 300


class NatsConfig(BaseModel):
    url: str = "nats://localhost:4222"
    stream: str = "goroad"


class DBConfig(BaseModel):
    host: str = "localhost"
    port: int = 5432
    user: str = "postgres"
    password: str = ""
    database: str = "goroad_dbvi1"


class RiderPosition(BaseModel):
    user_id: str
    room_id: str
    lat: float
    lon: float
    speed: float = 0.0
    heading: float = 0.0
    altitude: float = 0.0
    accuracy: float = 0.0
    battery: int = 100
    timestamp: datetime = None


class NotificationPayload(BaseModel):
    user_id: str
    title: str
    body: str
    data: Optional[Dict[str, Any]] = None
    type: str = "general"
