import '../models/chat_message.dart';
import '../models/pagination.dart';
import 'api_client.dart';

class ChatRepository {
  final ApiClient _client;

  ChatRepository(this._client);

  Future<CursorPage<ChatMessage>> getMessages(String roomId, {String? cursor, int limit = 50}) async {
    final params = <String, dynamic>{'limit': limit};
    if (cursor != null) params['cursor'] = cursor;

    final response = await _client.dio.get('/rooms/$roomId/messages', queryParameters: params);
    return CursorPage.fromJson(
      response.data as Map<String, dynamic>,
      (json) => ChatMessage.fromJson(json),
    );
  }

  Future<ChatMessage> sendMessage(String roomId, String message, {String messageType = 'text'}) async {
    final response = await _client.dio.post(
      '/rooms/$roomId/messages',
      data: {'message': message, 'message_type': messageType},
    );
    return ChatMessage.fromJson(response.data['data'] as Map<String, dynamic>);
  }

  Future<void> markRead(String roomId, String messageId) async {
    await _client.dio.post('/rooms/$roomId/messages/$messageId/read');
  }
}
