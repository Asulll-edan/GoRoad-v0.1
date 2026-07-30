import sentry_sdk
from sentry_sdk.integrations.fastapi import FastApiIntegration
from sentry_sdk.integrations.redis import RedisIntegration
from sentry_sdk.integrations.celery import CeleryIntegration
from sentry_sdk.integrations.grpc import GRPCIntegration


def init_sentry(dsn: str, environment: str = "production", traces_sample_rate: float = 0.1):
    sentry_sdk.init(
        dsn=dsn,
        environment=environment,
        traces_sample_rate=traces_sample_rate,
        profiles_sample_rate=0.05,
        integrations=[
            FastApiIntegration(),
            RedisIntegration(),
            CeleryIntegration(),
            GRPCIntegration(),
        ],
        send_default_pii=False,
    )
