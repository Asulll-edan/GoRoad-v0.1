from google.protobuf import descriptor_pb2, descriptor_pool, message_factory


DESCRIPTOR = descriptor_pb2.FileDescriptorProto()
DESCRIPTOR.name = "analytics_service.proto"
DESCRIPTOR.package = "analytics_service"
DESCRIPTOR.syntax = "proto3"

MSG = descriptor_pb2.DescriptorProto
FLD = descriptor_pb2.FieldDescriptorProto
TYPE_STRING = 9
TYPE_DOUBLE = 1
TYPE_INT32 = 5
TYPE_BOOL = 8
LABEL_OPTIONAL = 1
LABEL_REPEATED = 3

_RIDER_STATS_REQ = MSG(name="RiderStatsRequest")
_RIDER_STATS_REQ.field.add(name="user_id", number=1, type=TYPE_STRING, label=LABEL_OPTIONAL)
_RIDER_STATS_REQ.field.add(name="force_refresh", number=2, type=TYPE_BOOL, label=LABEL_OPTIONAL)
DESCRIPTOR.message_type.add().CopyFrom(_RIDER_STATS_REQ)

_RIDER_STATS_RESP = MSG(name="RiderStatsResponse")
_RIDER_STATS_RESP.field.add(name="stats_json", number=1, type=TYPE_STRING, label=LABEL_OPTIONAL)
_RIDER_STATS_RESP.field.add(name="cached", number=2, type=TYPE_STRING, label=LABEL_OPTIONAL)
DESCRIPTOR.message_type.add().CopyFrom(_RIDER_STATS_RESP)

_LEADERBOARD_REQ = MSG(name="LeaderboardRequest")
_LEADERBOARD_REQ.field.add(name="period", number=1, type=TYPE_STRING, label=LABEL_OPTIONAL)
_LEADERBOARD_REQ.field.add(name="limit", number=2, type=TYPE_INT32, label=LABEL_OPTIONAL)
_LEADERBOARD_REQ.field.add(name="cursor", number=3, type=TYPE_STRING, label=LABEL_OPTIONAL)
DESCRIPTOR.message_type.add().CopyFrom(_LEADERBOARD_REQ)

_LEADERBOARD_RESP = MSG(name="LeaderboardResponse")
_LEADERBOARD_RESP.field.add(name="leaderboard_json", number=1, type=TYPE_STRING, label=LABEL_OPTIONAL)
_LEADERBOARD_RESP.field.add(name="cursor", number=2, type=TYPE_STRING, label=LABEL_OPTIONAL)
_LEADERBOARD_RESP.field.add(name="has_more", number=3, type=TYPE_BOOL, label=LABEL_OPTIONAL)
DESCRIPTOR.message_type.add().CopyFrom(_LEADERBOARD_RESP)

_BADGE_EVAL_REQ = MSG(name="BadgeEvalRequest")
_BADGE_EVAL_REQ.field.add(name="user_id", number=1, type=TYPE_STRING, label=LABEL_OPTIONAL)
_BADGE_EVAL_REQ.field.add(name="room_id", number=2, type=TYPE_STRING, label=LABEL_OPTIONAL)
_BADGE_EVAL_REQ.field.add(name="touring_data_json", number=3, type=TYPE_STRING, label=LABEL_OPTIONAL)
DESCRIPTOR.message_type.add().CopyFrom(_BADGE_EVAL_REQ)

_BADGE_EVAL_RESP = MSG(name="BadgeEvalResponse")
_BADGE_EVAL_RESP.field.add(name="new_badge_codes", number=1, label=LABEL_REPEATED, type=TYPE_STRING)
_BADGE_EVAL_RESP.field.add(name="badges_json", number=2, type=TYPE_STRING, label=LABEL_OPTIONAL)
DESCRIPTOR.message_type.add().CopyFrom(_BADGE_EVAL_RESP)

