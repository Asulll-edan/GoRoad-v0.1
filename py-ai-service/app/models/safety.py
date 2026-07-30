from pydantic import BaseModel


class SafetyRequest(BaseModel):
    route_id: str
    weather_condition: str = "clear"
    rider_count: int = 1
    skill_level: str = "intermediate"


class SafetyResponse(BaseModel):
    advice_json: str
