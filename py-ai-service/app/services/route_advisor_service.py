import structlog
from app.services.gemini_service import GeminiService

logger = structlog.get_logger()

class RouteAdvisorService:
    def __init__(self):
        self.gemini = GeminiService()

    async def advise(self, origin: str, destination: str,
                    waypoints: list[str], preferences: list[str]) -> dict:
        prompt = f"""Berikan saran rute touring motor:
Dari: {origin}
Ke: {destination}
Waypoints: {', '.join(waypoints) if waypoints else 'Tidak ada'}
Preferensi: {', '.join(preferences) if preferences else 'Tidak ada preferensi khusus'}

Format JSON:
{{
  "alternatives": [
    {{
      "name": "Nama rute",
      "description": "Deskripsi singkat",
      "distance_km": 0,
      "estimated_duration": "X jam",
      "pros": ["kelebihan"],
      "cons": ["kekurangan"],
      "road_condition": "Kondisi jalan",
      "scenery": "Pemandangan"
    }}
  ],
  "recommended_index": 0,
  "tips": ["tips keamanan"]
}}"""

        try:
            import json
            import google.generativeai as genai
            from app.config import settings
            genai.configure(api_key=settings.gemini_api_key)
            model = genai.GenerativeModel("gemini-1.5-flash")
            response = model.generate_content(prompt)
            text = response.text.strip()
            if text.startswith("```json"):
                text = text[7:]
            if text.endswith("```"):
                text = text[:-3]
            return json.loads(text.strip())
        except Exception as e:
            logger.error("Route advice error", error=str(e))
            return {"error": "Gagal memberikan saran rute"}
