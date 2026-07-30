import grpc
from app.proto.analytics_service_pb2 import (
    RiderStatsRequest, RiderStatsResponse,
    LeaderboardRequest, LeaderboardResponse,
    BadgeEvalRequest, BadgeEvalResponse,
    ReportRequest, ReportResponse,
    RoomStatsRequest, RoomStatsResponse,
    FuelAnalyticsRequest, FuelAnalyticsResponse,
    ExpenseSummaryRequest, ExpenseSummaryResponse,
    AdminDashRequest, AdminDashResponse,
)


class AnalyticsServiceStub:
    def __init__(self, channel):
        self.compute_rider_stats = channel.unary_unary(
            "/analytics_service.AnalyticsService/ComputeRiderStats",
            request_serializer=RiderStatsRequest.SerializeToString,
            response_deserializer=RiderStatsResponse.FromString,
        )
        self.compute_leaderboard = channel.unary_unary(
            "/analytics_service.AnalyticsService/ComputeLeaderboard",
            request_serializer=LeaderboardRequest.SerializeToString,
            response_deserializer=LeaderboardResponse.FromString,
        )
        self.evaluate_badges = channel.unary_unary(
            "/analytics_service.AnalyticsService/EvaluateBadges",
            request_serializer=BadgeEvalRequest.SerializeToString,
            response_deserializer=BadgeEvalResponse.FromString,
        )
        self.generate_report = channel.unary_unary(
            "/analytics_service.AnalyticsService/GenerateReport",
            request_serializer=ReportRequest.SerializeToString,
            response_deserializer=ReportResponse.FromString,
        )
        self.compute_room_stats = channel.unary_unary(
            "/analytics_service.AnalyticsService/ComputeRoomStats",
            request_serializer=RoomStatsRequest.SerializeToString,
            response_deserializer=RoomStatsResponse.FromString,
        )
        self.compute_fuel_analytics = channel.unary_unary(
            "/analytics_service.AnalyticsService/ComputeFuelAnalytics",
            request_serializer=FuelAnalyticsRequest.SerializeToString,
            response_deserializer=FuelAnalyticsResponse.FromString,
        )
        self.compute_expense_summary = channel.unary_unary(
            "/analytics_service.AnalyticsService/ComputeExpenseSummary",
            request_serializer=ExpenseSummaryRequest.SerializeToString,
            response_deserializer=ExpenseSummaryResponse.FromString,
        )
        self.get_admin_dashboard = channel.unary_unary(
            "/analytics_service.AnalyticsService/GetAdminDashboard",
            request_serializer=AdminDashRequest.SerializeToString,
            response_deserializer=AdminDashResponse.FromString,
        )


class AnalyticsServiceServicer:
    def ComputeRiderStats(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details("Method not implemented!")
        raise NotImplementedError("Method not implemented!")

    def ComputeLeaderboard(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details("Method not implemented!")
        raise NotImplementedError("Method not implemented!")

    def EvaluateBadges(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details("Method not implemented!")
        raise NotImplementedError("Method not implemented!")

    def GenerateReport(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details("Method not implemented!")
        raise NotImplementedError("Method not implemented!")

    def ComputeRoomStats(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details("Method not implemented!")
        raise NotImplementedError("Method not implemented!")

    def ComputeFuelAnalytics(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details("Method not implemented!")
        raise NotImplementedError("Method not implemented!")

    def ComputeExpenseSummary(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details("Method not implemented!")
        raise NotImplementedError("Method not implemented!")

    def GetAdminDashboard(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details("Method not implemented!")
        raise NotImplementedError("Method not implemented!")


def add_AnalyticsServiceServicer_to_server(servicer, server):
    rpc_method_handlers = {
        "ComputeRiderStats": grpc.unary_unary_rpc_method_handler(
            servicer.ComputeRiderStats,
            request_deserializer=RiderStatsRequest.FromString,
            response_serializer=RiderStatsResponse.SerializeToString,
        ),
        "ComputeLeaderboard": grpc.unary_unary_rpc_method_handler(
            servicer.ComputeLeaderboard,
            request_deserializer=LeaderboardRequest.FromString,
            response_serializer=LeaderboardResponse.SerializeToString,
        ),
        "EvaluateBadges": grpc.unary_unary_rpc_method_handler(
            servicer.EvaluateBadges,
            request_deserializer=BadgeEvalRequest.FromString,
            response_serializer=BadgeEvalResponse.SerializeToString,
        ),
        "GenerateReport": grpc.unary_unary_rpc_method_handler(
            servicer.GenerateReport,
            request_deserializer=ReportRequest.FromString,
            response_serializer=ReportResponse.SerializeToString,
        ),
        "ComputeRoomStats": grpc.unary_unary_rpc_method_handler(
            servicer.ComputeRoomStats,
            request_deserializer=RoomStatsRequest.FromString,
            response_serializer=RoomStatsResponse.SerializeToString,
        ),
        "ComputeFuelAnalytics": grpc.unary_unary_rpc_method_handler(
            servicer.ComputeFuelAnalytics,
            request_deserializer=FuelAnalyticsRequest.FromString,
            response_serializer=FuelAnalyticsResponse.SerializeToString,
        ),
        "ComputeExpenseSummary": grpc.unary_unary_rpc_method_handler(
            servicer.ComputeExpenseSummary,
            request_deserializer=ExpenseSummaryRequest.FromString,
            response_serializer=ExpenseSummaryResponse.SerializeToString,
        ),
        "GetAdminDashboard": grpc.unary_unary_rpc_method_handler(
            servicer.GetAdminDashboard,
            request_deserializer=AdminDashRequest.FromString,
            response_serializer=AdminDashResponse.SerializeToString,
        ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
        "analytics_service.AnalyticsService", rpc_method_handlers
    )
    server.add_generic_rpc_handlers((generic_handler,))
