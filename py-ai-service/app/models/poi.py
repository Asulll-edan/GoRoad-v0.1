from pydantic import BaseModel
from typing import List, Optional


class POIRecommendRequest(BaseModel):
    lat: float
    lng: float
    radius_km: float = 10.0
    types: List[str] = []
    route_id: Optional[str] = None


class POIRecommendResponse(BaseModel):
    pois_json: str
    cached: Optional[str] = None
