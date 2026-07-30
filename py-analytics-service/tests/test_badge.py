import pytest
from unittest.mock import AsyncMock, patch

from app.services.badge_service import BadgeService


@pytest.fixture
def service():
    return BadgeService()


@pytest.mark.asyncio
async def test_evaluate_first_tour(service):
    with patch.object(service, "_get_user_badges", AsyncMock(return_value=[])):
        result = await service.evaluate("user1", "room1", {"touring_count": 1})
        assert "first_tour" in result["new_badge_codes"]


@pytest.mark.asyncio
async def test_evaluate_road_warrior(service):
    with patch.object(service, "_get_user_badges", AsyncMock(return_value=[])):
        result = await service.evaluate("user1", "room1", {"touring_count": 10})
        assert "road_warrior" in result["new_badge_codes"]


@pytest.mark.asyncio
async def test_evaluate_duplicate_badge(service):
    with patch.object(service, "_get_user_badges", AsyncMock(return_value=[{"code": "first_tour"}])):
        result = await service.evaluate("user1", "room1", {"touring_count": 1})
        assert "first_tour" not in result["new_badge_codes"]


@pytest.mark.asyncio
async def test_get_badge_progress(service):
    result = await service.get_badge_progress("user1")
    assert len(result) > 0
    assert all("badge_code" in b for b in result)
