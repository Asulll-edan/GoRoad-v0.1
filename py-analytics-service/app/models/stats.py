from pydantic import BaseModel
from typing import Optional


class RiderStatsRequest(BaseModel):
    user_id: str
    force_refresh: bool = False


class RiderStatsResponse(BaseModel):
    stats_json: str
    cached: Optional[str] = None