_REPORT_REQ = MSG(name="ReportRequest")
_REPORT_REQ.field.add(name="type", number=1, type=TYPE_STRING, label=LABEL_OPTIONAL)
_REPORT_REQ.field.add(name="params_json", number=2, type=TYPE_STRING, label=LABEL_OPTIONAL)
_REPORT_REQ.field.add(name="format", number=3, type=TYPE_STRING, label=LABEL_OPTIONAL)
DESCRIPTOR.message_type.add().CopyFrom(_REPORT_REQ)

_REPORT_RESP = MSG(name="ReportResponse")
_REPORT_RESP.field.add(name="report_url", number=1, type=TYPE_STRING, label=LABEL_OPTIONAL)
_REPORT_RESP.field.add(name="report_data_json", number=2, type=TYPE_STRING, label=LABEL_OPTIONAL)
DESCRIPTOR.message_type.add().CopyFrom(_REPORT_RESP)

_ROOM_STATS_REQ = MSG(name="RoomStatsRequest")
_ROOM_STATS_REQ.field.add(name="room_id", number=1, type=TYPE_STRING, label=LABEL_OPTIONAL)
DESCRIPTOR.message_type.add().CopyFrom(_ROOM_STATS_REQ)

_ROOM_STATS_RESP = MSG(name="RoomStatsResponse")
_ROOM_STATS_RESP.field.add(name="stats_json", number=1, type=TYPE_STRING, label=LABEL_OPTIONAL)
DESCRIPTOR.message_type.add().CopyFrom(_ROOM_STATS_RESP)

_FUEL_REQ = MSG(name="FuelAnalyticsRequest")
_FUEL_REQ.field.add(name="user_id", number=1, type=TYPE_STRING, label=LABEL_OPTIONAL)
_FUEL_REQ.field.add(name="motor_id", number=2, type=TYPE_STRING, label=LABEL_OPTIONAL)
_FUEL_REQ.field.add(name="months", number=3, type=TYPE_INT32, label=LABEL_OPTIONAL)
DESCRIPTOR.message_type.add().CopyFrom(_FUEL_REQ)

_FUEL_RESP = MSG(name="FuelAnalyticsResponse")
_FUEL_RESP.field.add(name="analytics_json", number=1, type=TYPE_STRING, label=LABEL_OPTIONAL)
DESCRIPTOR.message_type.add().CopyFrom(_FUEL_RESP)

_EXPENSE_REQ = MSG(name="ExpenseSummaryRequest")
_EXPENSE_REQ.field.add(name="user_id", number=1, type=TYPE_STRING, label=LABEL_OPTIONAL)
_EXPENSE_REQ.field.add(name="room_id", number=2, type=TYPE_STRING, label=LABEL_OPTIONAL)
_EXPENSE_REQ.field.add(name="period", number=3, type=TYPE_STRING, label=LABEL_OPTIONAL)
DESCRIPTOR.message_type.add().CopyFrom(_EXPENSE_REQ)

_EXPENSE_RESP = MSG(name="ExpenseSummaryResponse")
_EXPENSE_RESP.field.add(name="summary_json", number=1, type=TYPE_STRING, label=LABEL_OPTIONAL)
DESCRIPTOR.message_type.add().CopyFrom(_EXPENSE_RESP)

_ADMIN_DASH_REQ = MSG(name="AdminDashRequest")
_ADMIN_DASH_REQ.field.add(name="period", number=1, type=TYPE_STRING, label=LABEL_OPTIONAL)
DESCRIPTOR.message_type.add().CopyFrom(_ADMIN_DASH_REQ)

_ADMIN_DASH_RESP = MSG(name="AdminDashResponse")
_ADMIN_DASH_RESP.field.add(name="dashboard_json", number=1, type=TYPE_STRING, label=LABEL_OPTIONAL)
DESCRIPTOR.message_type.add().CopyFrom(_ADMIN_DASH_RESP)

