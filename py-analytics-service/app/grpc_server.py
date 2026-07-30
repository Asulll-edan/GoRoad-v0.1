import grpc
from concurrent import futures
from app.config import settings

class AnalyticsGrpcServer:
    def __init__(self):
        self.server = grpc.aio.server(futures.ThreadPoolExecutor(max_workers=10))

    async def start(self):
        self.server.add_insecure_port(f"0.0.0.0:{settings.grpc_port}")
        await self.server.start()

    async def stop(self):
        await self.server.stop(grace=5)

grpc_server = AnalyticsGrpcServer()
