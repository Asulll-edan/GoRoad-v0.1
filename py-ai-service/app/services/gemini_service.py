import structlog
import google.generativeai as genai
from app.config import settings
from app.cache import redis_client, make_hash_key

logger = structlog.get_logger()

class GeminiService:
    def __init__(self):
        self.model = None
        if settings.gemini_api_key:
            genai.configure(api_key=settings.gemini_api_key)
            self.model = genai.GenerativeModel("gemini-1.5-flash")
        else:
            logger.warning("Gemini API key not configured, AI features disabled")

    async def chat_stream(self, user_id: str, message: str, context: list = None):
        if not self.model:
            yield {"content": "AI assistant sedang tidak tersedia. Silakan coba lagi nanti.", "is_final": True}
            return

        cache_key = f"cache:ai:chat:{make_hash_key(message, context)}"
        cached = await redis_client.get(cache_key)
        if cached:
            yield {"content": cached, "is_final": True, "cached": True}
            return

        try:
            chat = self.model.start_chat()
            if context:
                for ctx in context[-10:]:
                    chat.send_message(ctx["content"])

            response = chat.send_message(message, stream=True)
            full_response = ""
            for chunk in response:
                if chunk.text:
                    full_response += chunk.text
                    yield {"content": chunk.text, "is_final": False}

            await redis_client.set(cache_key, full_response, ttl=3600)
            yield {"content": "", "is_final": True}

        except Exception as e:
            logger.error("Gemini chat error", error=str(e))
            yield {"content": "Maaf, terjadi kesalahan saat memproses permintaan.", "is_final": True}

    async def check_rate_limit(self, user_id: str) -> bool:
        key = f"rate:ai:{user_id}"
        count = await redis_client.incr_with_expiry(key, 3600)
        return count <= 20
