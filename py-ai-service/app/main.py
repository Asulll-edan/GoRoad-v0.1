import structlog
from fastapi import FastAPI
from prometheus_fastapi_instrumentator import Instrumentator

from app.config import settings
from app.cache import RedisClient

logger = structlog.get_logger()

app = FastAPI(title="Go Road AI Service", version="3.0.0")
redis_client = RedisClient()

Instrumentator().instrument(app).expose(app)

@app.on_event("startup")
async def startup():
    await redis_client.connect()
    logger.info("AI Service started", env=settings.app_env)

@app.on_event("shutdown")
async def shutdown():
    await redis_client.disconnect()

@app.get("/health")
async def health():
    return {"status": "ok", "service": "ai-service"}

@app.get("/metrics")
async def metrics():
    return {"message": "metrics endpoint"}
