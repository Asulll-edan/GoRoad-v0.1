from google.protobuf import descriptor_pb2, descriptor_pool, message_factory, reflection


DESCRIPTOR = descriptor_pb2.FileDescriptorProto()
DESCRIPTOR.name = "ai_service.proto"
DESCRIPTOR.package = "ai_service"
DESCRIPTOR.syntax = "proto3"

MSG = descriptor_pb2.DescriptorProto
FLD = descriptor_pb2.FieldDescriptorProto
TYPE_STRING = 9
TYPE_DOUBLE = 1
TYPE_BOOL = 8
TYPE_INT32 = 5
LABEL_OPTIONAL = 1
LABEL_REPEATED = 3

_CHAT_REQUEST = MSG(name="ChatRequest")
_CHAT_REQUEST.field.add(name="user_id", number=1, type=TYPE_STRING, label=LABEL_OPTIONAL)
_CHAT_REQUEST.field.add(name="room_id", number=2, type=TYPE_STRING, label=LABEL_OPTIONAL)
_CHAT_REQUEST.field.add(name="message", number=3, type=TYPE_STRING, label=LABEL_OPTIONAL)
_CHAT_REQUEST.field.add(name="context", number=4, label=LABEL_REPEATED, type=TYPE_STRING)
DESCRIPTOR.message_type.add().CopyFrom(_CHAT_REQUEST)

_CONTEXT_MSG = MSG(name="ContextMessage")
_CONTEXT_MSG.field.add(name="role", number=1, type=TYPE_STRING, label=LABEL_OPTIONAL)
_CONTEXT_MSG.field.add(name="content", number=2, type=TYPE_STRING, label=LABEL_OPTIONAL)
DESCRIPTOR.message_type.add().CopyFrom(_CONTEXT_MSG)

_CHAT_RESP = MSG(name="ChatResponse")
_CHAT_RESP.field.add(name="content", number=1, type=TYPE_STRING, label=LABEL_OPTIONAL)
_CHAT_RESP.field.add(name="is_final", number=2, type=TYPE_BOOL, label=LABEL_OPTIONAL)
DESCRIPTOR.message_type.add().CopyFrom(_CHAT_RESP)

_ITIN_REQ = MSG(name="ItineraryRequest")
_ITIN_REQ.field.add(name="route_id", number=1, type=TYPE_STRING, label=LABEL_OPTIONAL)
_ITIN_REQ.field.add(name="start_location", number=2, type=TYPE_STRING, label=LABEL_OPTIONAL)
_ITIN_REQ.field.add(name="end_location", number=3, type=TYPE_STRING, label=LABEL_OPTIONAL)
_ITIN_REQ.field.add(name="duration_days", number=4, type=TYPE_INT32, label=LABEL_OPTIONAL)
_ITIN_REQ.field.add(name="rider_count", number=5, type=TYPE_INT32, label=LABEL_OPTIONAL)
_ITIN_REQ.field.add(name="motor_ids", number=6, label=LABEL_REPEATED, type=TYPE_STRING)
_ITIN_REQ.field.add(name="preferences", number=7, label=LABEL_REPEATED, type=TYPE_STRING)
DESCRIPTOR.message_type.add().CopyFrom(_ITIN_REQ)

_ITIN_RESP = MSG(name="ItineraryResponse")
_ITIN_RESP.field.add(name="itinerary_json", number=1, type=TYPE_STRING, label=LABEL_OPTIONAL)
_ITIN_RESP.field.add(name="cached", number=2, type=TYPE_STRING, label=LABEL_OPTIONAL)
DESCRIPTOR.message_type.add().CopyFrom(_ITIN_RESP)

