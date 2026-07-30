from pydantic import BaseModel
from typing import List, Optional


class ItineraryRequest(BaseModel):
    route_id: str
    start_location: str
    end_location: str
    duration_days: int
    rider_count: int = 1
    motor_ids: List[str] = []
    preferences: List[str] = []


class ItineraryResponse(BaseModel):
    itinerary_json: str
    cached: Optional[str] = None
