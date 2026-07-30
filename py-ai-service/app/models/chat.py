from pydantic import BaseModel
from typing import List, Optional


class ContextMessage(BaseModel):
    role: str
    content: str


class ChatRequest(BaseModel):
    user_id: str
    room_id: str
    message: str
    context: List[ContextMessage] = []


class ChatResponse(BaseModel):
    content: str
    is_final: bool = True
