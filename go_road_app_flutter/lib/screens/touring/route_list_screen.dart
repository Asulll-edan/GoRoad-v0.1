import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../providers/api_client_provider.dart';
import '../../repositories/route_repository.dart';

class RouteListScreen extends ConsumerWidget {
  const RouteListScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Rute Saya'),
        actions: [
          IconButton(
            icon: const Icon(Icons.add),
            onPressed: () => context.push('/routes/create'),
          ),
        ],
      ),
      body: FutureBuilder(
        future: RouteRepository(ref.read(apiClientProvider)).getRoutes(),
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return const Center(child: CircularProgressIndicator());
          }
          if (snapshot.hasError) {
            return Center(child: Text('Error: ${snapshot.error}'));
          }
          final page = snapshot.data!;
          if (page.items.isEmpty) {
            return const Center(child: Text('Belum ada rute'));
          }
          return ListView.builder(
            padding: const EdgeInsets.all(12),
            itemCount: page.items.length,
            itemBuilder: (context, index) {
              final route = page.items[index];
              return Card(
                margin: const EdgeInsets.only(bottom: 12),
                child: ListTile(
                  title: Text(route.name, style: const TextStyle(fontWeight: FontWeight.bold)),
                  subtitle: Text('${route.distanceKm.toStringAsFixed(1)} km'),
                  trailing: Text(route.status, style: TextStyle(color: route.status == 'completed' ? Colors.green : Colors.orange)),
                ),
              );
            },
          );
        },
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => context.push('/routes/create'),
        child: const Icon(Icons.add),
      ),
    );
  }
}
