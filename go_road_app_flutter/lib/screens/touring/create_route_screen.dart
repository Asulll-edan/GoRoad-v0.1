import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../models/route_model.dart';
import '../../providers/api_client_provider.dart';
import '../../repositories/route_repository.dart';

class CreateRouteScreen extends ConsumerStatefulWidget {
  const CreateRouteScreen({super.key});

  @override
  ConsumerState<CreateRouteScreen> createState() => _CreateRouteScreenState();
}

class _CreateRouteScreenState extends ConsumerState<CreateRouteScreen> {
  final _formKey = GlobalKey<FormState>();
  final _nameController = TextEditingController();
  final _descController = TextEditingController();
  final _distanceController = TextEditingController();

  @override
  void dispose() {
    _nameController.dispose();
    _descController.dispose();
    _distanceController.dispose();
    super.dispose();
  }

  Future<void> _create() async {
    if (!_formKey.currentState!.validate()) return;
    try {
      final repo = RouteRepository(ref.read(apiClientProvider));
      await repo.createRoute(CreateRouteRequest(
        name: _nameController.text.trim(),
        description: _descController.text.trim(),
        distanceKm: double.tryParse(_distanceController.text) ?? 0,
      ));
      if (context.mounted) context.pop();
    } catch (e) {
      if (context.mounted) {
        ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Gagal: $e')));
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Buat Rute')),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(24),
        child: Form(
          key: _formKey,
          child: Column(
            children: [
              TextFormField(
                controller: _nameController,
                decoration: const InputDecoration(labelText: 'Nama Rute'),
                validator: (v) => v?.isEmpty ?? true ? 'Nama tidak boleh kosong' : null,
              ),
              const SizedBox(height: 16),
              TextFormField(
                controller: _descController,
                decoration: const InputDecoration(labelText: 'Deskripsi'),
                maxLines: 3,
              ),
              const SizedBox(height: 16),
              TextFormField(
                controller: _distanceController,
                decoration: const InputDecoration(labelText: 'Jarak (km)', suffixText: 'km'),
                keyboardType: TextInputType.number,
              ),
              const SizedBox(height: 32),
              SizedBox(
                width: double.infinity,
                height: 48,
                child: FilledButton(
                  onPressed: _create,
                  child: const Text('Buat Rute'),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