_SVC = descriptor_pb2.ServiceDescriptorProto()
_SVC.name = "AnalyticsService"
_SVC.method.add(name="ComputeRiderStats", input_type=".analytics_service.RiderStatsRequest", output_type=".analytics_service.RiderStatsResponse")
_SVC.method.add(name="ComputeLeaderboard", input_type=".analytics_service.LeaderboardRequest", output_type=".analytics_service.LeaderboardResponse")
_SVC.method.add(name="EvaluateBadges", input_type=".analytics_service.BadgeEvalRequest", output_type=".analytics_service.BadgeEvalResponse")
_SVC.method.add(name="GenerateReport", input_type=".analytics_service.ReportRequest", output_type=".analytics_service.ReportResponse")
_SVC.method.add(name="ComputeRoomStats", input_type=".analytics_service.RoomStatsRequest", output_type=".analytics_service.RoomStatsResponse")
_SVC.method.add(name="ComputeFuelAnalytics", input_type=".analytics_service.FuelAnalyticsRequest", output_type=".analytics_service.FuelAnalyticsResponse")
_SVC.method.add(name="ComputeExpenseSummary", input_type=".analytics_service.ExpenseSummaryRequest", output_type=".analytics_service.ExpenseSummaryResponse")
_SVC.method.add(name="GetAdminDashboard", input_type=".analytics_service.AdminDashRequest", output_type=".analytics_service.AdminDashResponse")
DESCRIPTOR.service.add().CopyFrom(_SVC)

_pool = descriptor_pool.Default()
_serialized = DESCRIPTOR.SerializeToString()
_pool.Add(_serialized)

RiderStatsRequest = message_factory.GetMessageClass(_pool.FindMessageTypeByName("analytics_service.RiderStatsRequest"))
RiderStatsResponse = message_factory.GetMessageClass(_pool.FindMessageTypeByName("analytics_service.RiderStatsResponse"))
LeaderboardRequest = message_factory.GetMessageClass(_pool.FindMessageTypeByName("analytics_service.LeaderboardRequest"))
LeaderboardResponse = message_factory.GetMessageClass(_pool.FindMessageTypeByName("analytics_service.LeaderboardResponse"))
BadgeEvalRequest = message_factory.GetMessageClass(_pool.FindMessageTypeByName("analytics_service.BadgeEvalRequest"))
BadgeEvalResponse = message_factory.GetMessageClass(_pool.FindMessageTypeByName("analytics_service.BadgeEvalResponse"))
ReportRequest = message_factory.GetMessageClass(_pool.FindMessageTypeByName("analytics_service.ReportRequest"))
ReportResponse = message_factory.GetMessageClass(_pool.FindMessageTypeByName("analytics_service.ReportResponse"))
RoomStatsRequest = message_factory.GetMessageClass(_pool.FindMessageTypeByName("analytics_service.RoomStatsRequest"))
RoomStatsResponse = message_factory.GetMessageClass(_pool.FindMessageTypeByName("analytics_service.RoomStatsResponse"))
FuelAnalyticsRequest = message_factory.GetMessageClass(_pool.FindMessageTypeByName("analytics_service.FuelAnalyticsRequest"))
FuelAnalyticsResponse = message_factory.GetMessageClass(_pool.FindMessageTypeByName("analytics_service.FuelAnalyticsResponse"))
ExpenseSummaryRequest = message_factory.GetMessageClass(_pool.FindMessageTypeByName("analytics_service.ExpenseSummaryRequest"))
ExpenseSummaryResponse = message_factory.GetMessageClass(_pool.FindMessageTypeByName("analytics_service.ExpenseSummaryResponse"))
AdminDashRequest = message_factory.GetMessageClass(_pool.FindMessageTypeByName("analytics_service.AdminDashRequest"))
AdminDashResponse = message_factory.GetMessageClass(_pool.FindMessageTypeByName("analytics_service.AdminDashResponse"))
