import structlog
from fastapi import FastAPI
from prometheus_fastapi_instrumentator import Instrumentator
from app.config import settings
from app.cache import redis_client

logger = structlog.get_logger()

app = FastAPI(title="Go Road Analytics Service", version="3.0.0")

Instrumentator().instrument(app).expose(app)

@app.on_event("startup")
async def startup():
    await redis_client.connect()
    logger.info("Analytics Service started", env=settings.app_env)

@app.on_event("shutdown")
async def shutdown():
    await redis_client.disconnect()

@app.get("/health")
async def health():
    return {"status": "ok", "service": "analytics-service"}
