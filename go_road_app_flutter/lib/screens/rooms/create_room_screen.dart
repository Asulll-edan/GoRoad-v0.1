import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../models/room.dart';
import '../../providers/api_client_provider.dart';
import '../../repositories/room_repository.dart';

class CreateRoomScreen extends ConsumerStatefulWidget {
  const CreateRoomScreen({super.key});

  @override
  ConsumerState<CreateRoomScreen> createState() => _CreateRoomScreenState();
}

class _CreateRoomScreenState extends ConsumerState<CreateRoomScreen> {
  final _formKey = GlobalKey<FormState>();
  final _nameController = TextEditingController();
  final _descController = TextEditingController();
  final _originController = TextEditingController();
  final _destController = TextEditingController();
  bool _isLoading = false;

  @override
  void dispose() {
    _nameController.dispose();
    _descController.dispose();
    _originController.dispose();
    _destController.dispose();
    super.dispose();
  }

  Future<void> _create() async {
    if (!_formKey.currentState!.validate()) return;
    setState(() => _isLoading = true);
    try {
      final repo = RoomRepository(ref.read(apiClientProvider));
      await repo.createRoom(CreateRoomRequest(
        name: _nameController.text.trim(),
        description: _descController.text.trim(),
        originCity: _originController.text.trim(),
        destinationCity: _destController.text.trim(),
      ));
      if (context.mounted) context.pop();
    } catch (e) {
      if (context.mounted) {
        ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Gagal: $e')));
      }
    } finally {
      setState(() => _isLoading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Buat Room Baru')),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(24),
        child: Form(
          key: _formKey,
          child: Column(
            children: [
              TextFormField(
                controller: _nameController,
                decoration: const InputDecoration(labelText: 'Nama Room'),
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
                controller: _originController,
                decoration: const InputDecoration(labelText: 'Kota Asal', prefixIcon: Icon(Icons.location_on)),
              ),
              const SizedBox(height: 16),
              TextFormField(
                controller: _destController,
                decoration: const InputDecoration(labelText: 'Kota Tujuan', prefixIcon: Icon(Icons.flag)),
              ),
              const SizedBox(height: 32),
              SizedBox(
                width: double.infinity,
                height: 48,
                child: FilledButton(
                  onPressed: _isLoading ? null : _create,
                  child: _isLoading
                      ? const SizedBox(width: 24, height: 24, child: CircularProgressIndicator(strokeWidth: 2))
                      : const Text('Buat Room'),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
