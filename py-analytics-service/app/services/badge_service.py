import structlog
from app.cache import redis_client

logger = structlog.get_logger()

BADGE_CRITERIA = {
    "first_tour": {"type": "touring_count", "value": 1},
    "road_warrior": {"type": "touring_count", "value": 10},
    "century_rider": {"type": "single_distance_km", "value": 100},
    "iron_butt": {"type": "single_distance_km", "value": 500},
    "thousand_milestone": {"type": "total_distance_km", "value": 1000},
    "eagle_eye": {"type": "single_elevation_m", "value": 1000},
    "speed_demon": {"type": "avg_speed_kmh", "value": 80},
    "night_owl": {"type": "night_touring", "value": True},
    "rain_rider": {"type": "rain_touring", "value": True},
    "convoy_leader": {"type": "lead_count", "value": 5},
    "social_butterfly": {"type": "unique_rooms", "value": 20},
    "perfect_attendance": {"type": "perfect_attendance", "value": True},
}

class BadgeService:
    async def evaluate(self, user_id: str, room_id: str, touring_data: dict) -> dict:
        new_badges = []
        user_badges = await self._get_user_badges(user_id)
        owned_codes = {b["code"] for b in user_badges}

        for code, criteria in BADGE_CRITERIA.items():
            if code in owned_codes:
                continue
            if self._check_criteria(criteria, touring_data):
                new_badges.append(code)

        return {
            "new_badge_codes": new_badges,
            "badges_json": user_badges,
        }

    def _check_criteria(self, criteria: dict, data: dict) -> bool:
        value_type = criteria["type"]
        required_value = criteria["value"]
        actual_value = data.get(value_type, 0)
        return actual_value >= required_value if isinstance(actual_value, (int, float)) else actual_value == required_value

    async def _get_user_badges(self, user_id: str) -> list:
        cache_key = f"cache:user:{user_id}:badges"
        cached = await redis_client.get_json(cache_key)
        return cached or []

    async def get_badge_progress(self, user_id: str) -> list:
        progress = []
        for code, criteria in BADGE_CRITERIA.items():
            progress.append({
                "badge_code": code,
                "badge_name": code.replace("_", " ").title(),
                "progress": 0,
                "target": criteria["value"],
                "is_awarded": False,
            })
        return progress
