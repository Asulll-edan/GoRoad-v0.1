import asyncio
import json
from typing import Any, Callable, Optional
import nats
from nats.js import JetStreamContext
from structlog import get_logger

logger = get_logger()


class NatsClient:
    def __init__(self, url: str = "nats://localhost:4222"):
        self.url = url
        self.nc: Optional[nats.NATS] = None
        self.js: Optional[JetStreamContext] = None

    async def connect(self):
        self.nc = await nats.connect(self.url)
        self.js = self.nc.jetstream()
        logger.info("nats_client.connected", url=self.url)

    async def close(self):
        if self.nc:
            await self.nc.drain()
            logger.info("nats_client.disconnected")

    async def publish(self, subject: str, data: dict):
        if not self.nc:
            raise RuntimeError("NATS not connected")
        payload = json.dumps(data, default=str).encode()
        await self.nc.publish(subject, payload)

    async def subscribe(
        self, subject: str, callback: Callable, queue: Optional[str] = None
    ):
        if not self.nc:
            raise RuntimeError("NATS not connected")
        sub = await self.nc.subscribe(
            subject, cb=lambda msg: self._handle(msg, callback), queue=queue
        )
        return sub

    async def js_publish(self, subject: str, data: dict, stream: str = "goroad"):
        if not self.js:
            raise RuntimeError("JetStream not connected")
        payload = json.dumps(data, default=str).encode()
        ack = await self.js.publish(subject, payload, stream=stream)
        return ack

    async def js_subscribe(
        self, subject: str, callback: Callable, stream: str = "goroad", durable: Optional[str] = None
    ):
        if not self.js:
            raise RuntimeError("JetStream not connected")
        sub = await self.js.subscribe(
            subject,
            cb=lambda msg: self._handle(msg, callback),
            stream=stream,
            durable=durable,
        )
        return sub

    async def _handle(self, msg, callback: Callable):
        try:
            data = json.loads(msg.data.decode())
            await callback(data)
            await msg.ack()
        except Exception as e:
            logger.error("nats_client.handle_error", error=str(e))
