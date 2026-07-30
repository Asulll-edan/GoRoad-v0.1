from pydantic import BaseModel
from typing import Optional


class FuelAnalyticsRequest(BaseModel):
    user_id: str
    motor_id: Optional[str] = None
    months: int = 6


class FuelAnalyticsResponse(BaseModel):
    analytics_json: str
