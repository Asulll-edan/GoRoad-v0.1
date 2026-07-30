from pydantic import BaseModel
from typing import List, Optional


class RouteAdviceRequest(BaseModel):
    origin: str
    destination: str
    waypoints: List[str] = []
    preferences: List[str] = []


class RouteAdviceResponse(BaseModel):
    advice_json: str
