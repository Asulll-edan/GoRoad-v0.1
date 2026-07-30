import structlog
from app.cache import redis_client, make_hash_key
from app.services.gemini_service import GeminiService

logger = structlog.get_logger()

class ItineraryService:
    def __init__(self):
        self.gemini = GeminiService()

    async def generate(self, route_id: str, start_location: str, end_location: str,
                      duration_days: int, rider_count: int, motor_ids: list[str],
                      preferences: list[str]) -> dict:
        cache_key = f"cache:ai:itinerary:{make_hash_key(route_id, start_location, end_location, duration_days)}"
        cached = await redis_client.get_json(cache_key)
        if cached:
            return {**cached, "cached": True}

        prompt = f"""Buat itinerary touring motor detail:
Rute: {start_location} → {end_location}
Durasi: {duration_days} hari
Jumlah rider: {rider_count} orang
Preferensi: {', '.join(preferences) if preferences else 'Tidak ada preferensi khusus'}

Format response dalam JSON dengan struktur:
{{
  "days": [
    {{
      "day": 1,
      "date": "estimasi tanggal",
      "route": "lokasi start → lokasi finish",
      "distance_km": 0,
      "duration_hours": 0,
      "stops": [{{"name": "nama", "type": "rest/fuel/eat", "duration_minutes": 0, "notes": ""}}],
      "weather_note": "",
      "road_condition": "",
      "tips": ""
    }}
  ],
  "total_distance_km": 0,
  "total_duration_hours": 0,
  "recommendations": []
}}"""

        itinerary = await self._call_gemini(prompt)
        if itinerary:
            await redis_client.set_json(cache_key, itinerary, ttl=21600)
        return itinerary or {"error": "Gagal menghasilkan itinerary"}

    async def _call_gemini(self, prompt: str) -> dict:
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
            logger.error("Itinerary generation error", error=str(e))
            return None
