import structlog
import json
from datetime import datetime, timedelta
from app.config import settings
from app.tasks.celery_app import celery_app

logger = structlog.get_logger()

@celery_app.task(bind=True, max_retries=3, default_retry_delay=60)
def service_reminder_check(self):
    import asyncio
    loop = asyncio.new_event_loop()
    asyncio.set_event_loop(loop)
    try:
        return loop.run_until_complete(_check_reminders())
    finally:
        loop.close()

async def _check_reminders():
    import asyncpg
    import redis.asyncio as aioredis

    conn = await asyncpg.connect(
        host=settings.db_host,
        port=settings.db_port,
        user=settings.db_user,
        password=settings.db_password,
        database=settings.db_name,
    )
    r = aioredis.from_url(settings.redis_url, decode_responses=True)
    try:
        today = datetime.utcnow().date()
        target_dates = [today + timedelta(days=7), today + timedelta(days=1)]

        rows = await conn.fetch("""
            SELECT sr.id, sr.user_id, sr.title, sr.due_date, sr.motor_id,
                   m.name AS motor_name, m.plate_number
            FROM service_reminders sr
            JOIN motors m ON m.id = sr.motor_id
            WHERE sr.due_date = ANY($1::date[])
              AND sr.status = 'active'
        """, target_dates)

        notified = 0
        for row in rows:
            days_left = (row["due_date"] - today).days
            notification = {
                "user_id": str(row["user_id"]),
                "title": f"Service Reminder: {row['title']}",
                "body": (f"Motor {row['motor_name']} ({row['plate_number']}) "
                        f"perlu service dalam {days_left} hari"),
                "data": {
                    "type": "service_reminder",
                    "reminder_id": str(row["id"]),
                    "motor_id": str(row["motor_id"]),
                    "days_left": days_left,
                },
            }
            await r.lpush(f"cache:notifications:queue:{row['user_id']}", json.dumps(notification))
            await r.publish("notification.send", json.dumps(notification))
            notified += 1

        logger.info("service_reminder_checked", notified=notified, total=len(rows))
        return {"notified": notified, "total_reminders": len(rows)}
    except Exception as e:
        logger.error("service_reminder_check_failed", error=str(e))
        raise
    finally:
        await conn.close()
        await r.aclose()
