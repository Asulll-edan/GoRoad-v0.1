from celery import Celery
from celery.schedules import crontab
from app.config import settings

celery_app = Celery(
    "goroad_pipeline",
    broker=settings.celery_broker_url,
    backend=settings.celery_result_backend,
)

celery_app.conf.update(
    task_serializer="json",
    accept_content=["json"],
    result_serializer="json",
    timezone="Asia/Jakarta",
    enable_utc=True,
    task_track_started=True,
    task_acks_late=True,
    worker_prefetch_multiplier=1,
)

celery_app.conf.beat_schedule = {
    "weather-sync": {
        "task": "app.tasks.weather_sync.weather_sync",
        "schedule": crontab(minute="*/30"),
    },
    "poi-cache-refresh": {
        "task": "app.tasks.poi_cache_refresh.poi_cache_refresh",
        "schedule": crontab(hour=3, minute=0),
    },
    "leaderboard-recompute": {
        "task": "app.tasks.leaderboard_recompute.leaderboard_recompute",
        "schedule": crontab(minute=0),
    },
    "service-reminder-check": {
        "task": "app.tasks.service_reminder_check.service_reminder_check",
        "schedule": crontab(hour=8, minute=0),
    },
    "data-retention-cleanup": {
        "task": "app.tasks.data_retention.data_retention_cleanup",
        "schedule": crontab(hour=2, minute=0),
    },
    "pre-touring-weather": {
        "task": "app.tasks.pre_touring_weather.pre_touring_weather",
        "schedule": crontab(minute="*/15"),
    },
}
