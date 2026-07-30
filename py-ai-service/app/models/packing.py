from pydantic import BaseModel


class PackingListRequest(BaseModel):
    duration_days: int
    weather_condition: str = "clear"
    touring_type: str = "adventure"


class PackingListResponse(BaseModel):
    packing_list_json: str
