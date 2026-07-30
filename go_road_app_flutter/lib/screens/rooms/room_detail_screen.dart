import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:cached_network_image/cached_network_image.dart';
import '../../repositories/room_repository.dart';
import '../../providers/api_client_provider.dart';

class RoomDetailScreen extends ConsumerWidget {
  final String roomId;

  const RoomDetailScreen({super.key, required this.roomId});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final apiClient = ref.watch(apiClientProvider);
    final roomRepo = RoomRepository(apiClient);

    return FutureBuilder(
      future: roomRepo.getRoom(roomId, include: 'members'),
      builder: (context, snapshot) {
        if (snapshot.connectionState == ConnectionState.waiting) {
          return Scaffold(
            appBar: AppBar(title: const Text('Room Detail')),
            body: const Center(child: CircularProgressIndicator()),
          );
        }

        if (snapshot.hasError) {
          return Scaffold(
            appBar: AppBar(title: const Text('Room Detail')),
            body: Center(child: Text('Error: ${snapshot.error}')),
          );
        }

        final room = snapshot.data!;

        return Scaffold(
          appBar: AppBar(
            title: Text(room.name),
            actions: [
              IconButton(
                icon: const Icon(Icons.chat),
                onPressed: () => context.push('/rooms/$roomId/chat'),
              ),
            ],
          ),
          body: SingleChildScrollView(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                if (room.coverUrl != null)
                  CachedNetworkImage(
                    imageUrl: room.coverUrl!,
                    width: double.infinity,
                    height: 200,
                    fit: BoxFit.cover,
                    placeholder: (_, __) => Container(height: 200, color: Colors.grey[200]),
                  ),
                Padding(
                  padding: const EdgeInsets.all(16),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(room.name, style: theme.textTheme.headlineSmall?.copyWith(fontWeight: FontWeight.bold)),
                      if (room.description != null) ...[
                        const SizedBox(height: 8),
                        Text(room.description!, style: TextStyle(color: Colors.grey[600], fontSize: 16)),
                      ],
                      const SizedBox(height: 16),
                      Row(
                        children: [
                          _infoChip(Icons.people, '${room.memberCount} Member'),
                          if (room.originCity != null) _infoChip(Icons.location_on, room.originCity!),
                          if (room.status != 'active') _infoChip(Icons.info, room.status),
                        ],
                      ),
                      const SizedBox(height: 24),
                      Row(
                        children: [
                          Text('Members', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
                          const Spacer(),
                          TextButton(
                            onPressed: () => context.push('/rooms/$roomId/chat'),
                            child: const Text('Chat Room'),
                          ),
                        ],
                      ),
                      if (room.members != null)
                        ...room.members!.map((m) => ListTile(
                              leading: CircleAvatar(
                                backgroundImage: m.avatarUrl != null ? CachedNetworkImageProvider(m.avatarUrl!) : null,
                                child: m.avatarUrl == null ? Text(m.username?[0].toUpperCase() ?? '?') : null,
                              ),
                              title: Text(m.username ?? 'Unknown'),
                              subtitle: Text(m.role),
                            )),
                    ],
                  ),
                ),
              ],
            ),
          ),
          bottomNavigationBar: SafeArea(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: FilledButton.icon(
                onPressed: () async {
                  await roomRepo.joinRoom(roomId);
                  if (context.mounted) {
                    ScaffoldMessenger.of(context).showSnackBar(
                      const SnackBar(content: Text('Bergabung ke room')),
                    );
                  }
                },
                icon: const Icon(Icons.login),
                label: const Text('Join Room'),
              ),
            ),
          ),
        );
      },
    );
  }

  Widget _infoChip(IconData icon, String label) {
    return Padding(
      padding: const EdgeInsets.only(right: 12),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(icon, size: 16, color: Colors.grey),
          const SizedBox(width: 4),
          Text(label, style: const TextStyle(color: Colors.grey)),
        ],
      ),
    );
  }
}
