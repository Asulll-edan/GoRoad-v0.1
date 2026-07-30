import structlog
from app.cache import redis_client, make_hash_key
from app.services.gemini_service import GeminiService

logger = structlog.get_logger()

FUEL_PRICES = {
    "pertalite": 10000,
    "pertamax": 12950,
    "pertamax_turbo": 14000,
    "solar": 6800,
    "dexlite": 13100,
    "pertamina_dex": 14050,
}

class CostEstimationService:
    def __init__(self):
        self.gemini = GeminiService()

    async def estimate(self, route_id: str, motor_ids: list[str],
                      rider_count: int, duration_days: int, fuel_type: str) -> dict:
        cache_key = f"cache:ai:cost:{make_hash_key(route_id, motor_ids, duration_days)}"
        cached = await redis_client.get_json(cache_key)
        if cached:
            return {**cached, "cached": True}

        price_per_liter = FUEL_PRICES.get(fuel_type, 12950)

        prompt = f"""Estimasi biaya touring motor:
Jumlah motor: {len(motor_ids)}
Jumlah rider: {rider_count}
Durasi: {duration_days} hari
Tipe BBM: {fuel_type} (Rp{price_per_liter}/L)

Berikan estimasi biaya dalam format JSON:
{{
  "fuel_cost": 0,
  "food_cost": 0,
  "accommodation_cost": 0,
  "toll_cost": 0,
  "parking_cost": 0,
  "total_cost": 0,
  "per_person_cost": 0,
  "breakdown": {{}}
}}"""

        estimate = await self._call_gemini(prompt)
        if estimate:
            await redis_client.set_json(cache_key, estimate, ttl=21600)
        return estimate or {"error": "Gagal menghitung estimasi biaya"}

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
            logger.error("Cost estimation error", error=str(e))
            return None
