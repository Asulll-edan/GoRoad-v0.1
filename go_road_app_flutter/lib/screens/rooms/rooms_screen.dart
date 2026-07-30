import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:cached_network_image/cached_network_image.dart';
import '../../providers/room_provider.dart';
import '../../providers/auth_provider.dart';

class RoomsScreen extends ConsumerStatefulWidget {
  const RoomsScreen({super.key});

  @override
  ConsumerState<RoomsScreen> createState() => _RoomsScreenState();
}

class _RoomsScreenState extends ConsumerState<RoomsScreen> {
  final _scrollController = ScrollController();

  @override
  void initState() {
    super.initState();
    Future.microtask(() => ref.read(roomsProvider.notifier).loadRooms());
    _scrollController.addListener(_onScroll);
  }

  @override
  void dispose() {
    _scrollController.dispose();
    super.dispose();
  }

  void _onScroll() {
    if (_scrollController.position.pixels >= _scrollController.position.maxScrollExtent - 200) {
      ref.read(roomsProvider.notifier).loadMore();
    }
  }

  @override
  Widget build(BuildContext context) {
    final roomsState = ref.watch(roomsProvider);
    final authState = ref.watch(authStateProvider);
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Go Road'),
        actions: [
          IconButton(
            icon: const Icon(Icons.person),
            onPressed: () => context.push('/profile'),
          ),
          if (authState.isAuthenticated)
            IconButton(
              icon: const Icon(Icons.logout),
              onPressed: () => ref.read(authStateProvider.notifier).logout(),
            ),
        ],
      ),
      body: RefreshIndicator(
        onRefresh: () => ref.read(roomsProvider.notifier).refresh(),
        child: roomsState.isLoading && roomsState.rooms.isEmpty
            ? const Center(child: CircularProgressIndicator())
            : roomsState.rooms.isEmpty
                ? ListView(
                    children: [
                      SizedBox(
                        height: MediaQuery.of(context).size.height * 0.6,
                        child: const Center(
                          child: Column(
                            mainAxisSize: MainAxisSize.min,
                            children: [
                              Icon(Icons.motorcycle, size: 64, color: Colors.grey),
                              SizedBox(height: 16),
                              Text('Belum ada room', style: TextStyle(fontSize: 18, color: Colors.grey)),
                              Text('Buat room baru untuk mulai touring', style: TextStyle(color: Colors.grey)),
                            ],
                          ),
                        ),
                      ),
                    ],
                  )
                : ListView.builder(
                    controller: _scrollController,
                    padding: const EdgeInsets.all(12),
                    itemCount: roomsState.rooms.length + (roomsState.isLoading ? 1 : 0),
                    itemBuilder: (context, index) {
                      if (index >= roomsState.rooms.length) {
                        return const Padding(
                          padding: EdgeInsets.all(16),
                          child: Center(child: CircularProgressIndicator()),
                        );
                      }
                      final room = roomsState.rooms[index];
                      return Card(
                        margin: const EdgeInsets.only(bottom: 12),
                        child: InkWell(
                          borderRadius: BorderRadius.circular(12),
                          onTap: () => context.push('/rooms/${room.id}'),
                          child: Padding(
                            padding: const EdgeInsets.all(16),
                            child: Row(
                              children: [
                                ClipRRect(
                                  borderRadius: BorderRadius.circular(8),
                                  child: room.coverUrl != null
                                      ? CachedNetworkImage(
                                          imageUrl: room.coverUrl!,
                                          width: 64,
                                          height: 64,
                                          fit: BoxFit.cover,
                                          placeholder: (_, __) => Container(color: Colors.grey[200]),
                                        )
                                      : Container(
                                          width: 64,
                                          height: 64,
                                          color: theme.colorScheme.primaryContainer,
                                          child: Icon(Icons.motorcycle, color: theme.colorScheme.primary),
                                        ),
                                ),
                                const SizedBox(width: 16),
                                Expanded(
                                  child: Column(
                                    crossAxisAlignment: CrossAxisAlignment.start,
                                    children: [
                                      Text(room.name, style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
                                      if (room.description != null) ...[
                                        const SizedBox(height: 4),
                                        Text(room.description!, maxLines: 2, overflow: TextOverflow.ellipsis, style: TextStyle(color: Colors.grey[600])),
                                      ],
                                      const SizedBox(height: 8),
                                      Row(
                                        children: [
                                          const Icon(Icons.people, size: 16, color: Colors.grey),
                                          const SizedBox(width: 4),
                                          Text('${room.memberCount}/${room.maxMembers}', style: const TextStyle(fontSize: 12, color: Colors.grey)),
                                          if (room.originCity != null) ...[
                                            const SizedBox(width: 16),
                                            const Icon(Icons.location_on, size: 16, color: Colors.grey),
                                            const SizedBox(width: 4),
                                            Text(room.originCity!, style: const TextStyle(fontSize: 12, color: Colors.grey)),
                                          ],
                                        ],
                                      ),
                                    ],
                                  ),
                                ),
                                const Icon(Icons.chevron_right),
                              ],
                            ),
                          ),
                        ),
                      );
                    },
                  ),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => context.push('/rooms/create'),
        child: const Icon(Icons.add),
      ),
    );
  }
}
