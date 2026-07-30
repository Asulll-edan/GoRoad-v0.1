from typing import Optional
import httpx
from app.config import settings
from app.cache import cache_get_json, cache_set_json


class SafetyService:
    def __init__(self):
        self.api_key = settings.GEMINI_API_KEY

    async def get_advice(
        self,
        route_id: str,
        weather_condition: str = "clear",
        rider_count: int = 1,
        skill_level: str = "intermediate",
    ) -> str:
        cache_key = f"cache:ai:safety:{route_id}:{weather_condition}:{rider_count}:{skill_level}"
        cached = await cache_get_json(cache_key)
        if cached:
            return cached["advice_json"]

        route_info = await self._get_route_info(route_id)
        prompt = self._build_prompt(route_info, weather_condition, rider_count, skill_level)
        advice = await self._call_gemini(prompt)

        await cache_set_json(cache_key, {"advice_json": advice}, ttl=3600)
        return advice

    def _build_prompt(self, route: Optional[dict], weather: str, riders: int, skill: str) -> str:
        route_info = f"Route: {route}" if route else "Unknown route"
        return (
            f"As a motor touring safety expert, provide safety advice for:\n"
            f"Route info: {route_info}\n"
            f"Weather: {weather}\n"
            f"Rider count: {riders}\n"
            f"Rider skill level: {skill}\n\n"
            f"Include:\n"
            f"1. Pre-ride safety checklist\n"
            f"2. Weather-related precautions\n"
            f"3. Group riding formation tips\n"
            f"4. Emergency procedures\n"
            f"5. Essential gear recommendations\n"
            f"Format as JSON with sections."
        )

    async def _get_route_info(self, route_id: str) -> Optional[dict]:
        return None

    async def _call_gemini(self, prompt: str) -> str:
        if not self.api_key:
            return '{"error": "AI service unavailable", "advice": []}'
        import google.generativeai as genai
        genai.configure(api_key=self.api_key)
        model = genai.GenerativeModel("gemini-2.0-flash")
        response = await model.generate_content_async(prompt)
        return response.text
