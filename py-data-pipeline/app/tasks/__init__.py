from app.tasks import weather_sync
from app.tasks import poi_cache_refresh
from app.tasks import smart_detection
from app.tasks import leaderboard_recompute
from app.tasks import service_reminder_check
from app.tasks import data_retention
from app.tasks import feed_aggregator
from app.tasks import pre_touring_weather

__all__ = [
    "weather_sync",
    "poi_cache_refresh",
    "smart_detection",
    "leaderboard_recompute",
    "service_reminder_check",
    "data_retention",
    "feed_aggregator",
    "pre_touring_weather",
]