_COST_REQ = MSG(name="CostEstimationRequest")
_COST_REQ.field.add(name="route_id", number=1, type=TYPE_STRING, label=LABEL_OPTIONAL)
_COST_REQ.field.add(name="motor_ids", number=2, label=LABEL_REPEATED, type=TYPE_STRING)
_COST_REQ.field.add(name="rider_count", number=3, type=TYPE_INT32, label=LABEL_OPTIONAL)
_COST_REQ.field.add(name="duration_days", number=4, type=TYPE_INT32, label=LABEL_OPTIONAL)
_COST_REQ.field.add(name="fuel_type", number=5, type=TYPE_STRING, label=LABEL_OPTIONAL)
DESCRIPTOR.message_type.add().CopyFrom(_COST_REQ)

_COST_RESP = MSG(name="CostEstimationResponse")
_COST_RESP.field.add(name="estimate_json", number=1, type=TYPE_STRING, label=LABEL_OPTIONAL)
_COST_RESP.field.add(name="cached", number=2, type=TYPE_STRING, label=LABEL_OPTIONAL)
DESCRIPTOR.message_type.add().CopyFrom(_COST_RESP)

_ROUTE_REQ = MSG(name="RouteAdviceRequest")
_ROUTE_REQ.field.add(name="origin", number=1, type=TYPE_STRING, label=LABEL_OPTIONAL)
_ROUTE_REQ.field.add(name="destination", number=2, type=TYPE_STRING, label=LABEL_OPTIONAL)
_ROUTE_REQ.field.add(name="waypoints", number=3, label=LABEL_REPEATED, type=TYPE_STRING)
_ROUTE_REQ.field.add(name="preferences", number=4, label=LABEL_REPEATED, type=TYPE_STRING)
DESCRIPTOR.message_type.add().CopyFrom(_ROUTE_REQ)

_ROUTE_RESP = MSG(name="RouteAdviceResponse")
_ROUTE_RESP.field.add(name="advice_json", number=1, type=TYPE_STRING, label=LABEL_OPTIONAL)
DESCRIPTOR.message_type.add().CopyFrom(_ROUTE_RESP)

_PACK_REQ = MSG(name="PackingListRequest")
_PACK_REQ.field.add(name="duration_days", number=1, type=TYPE_INT32, label=LABEL_OPTIONAL)
_PACK_REQ.field.add(name="weather_condition", number=2, type=TYPE_STRING, label=LABEL_OPTIONAL)
_PACK_REQ.field.add(name="touring_type", number=3, type=TYPE_STRING, label=LABEL_OPTIONAL)
DESCRIPTOR.message_type.add().CopyFrom(_PACK_REQ)

_PACK_RESP = MSG(name="PackingListResponse")
_PACK_RESP.field.add(name="packing_list_json", number=1, type=TYPE_STRING, label=LABEL_OPTIONAL)
DESCRIPTOR.message_type.add().CopyFrom(_PACK_RESP)

_SAFETY_REQ = MSG(name="SafetyRequest")
_SAFETY_REQ.field.add(name="route_id", number=1, type=TYPE_STRING, label=LABEL_OPTIONAL)
_SAFETY_REQ.field.add(name="weather_condition", number=2, type=TYPE_STRING, label=LABEL_OPTIONAL)
_SAFETY_REQ.field.add(name="rider_count", number=3, type=TYPE_INT32, label=LABEL_OPTIONAL)
_SAFETY_REQ.field.add(name="skill_level", number=4, type=TYPE_STRING, label=LABEL_OPTIONAL)
DESCRIPTOR.message_type.add().CopyFrom(_SAFETY_REQ)

_SAFETY_RESP = MSG(name="SafetyResponse")
_SAFETY_RESP.field.add(name="advice_json", number=1, type=TYPE_STRING, label=LABEL_OPTIONAL)
DESCRIPTOR.message_type.add().CopyFrom(_SAFETY_RESP)

_POI_REQ = MSG(name="POIRecommendRequest")
_POI_REQ.field.add(name="lat", number=1, type=TYPE_DOUBLE, label=LABEL_OPTIONAL)
_POI_REQ.field.add(name="lng", number=2, type=TYPE_DOUBLE, label=LABEL_OPTIONAL)
_POI_REQ.field.add(name="radius_km", number=3, type=TYPE_DOUBLE, label=LABEL_OPTIONAL)
_POI_REQ.field.add(name="types", number=4, label=LABEL_REPEATED, type=TYPE_STRING)
_POI_REQ.field.add(name="route_id", number=5, type=TYPE_STRING, label=LABEL_OPTIONAL)
DESCRIPTOR.message_type.add().CopyFrom(_POI_REQ)

