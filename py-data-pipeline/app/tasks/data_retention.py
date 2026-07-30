import structlog
from datetime import datetime, timedelta
from app.config import settings
from app.tasks.celery_app import celery_app

logger = structlog.get_logger()

@celery_app.task(bind=True, max_retries=2, default_retry_delay=120)
def data_retention_cleanup(self):
    import asyncio
    loop = asyncio.new_event_loop()
    asyncio.set_event_loop(loop)
    try:
        return loop.run_until_complete(_cleanup())
    finally:
        loop.close()

async def _cleanup():
    import asyncpg

    conn = await asyncpg.connect(
        host=settings.db_host,
        port=settings.db_port,
        user=settings.db_user,
        password=settings.db_password,
        database=settings.db_name,
    )
    try:
        now = datetime.utcnow()
        cutoff_chat = now - timedelta(days=365)
        cutoff_logs = now - timedelta(days=30)
        cutoff_analytics = now - timedelta(days=730)
        cutoff_soft_delete = now - timedelta(days=30)

        chat_deleted = await conn.execute("""
            DELETE FROM chat_messages WHERE created_at < $1
        """, cutoff_chat)

        analytics_deleted = await conn.execute("""
            DELETE FROM analytics_raw_data WHERE created_at < $1
        """, cutoff_analytics)

        users_hard_deleted = await conn.execute("""
            DELETE FROM users WHERE deleted_at IS NOT NULL AND deleted_at < $1
        """, cutoff_soft_delete)

        location_deleted = await conn.execute("""
            SELECT drop_chunks(older_than => INTERVAL '90 days')
        """)

        logger.info("data_retention_cleanup_complete",
                    chat_deleted=chat_deleted,
                    analytics_deleted=analytics_deleted,
                    users_hard_deleted=users_hard_deleted)
        return {
            "chat_deleted": chat_deleted,
            "analytics_deleted": analytics_deleted,
            "users_hard_deleted": users_hard_deleted,
        }
    except Exception as e:
        logger.error("data_retention_cleanup_failed", error=str(e))
        raise
    finally:
        await conn.close()
