import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/room.dart';
import '../repositories/room_repository.dart';
import 'api_client_provider.dart';

final roomRepositoryProvider = Provider<RoomRepository>((ref) {
  return RoomRepository(ref.watch(apiClientProvider));
});

final roomsProvider = StateNotifierProvider<RoomsNotifier, RoomsState>((ref) {
  return RoomsNotifier(ref.watch(roomRepositoryProvider));
});

class RoomsState {
  final List<Room> rooms;
  final String? nextCursor;
  final bool hasMore;
  final bool isLoading;
  final String? error;

  const RoomsState({
    this.rooms = const [],
    this.nextCursor,
    this.hasMore = false,
    this.isLoading = false,
    this.error,
  });
}

class RoomsNotifier extends StateNotifier<RoomsState> {
  final RoomRepository _repository;

  RoomsNotifier(this._repository) : super(const RoomsState());

  Future<void> loadRooms() async {
    if (state.isLoading) return;
    state = RoomsState(isLoading: true);
    try {
      final page = await _repository.getRooms();
      state = RoomsState(
        rooms: page.items,
        nextCursor: page.nextCursor,
        hasMore: page.hasMore,
      );
    } catch (e) {
      state = RoomsState(error: e.toString());
    }
  }

  Future<void> loadMore() async {
    if (state.isLoading || !state.hasMore) return;
    state = RoomsState(rooms: state.rooms, isLoading: true, hasMore: state.hasMore);
    try {
      final page = await _repository.getRooms(cursor: state.nextCursor);
      state = RoomsState(
        rooms: [...state.rooms, ...page.items],
        nextCursor: page.nextCursor,
        hasMore: page.hasMore,
      );
    } catch (e) {
      state = RoomsState(rooms: state.rooms, error: e.toString());
    }
  }

  Future<void> refresh() async {
    state = const RoomsState(isLoading: true);
    await loadRooms();
  }
}