_POI_RESP = MSG(name="POIRecommendResponse")
_POI_RESP.field.add(name="pois_json", number=1, type=TYPE_STRING, label=LABEL_OPTIONAL)
_POI_RESP.field.add(name="cached", number=2, type=TYPE_STRING, label=LABEL_OPTIONAL)
DESCRIPTOR.message_type.add().CopyFrom(_POI_RESP)

_SVC = descriptor_pb2.ServiceDescriptorProto()
_SVC.name = "AIService"
_SVC.method.add(name="ChatStream", input_type=".ai_service.ChatRequest", output_type=".ai_service.ChatResponse", server_streaming=True)
_SVC.method.add(name="GenerateItinerary", input_type=".ai_service.ItineraryRequest", output_type=".ai_service.ItineraryResponse")
_SVC.method.add(name="EstimateCost", input_type=".ai_service.CostEstimationRequest", output_type=".ai_service.CostEstimationResponse")
_SVC.method.add(name="AdviseRoute", input_type=".ai_service.RouteAdviceRequest", output_type=".ai_service.RouteAdviceResponse")
_SVC.method.add(name="GeneratePackingList", input_type=".ai_service.PackingListRequest", output_type=".ai_service.PackingListResponse")
_SVC.method.add(name="AdviseSafety", input_type=".ai_service.SafetyRequest", output_type=".ai_service.SafetyResponse")
_SVC.method.add(name="RecommendPOI", input_type=".ai_service.POIRecommendRequest", output_type=".ai_service.POIRecommendResponse")
DESCRIPTOR.service.add().CopyFrom(_SVC)

_pool = descriptor_pool.Default()
_serialized = DESCRIPTOR.SerializeToString()
_file_desc = _pool.Add(_serialized)

ChatRequest = message_factory.GetMessageClass(_pool.FindMessageTypeByName("ai_service.ChatRequest"))
ContextMessage = message_factory.GetMessageClass(_pool.FindMessageTypeByName("ai_service.ContextMessage"))
ChatResponse = message_factory.GetMessageClass(_pool.FindMessageTypeByName("ai_service.ChatResponse"))
ItineraryRequest = message_factory.GetMessageClass(_pool.FindMessageTypeByName("ai_service.ItineraryRequest"))
ItineraryResponse = message_factory.GetMessageClass(_pool.FindMessageTypeByName("ai_service.ItineraryResponse"))
CostEstimationRequest = message_factory.GetMessageClass(_pool.FindMessageTypeByName("ai_service.CostEstimationRequest"))
CostEstimationResponse = message_factory.GetMessageClass(_pool.FindMessageTypeByName("ai_service.CostEstimationResponse"))
RouteAdviceRequest = message_factory.GetMessageClass(_pool.FindMessageTypeByName("ai_service.RouteAdviceRequest"))
RouteAdviceResponse = message_factory.GetMessageClass(_pool.FindMessageTypeByName("ai_service.RouteAdviceResponse"))
PackingListRequest = message_factory.GetMessageClass(_pool.FindMessageTypeByName("ai_service.PackingListRequest"))
PackingListResponse = message_factory.GetMessageClass(_pool.FindMessageTypeByName("ai_service.PackingListResponse"))
SafetyRequest = message_factory.GetMessageClass(_pool.FindMessageTypeByName("ai_service.SafetyRequest"))
SafetyResponse = message_factory.GetMessageClass(_pool.FindMessageTypeByName("ai_service.SafetyResponse"))
POIRecommendRequest = message_factory.GetMessageClass(_pool.FindMessageTypeByName("ai_service.POIRecommendRequest"))
POIRecommendResponse = message_factory.GetMessageClass(_pool.FindMessageTypeByName("ai_service.POIRecommendResponse"))
