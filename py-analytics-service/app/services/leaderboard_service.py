import structlog
from app.cache import redis_client

logger = structlog.get_logger()

class LeaderboardService:
    async def compute_leaderboard(self, period: str = "monthly", limit: int = 20, cursor: str = None) -> dict:
        cache_key = f"cache:leaderboard:{'monthly' if period == 'monthly' else 'alltime'}:1"

        if period == "monthly":
            import datetime
            now = datetime.datetime.now()
            cache_key = f"cache:leaderboard:monthly:{now.year}_{now.month:02d}"

        cached = await redis_client.get_json(cache_key)
        if cached:
            return cached

        data = {
            "entries": [
                {"rank": 1, "user_id": "", "username": "", "points": 0, "avatar_url": ""}
            ],
            "cursor": None,
            "has_more": False,
        }

        await redis_client.set_json(cache_key, data, ttl=3600)
        return data

    async def update_score(self, user_id: str, points: float):
        import datetime
        now = datetime.datetime.now()
        monthly_key = f"cache:leaderboard:monthly:{now.year}_{now.month:02d}"
        alltime_key = "cache:leaderboard:alltime"

        await redis_client.zadd(monthly_key, {user_id: points})
        await redis_client.zadd(alltime_key, {user_id: points})
