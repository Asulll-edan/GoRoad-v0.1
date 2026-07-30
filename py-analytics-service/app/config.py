from pydantic_settings import BaseSettings

class Settings(BaseSettings):
    app_env: str = "development"
    db_host: str = "localhost"
    db_port: int = 5432
    db_user: str = "postgres"
    db_password: str = "sandiandasalah"
    db_name: str = "goroad_dbvi1"
    redis_url: str = "redis://localhost:6379/1"
    nats_url: str = "nats://localhost:4222"
    grpc_port: int = 50052

    model_config = {"env_prefix": ""}

    @property
    def database_url(self) -> str:
        return f"postgresql+asyncpg://{self.db_user}:{self.db_password}@{self.db_host}:{self.db_port}/{self.db_name}"

settings = Settings()
