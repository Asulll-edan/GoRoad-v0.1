import pytest
from unittest.mock import AsyncMock, patch

from app.services.stats_service import StatsService


@pytest.fixture
def service():
    return StatsService()


@pytest.mark.asyncio
async def test_get_rider_stats_cached(service):
    with patch.object(service, "redis", AsyncMock()) as mock_redis:
        mock_redis.get_json = AsyncMock(return_value={"total_km": 1000})
        result = await service.get_rider_stats("user1")
        assert result["total_km"] == 1000


@pytest.mark.asyncio
async def test_get_rider_stats_new(service):
    with patch.object(service, "redis", AsyncMock()) as mock_redis:
        mock_redis.get_json = AsyncMock(return_value=None)
        mock_redis.set_json = AsyncMock()
        result = await service.get_rider_stats("user1")
        assert result is not None
