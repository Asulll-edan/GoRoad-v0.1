from pydantic import BaseModel
from typing import List, Optional


class BadgeEvalRequest(BaseModel):
    user_id: str
    room_id: Optional[str] = None
    touring_data_json: str = "{}"


class BadgeEvalResponse(BaseModel):
    new_badge_codes: List[str] = []
    badges_json: str = "[]"
