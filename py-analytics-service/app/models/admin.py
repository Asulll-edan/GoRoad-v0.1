from pydantic import BaseModel


class AdminDashRequest(BaseModel):
    period: str = "7d"


class AdminDashResponse(BaseModel):
    dashboard_json: str
