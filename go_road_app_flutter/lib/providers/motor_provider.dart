import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/motor.dart';
import '../repositories/motor_repository.dart';
import 'api_client_provider.dart';

final motorRepositoryProvider = Provider<MotorRepository>((ref) {
  return MotorRepository(ref.watch(apiClientProvider));
});

final motorsProvider = StateNotifierProvider<MotorsNotifier, MotorsState>((ref) {
  return MotorsNotifier(ref.watch(motorRepositoryProvider));
});

class MotorsState {
  final List<Motor> motors;
  final bool isLoading;
  final String? error;

  const MotorsState({this.motors = const [], this.isLoading = false, this.error});
}

class MotorsNotifier extends StateNotifier<MotorsState> {
  final MotorRepository _repository;

  MotorsNotifier(this._repository) : super(const MotorsState());

  Future<void> loadMotors() async {
    state = const MotorsState(isLoading: true);
    try {
      final motors = await _repository.getMotors();
      state = MotorsState(motors: motors);
    } catch (e) {
      state = MotorsState(error: e.toString());
    }
  }

  Future<void> addMotor(Map<String, dynamic> data) async {
    try {
      await _repository.createMotor(data);
      await loadMotors();
    } catch (e) {
      state = MotorsState(motors: state.motors, error: e.toString());
    }
  }
}
