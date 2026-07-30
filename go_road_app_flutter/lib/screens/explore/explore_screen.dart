import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:cached_network_image/cached_network_image.dart';
import '../../providers/room_provider.dart';

class ExploreScreen extends ConsumerStatefulWidget {
  const ExploreScreen({super.key});

  @override
  ConsumerState<ExploreScreen> createState() => _ExploreScreenState();
}

class _ExploreScreenState extends ConsumerState<ExploreScreen> {
  final _searchController = TextEditingController();

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: TextField(
          controller: _searchController,
          decoration: const InputDecoration(
            hintText: 'Cari room, rute, atau tempat...',
            border: InputBorder.none,
          ),
        ),
      ),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('Jelajahi', style: theme.textTheme.headlineSmall?.copyWith(fontWeight: FontWeight.bold)),
            const SizedBox(height: 16),
            _categoryChip('Touring'),
            _categoryChip('Weekly Meetup'),
            _categoryChip('Adventure'),
            _categoryChip('Off-road'),
            const SizedBox(height: 24),
            Text('Room Populer', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
            const SizedBox(height: 12),
            Expanded(
              child: ref.watch(roomsProvider).rooms.isEmpty
                  ? const Center(child: Text('Belum ada room'))
                  : ListView.builder(
                      itemCount: ref.watch(roomsProvider).rooms.length,
                      itemBuilder: (context, index) {
                        final room = ref.watch(roomsProvider).rooms[index];
                        return ListTile(
                          leading: ClipRRect(
                            borderRadius: BorderRadius.circular(8),
                            child: room.coverUrl != null
                                ? CachedNetworkImage(imageUrl: room.coverUrl!, width: 48, height: 48, fit: BoxFit.cover)
                                : Container(width: 48, height: 48, color: Colors.grey[200], child: const Icon(Icons.motorcycle)),
                          ),
                          title: Text(room.name),
                          subtitle: Text('${room.memberCount} members'),
                        );
                      },
                    ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _categoryChip(String label) {
    return Padding(
      padding: const EdgeInsets.only(right: 8, bottom: 8),
      child: FilterChip(
        label: Text(label),
        onSelected: (_) {},
      ),
    );
  }
}
