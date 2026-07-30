from pydantic import BaseModel
from typing import Optional


class LeaderboardRequest(BaseModel):
    period: str = "all"
    limit: int = 20
    cursor: Optional[str] = None


class LeaderboardResponse(BaseModel):
    leaderboard_json: str
    cursor: Optional[str] = None
    has_more: bool = False
