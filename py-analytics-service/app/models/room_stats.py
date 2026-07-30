from pydantic import BaseModel


class RoomStatsRequest(BaseModel):
    room_id: str


class RoomStatsResponse(BaseModel):
    stats_json: str
