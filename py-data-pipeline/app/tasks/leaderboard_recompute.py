import structlog
import json
from app.config import settings
from app.tasks.celery_app import celery_app

logger = structlog.get_logger()

@celery_app.task(bind=True, max_retries=3, default_retry_delay=60)
def leaderboard_recompute(self):
    import asyncio
    loop = asyncio.new_event_loop()
    asyncio.set_event_loop(loop)
    try:
        return loop.run_until_complete(_recompute_all())
    finally:
        loop.close()

async def _recompute_all():
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
        rows = await conn.fetch("""
            SELECT
                u.id,
                COALESCE(SUM(r.distance_km), 0) AS total_distance,
                COUNT(DISTINCT r.id) AS total_tourings,
                COALESCE(SUM(r.duration_hours), 0) AS total_duration,
                (SELECT COUNT(*) FROM badges WHERE user_id = u.id) AS badge_count
            FROM users u
            LEFT JOIN routes r ON r.user_id = u.id AND r.status = 'completed'
            GROUP BY u.id
        """)

        monthly_cutoff = "2026-01-01"
        monthly_rows = await conn.fetch(f"""
            SELECT
                u.id,
                COALESCE(SUM(r.distance_km), 0) AS monthly_distance,
                COUNT(DISTINCT r.id) AS monthly_tourings
            FROM users u
            LEFT JOIN routes r ON r.user_id = u.id AND r.status = 'completed'
                AND r.updated_at >= $1::date
            GROUP BY u.id
        """, monthly_cutoff)

        all_time_key = "cache:leaderboard:all_time"
        monthly_key = "cache:leaderboard:monthly"

        pipe = r.pipeline()
        await pipe.delete(all_time_key, monthly_key)

        for row in rows:
            score = int(row["total_distance"] + row["total_tourings"] * 50 + row["badge_count"] * 100)
            await pipe.zadd(all_time_key, {str(row["id"]): score})

        for row in monthly_rows:
            score = int(row["monthly_distance"] * 2 + row["monthly_tourings"] * 100)
            await pipe.zadd(monthly_key, {str(row["id"]): score})

        await pipe.expire(all_time_key, 3600)
        await pipe.expire(monthly_key, 3600)
        await pipe.execute()

        logger.info("leaderboard_recomputed",
                    all_time=len(rows), monthly=len(monthly_rows))
        return {"all_time_count": len(rows), "monthly_count": len(monthly_rows)}
    except Exception as e:
        logger.error("leaderboard_recompute_failed", error=str(e))
        raise
    finally:
        await conn.close()
        await r.aclose()
