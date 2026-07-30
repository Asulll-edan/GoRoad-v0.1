import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/user.dart';
import '../repositories/auth_repository.dart';
import 'api_client_provider.dart';

final authRepositoryProvider = Provider<AuthRepository>((ref) {
  return AuthRepository(ref.watch(apiClientProvider));
});

final authStateProvider = StateNotifierProvider<AuthNotifier, AuthState>((ref) {
  return AuthNotifier(ref.watch(authRepositoryProvider));
});

class AuthState {
  final User? user;
  final bool isLoading;
  final String? error;

  const AuthState({this.user, this.isLoading = false, this.error});

  bool get isAuthenticated => user != null;
}

class AuthNotifier extends StateNotifier<AuthState> {
  final AuthRepository _repository;

  AuthNotifier(this._repository) : super(const AuthState());

  Future<void> login(String email, String password) async {
    state = const AuthState(isLoading: true);
    try {
      final response = await _repository.login(LoginRequest(email: email, password: password));
      state = AuthState(user: response.user);
    } catch (e) {
      state = AuthState(error: e.toString());
    }
  }

  Future<void> register(RegisterRequest request) async {
    state = const AuthState(isLoading: true);
    try {
      final response = await _repository.register(request);
      state = AuthState(user: response.user);
    } catch (e) {
      state = AuthState(error: e.toString());
    }
  }

  Future<void> loadProfile() async {
    try {
      final user = await _repository.getProfile();
      state = AuthState(user: user);
    } catch (_) {
      state = const AuthState();
    }
  }

  Future<void> logout() async {
    await _repository.logout();
    state = const AuthState();
  }

  void clearError() {
    state = AuthState(user: state.user);
  }
}
