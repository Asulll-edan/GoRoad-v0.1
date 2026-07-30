import grpc
from google.protobuf import empty_pb2
from app.proto.ai_service_pb2 import (
    ChatRequest, ChatResponse,
    ItineraryRequest, ItineraryResponse,
    CostEstimationRequest, CostEstimationResponse,
    RouteAdviceRequest, RouteAdviceResponse,
    PackingListRequest, PackingListResponse,
    SafetyRequest, SafetyResponse,
    POIRecommendRequest, POIRecommendResponse,
)


class AIServiceStub:
    def __init__(self, channel):
        self.chat_stream = channel.unary_stream(
            "/ai_service.AIService/ChatStream",
            request_serializer=ChatRequest.SerializeToString,
            response_deserializer=ChatResponse.FromString,
        )
        self.generate_itinerary = channel.unary_unary(
            "/ai_service.AIService/GenerateItinerary",
            request_serializer=ItineraryRequest.SerializeToString,
            response_deserializer=ItineraryResponse.FromString,
        )
        self.estimate_cost = channel.unary_unary(
            "/ai_service.AIService/EstimateCost",
            request_serializer=CostEstimationRequest.SerializeToString,
            response_deserializer=CostEstimationResponse.FromString,
        )
        self.advise_route = channel.unary_unary(
            "/ai_service.AIService/AdviseRoute",
            request_serializer=RouteAdviceRequest.SerializeToString,
            response_deserializer=RouteAdviceResponse.FromString,
        )
        self.generate_packing_list = channel.unary_unary(
            "/ai_service.AIService/GeneratePackingList",
            request_serializer=PackingListRequest.SerializeToString,
            response_deserializer=PackingListResponse.FromString,
        )
        self.advise_safety = channel.unary_unary(
            "/ai_service.AIService/AdviseSafety",
            request_serializer=SafetyRequest.SerializeToString,
            response_deserializer=SafetyResponse.FromString,
        )
        self.recommend_poi = channel.unary_unary(
            "/ai_service.AIService/RecommendPOI",
            request_serializer=POIRecommendRequest.SerializeToString,
            response_deserializer=POIRecommendResponse.FromString,
        )


class AIServiceServicer:
    def ChatStream(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details("Method not implemented!")
        raise NotImplementedError("Method not implemented!")

    def GenerateItinerary(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details("Method not implemented!")
        raise NotImplementedError("Method not implemented!")

    def EstimateCost(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details("Method not implemented!")
        raise NotImplementedError("Method not implemented!")

    def AdviseRoute(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details("Method not implemented!")
        raise NotImplementedError("Method not implemented!")

    def GeneratePackingList(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details("Method not implemented!")
        raise NotImplementedError("Method not implemented!")

    def AdviseSafety(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details("Method not implemented!")
        raise NotImplementedError("Method not implemented!")

    def RecommendPOI(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details("Method not implemented!")
        raise NotImplementedError("Method not implemented!")


def add_AIServiceServicer_to_server(servicer, server):
    rpc_method_handlers = {
        "ChatStream": grpc.unary_stream_rpc_method_handler(
            servicer.ChatStream,
            request_deserializer=ChatRequest.FromString,
            response_serializer=ChatResponse.SerializeToString,
        ),
        "GenerateItinerary": grpc.unary_unary_rpc_method_handler(
            servicer.GenerateItinerary,
            request_deserializer=ItineraryRequest.FromString,
            response_serializer=ItineraryResponse.SerializeToString,
        ),
        "EstimateCost": grpc.unary_unary_rpc_method_handler(
            servicer.EstimateCost,
            request_deserializer=CostEstimationRequest.FromString,
            response_serializer=CostEstimationResponse.SerializeToString,
        ),
        "AdviseRoute": grpc.unary_unary_rpc_method_handler(
            servicer.AdviseRoute,
            request_deserializer=RouteAdviceRequest.FromString,
            response_serializer=RouteAdviceResponse.SerializeToString,
        ),
        "GeneratePackingList": grpc.unary_unary_rpc_method_handler(
            servicer.GeneratePackingList,
            request_deserializer=PackingListRequest.FromString,
            response_serializer=PackingListResponse.SerializeToString,
        ),
        "AdviseSafety": grpc.unary_unary_rpc_method_handler(
            servicer.AdviseSafety,
            request_deserializer=SafetyRequest.FromString,
            response_serializer=SafetyResponse.SerializeToString,
        ),
        "RecommendPOI": grpc.unary_unary_rpc_method_handler(
            servicer.RecommendPOI,
            request_deserializer=POIRecommendRequest.FromString,
            response_serializer=POIRecommendResponse.SerializeToString,
        ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
        "ai_service.AIService", rpc_method_handlers
    )
    server.add_generic_rpc_handlers((generic_handler,))
