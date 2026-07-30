from pydantic_settings import BaseSettings

class Settings(BaseSettings):
    app_env: str = "development"
    redis_url: str = "redis://localhost:6379/0"
    nats_url: str = "nats://localhost:4222"
    gemini_api_key: str = ""
    grpc_port: int = 50051
    http_port: int = 8000
    
    model_config = {"env_prefix": ""}

settings = Settings()
