from pydantic import BaseModel
from typing import List, Optional


class CostEstimationRequest(BaseModel):
    route_id: str
    motor_ids: List[str] = []
    rider_count: int = 1
    duration_days: int
    fuel_type: str = "pertalite"


class CostEstimationResponse(BaseModel):
    estimate_json: str
    cached: Optional[str] = None
