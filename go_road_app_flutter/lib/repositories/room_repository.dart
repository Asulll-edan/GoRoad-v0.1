import '../models/room.dart';
import '../models/pagination.dart';
import 'api_client.dart';

class RoomRepository {
  final ApiClient _client;

  RoomRepository(this._client);

  Future<CursorPage<Room>> getRooms({String? cursor, int limit = 20}) async {
    final queryParams = <String, dynamic>{'limit': limit};
    if (cursor != null) queryParams['cursor'] = cursor;

    final response = await _client.dio.get('/rooms', queryParameters: queryParams);
    return CursorPage.fromJson(
      response.data as Map<String, dynamic>,
      (json) => Room.fromJson(json),
    );
  }

  Future<CursorPage<Room>> getMyRooms({String? cursor, int limit = 20}) async {
    final queryParams = <String, dynamic>{'limit': limit};
    if (cursor != null) queryParams['cursor'] = cursor;

    final response = await _client.dio.get('/rooms/mine', queryParameters: queryParams);
    return CursorPage.fromJson(
      response.data as Map<String, dynamic>,
      (json) => Room.fromJson(json),
    );
  }

  Future<Room> getRoom(String id, {String? include}) async {
    final queryParams = <String, dynamic>{};
    if (include != null) queryParams['include'] = include;

    final response = await _client.dio.get('/rooms/$id', queryParameters: queryParams);
    return Room.fromJson(response.data['data'] as Map<String, dynamic>);
  }

  Future<Room> createRoom(CreateRoomRequest request) async {
    final response = await _client.dio.post('/rooms', data: request.toJson());
    return Room.fromJson(response.data['data'] as Map<String, dynamic>);
  }

  Future<void> joinRoom(String roomId) async {
    await _client.dio.post('/rooms/$roomId/join');
  }

  Future<void> leaveRoom(String roomId) async {
    await _client.dio.post('/rooms/$roomId/leave');
  }

  Future<List<RoomMember>> getMembers(String roomId) async {
    final response = await _client.dio.get('/rooms/$roomId/members');
    return (response.data['data'] as List<dynamic>)
        .map((e) => RoomMember.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  Future<CursorPage<Room>> searchRooms(String query, {String? cursor, int limit = 20}) async {
    final params = <String, dynamic>{'q': query, 'limit': limit};
    if (cursor != null) params['cursor'] = cursor;

    final response = await _client.dio.get('/rooms/search', queryParameters: params);
    return CursorPage.fromJson(
      response.data as Map<String, dynamic>,
      (json) => Room.fromJson(json),
    );
  }
}
