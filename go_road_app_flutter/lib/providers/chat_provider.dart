import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/chat_message.dart';
import '../repositories/chat_repository.dart';
import 'api_client_provider.dart';

final chatRepositoryProvider = Provider<ChatRepository>((ref) {
  return ChatRepository(ref.watch(apiClientProvider));
});

final messagesProvider = StateNotifierProvider.family<MessagesNotifier, MessagesState, String>((ref, roomId) {
  return MessagesNotifier(ref.watch(chatRepositoryProvider), roomId);
});

class MessagesState {
  final List<ChatMessage> messages;
  final String? nextCursor;
  final bool hasMore;
  final bool isLoading;
  final String? error;

  const MessagesState({
    this.messages = const [],
    this.nextCursor,
    this.hasMore = false,
    this.isLoading = false,
    this.error,
  });
}

class MessagesNotifier extends StateNotifier<MessagesState> {
  final ChatRepository _repository;
  final String _roomId;

  MessagesNotifier(this._repository, this._roomId) : super(const MessagesState());

  Future<void> loadMessages() async {
    if (state.isLoading) return;
    state = const MessagesState(isLoading: true);
    try {
      final page = await _repository.getMessages(_roomId);
      state = MessagesState(
        messages: page.items,
        nextCursor: page.nextCursor,
        hasMore: page.hasMore,
      );
    } catch (e) {
      state = MessagesState(error: e.toString());
    }
  }

  Future<void> loadMore() async {
    if (state.isLoading || !state.hasMore) return;
    state = MessagesState(messages: state.messages, isLoading: true, hasMore: state.hasMore);
    try {
      final page = await _repository.getMessages(_roomId, cursor: state.nextCursor);
      state = MessagesState(
        messages: [...state.messages, ...page.items],
        nextCursor: page.nextCursor,
        hasMore: page.hasMore,
      );
    } catch (e) {
      state = MessagesState(messages: state.messages, error: e.toString());
    }
  }

  Future<void> sendMessage(String message) async {
    try {
      final msg = await _repository.sendMessage(_roomId, message);
      state = MessagesState(
        messages: [...state.messages, msg],
        hasMore: state.hasMore,
      );
    } catch (e) {
      state = MessagesState(messages: state.messages, error: e.toString());
    }
  }

  void addRealtimeMessage(ChatMessage message) {
    state = MessagesState(
      messages: [...state.messages, message],
      hasMore: state.hasMore,
    );
  }
}
