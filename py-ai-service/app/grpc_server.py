import asyncio
import grpc
from concurrent import futures
from app.config import settings

class AIGrpcServer:
    def __init__(self):
        self.server = grpc.aio.server(futures.ThreadPoolExecutor(max_workers=10))

    async def start(self):
        # gRPC service implementations will be added here
        self.server.add_insecure_port(f"0.0.0.0:{settings.grpc_port}")
        await self.server.start()
        print(f"gRPC server started on port {settings.grpc_port}")

    async def stop(self):
        await self.server.stop(grace=5)

grpc_server = AIGrpcServer()
