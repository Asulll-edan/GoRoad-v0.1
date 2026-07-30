from pydantic_settings import BaseSettings

class Settings(BaseSettings):
    app_env: str = "development"
    db_host: str = "localhost"
    db_port: int = 5432
    db_user: str = "postgres"
    db_password: str = "sandiandasalah"
    db_name: str = "goroad_dbvi1"
    redis_url: str = "redis://localhost:6379/2"
    nats_url: str = "nats://localhost:4222"
    openweather_api_key: str = ""
    
    model_config = {"env_prefix": ""}

    @property
    def database_url(self) -> str:
        return f"postgresql+asyncpg://{self.db_user}:{self.db_password}@{self.db_host}:{self.db_port}/{self.db_name}"

    @property
    def celery_broker_url(self) -> str:
        return self.redis_url

    @property
    def celery_result_backend(self) -> str:
        return self.redis_url

settings = Settings()