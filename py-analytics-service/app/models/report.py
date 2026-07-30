from pydantic import BaseModel
from typing import Optional


class ReportRequest(BaseModel):
    type: str
    params_json: str = "{}"
    format: str = "json"


class ReportResponse(BaseModel):
    report_url: Optional[str] = None
    report_data_json: str = "{}"
