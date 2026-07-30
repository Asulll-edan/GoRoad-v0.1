import pytest
from unittest.mock import AsyncMock, patch

from app.services.gemini_service import GeminiService


@pytest.fixture
def service():
    with patch("app.services.gemini_service.settings") as mock_settings:
        mock_settings.gemini_api_key = "test-key"
        svc = GeminiService()
        svc.model = AsyncMock()
        return svc


@pytest.mark.asyncio
async def test_chat_stream_no_key():
    with patch("app.services.gemini_service.settings") as mock:
        mock.gemini_api_key = None
        svc = GeminiService()
        results = [r async for r in svc.chat_stream("user1", "hello")]
        assert len(results) == 1
        assert results[0]["is_final"] is True


@pytest.mark.asyncio
async def test_rate_limit():
    with patch("app.services.gemini_service.redis_client") as mock_redis:
        mock_redis.incr_with_expiry = AsyncMock(return_value=5)
        svc = GeminiService()
        result = await svc.check_rate_limit("user1")
        assert result is True
        mock_redis.incr_with_expiry.assert_called_once()


@pytest.mark.asyncio
async def test_rate_limit_exceeded():
    with patch("app.services.gemini_service.redis_client") as mock_redis:
        mock_redis.incr_with_expiry = AsyncMock(return_value=21)
        svc = GeminiService()
        result = await svc.check_rate_limit("user1")
        assert result is False
