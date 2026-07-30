import asyncio
import json
import structlog
from app.config import settings
from app.tasks.feed_aggregator import aggregate_touring_feed
from app.tasks.smart_detection import detect_anomalies

logger = structlog.get_logger()

class NATSSubscriber:
    def __init__(self):
        self.nc = None
        self.js = None
        self._running = False

    async def connect(self):
        try:
            import nats
            self.nc = await nats.connect(settings.nats_url)
            self.js = self.nc.jetstream()
            self._running = True
            logger.info("nats_subscriber_connected", url=settings.nats_url)
        except Exception as e:
            logger.warning("nats_subscriber_connect_failed", error=str(e))
            self._running = False

    async def subscribe_all(self):
        if not self._running:
            logger.warning("nats_not_connected_skipping_subscriptions")
            return

        async def handle_touring_completed(msg):
            try:
                data = json.loads(msg.data.decode())
                room_id = data.get("room_id")
                touring_data = data.get("touring_data", {})
                if room_id:
                    aggregate_touring_feed.delay(room_id, touring_data)
                await msg.ack()
            except Exception as e:
                logger.error("handle_touring_completed_error", error=str(e))

        async def handle_location_batch(msg):
            try:
                data = json.loads(msg.data.decode())
                room_id = data.get("room_id")
                positions = data.get("positions", [])
                if room_id and positions:
                    detect_anomalies.delay(room_id, positions)
                await msg.ack()
            except Exception as e:
                logger.error("handle_location_batch_error", error=str(e))

        async def handle_cache_invalidation(msg):
            try:
                data = json.loads(msg.data.decode())
                keys = data.get("keys", [])
                if keys:
                    import redis.asyncio as aioredis
                    r = aioredis.from_url(settings.redis_url, decode_responses=True)
                    try:
                        await r.delete(*keys)
                    finally:
                        await r.aclose()
                await msg.ack()
            except Exception as e:
                logger.error("handle_cache_invalidation_error", error=str(e))

        try:
            await self.js.subscribe("touring.completed", cb=handle_touring_completed)
            await self.js.subscribe("room.*.location.batch", cb=handle_location_batch)
            await self.js.subscribe("cache.invalidate", cb=handle_cache_invalidation)
            logger.info("nats_subscriptions_active")
        except Exception as e:
            logger.warning("nats_subscribe_error", error=str(e))

    async def run(self):
        await self.connect()
        await self.subscribe_all()
        while self._running:
            await asyncio.sleep(1)

    async def close(self):
        self._running = False
        if self.nc:
            await self.nc.close()
