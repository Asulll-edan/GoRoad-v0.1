import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:cached_network_image/cached_network_image.dart';
import '../../providers/auth_provider.dart';

class ProfileScreen extends ConsumerWidget {
  const ProfileScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final authState = ref.watch(authStateProvider);
    final user = authState.user;
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Profil'),
        actions: [
          IconButton(
            icon: const Icon(Icons.settings),
            onPressed: () {},
          ),
        ],
      ),
      body: user == null
          ? const Center(child: CircularProgressIndicator())
          : ListView(
              padding: const EdgeInsets.all(16),
              children: [
                Center(
                  child: Column(
                    children: [
                      CircleAvatar(
                        radius: 48,
                        backgroundImage: user.avatarUrl != null ? CachedNetworkImageProvider(user.avatarUrl!) : null,
                        child: user.avatarUrl == null
                            ? Text(user.username[0].toUpperCase(), style: const TextStyle(fontSize: 32))
                            : null,
                      ),
                      const SizedBox(height: 16),
                      Text(user.displayName ?? user.username, style: theme.textTheme.headlineSmall?.copyWith(fontWeight: FontWeight.bold)),
                      Text('@${user.username}', style: TextStyle(color: Colors.grey[600])),
                      if (user.bio != null) ...[
                        const SizedBox(height: 8),
                        Text(user.bio!, style: TextStyle(color: Colors.grey[600])),
                      ],
                    ],
                  ),
                ),
                const SizedBox(height: 32),
                Card(
                  child: Column(
                    children: [
                      _menuItem(Icons.motorcycle, 'Motor Saya', () => context.push('/routes')),
                      const Divider(height: 1),
                      _menuItem(Icons.route, 'Rute Saya', () => context.push('/routes')),
                      const Divider(height: 1),
                      _menuItem(Icons.explore, 'Jelajahi', () => context.push('/explore')),
                      const Divider(height: 1),
                      _menuItem(Icons.person, 'Edit Profil', () => context.push('/profile/edit')),
                    ],
                  ),
                ),
                const SizedBox(height: 16),
                Card(
                  child: _menuItem(Icons.logout, 'Keluar', () {
                    ref.read(authStateProvider.notifier).logout();
                  }, color: Colors.red),
                ),
              ],
            ),
    );
  }

  Widget _menuItem(IconData icon, String title, VoidCallback onTap, {Color? color}) {
    return ListTile(
      leading: Icon(icon, color: color),
      title: Text(title, style: color != null ? TextStyle(color: color) : null),
      trailing: const Icon(Icons.chevron_right),
      onTap: onTap,
    );
  }
}
