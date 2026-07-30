import pytest
from unittest.mock import AsyncMock, patch

from app.services.itinerary_service import ItineraryService


@pytest.fixture
def service():
    return ItineraryService()


@pytest.mark.asyncio
async def test_generate_itinerary_cached(service):
    with patch.object(service, "redis", AsyncMock()) as mock_redis:
        mock_redis.get_json = AsyncMock(return_value={"cached": True})
        result = await service.generate_itinerary("route1", 3, ["motor1"])
        assert result["cached"] is True


@pytest.mark.asyncio
async def test_generate_itinerary_new(service):
    with patch.object(service, "redis", AsyncMock()) as mock_redis:
        mock_redis.get_json = AsyncMock(return_value=None)
        mock_redis.set_json = AsyncMock()
        result = await service.generate_itinerary("route1", 3, ["motor1"])
        assert result is not None
