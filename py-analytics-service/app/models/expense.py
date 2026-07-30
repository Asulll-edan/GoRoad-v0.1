from pydantic import BaseModel
from typing import Optional


class ExpenseSummaryRequest(BaseModel):
    user_id: str
    room_id: Optional[str] = None
    period: str = "all"


class ExpenseSummaryResponse(BaseModel):
    summary_json: str
