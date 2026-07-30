import '../models/user.dart';
import 'api_client.dart';

class AuthRepository {
  final ApiClient _client;

  AuthRepository(this._client);

  Future<LoginResponse> login(LoginRequest request) async {
    final response = await _client.dio.post('/auth/login', data: request.toJson());
    final loginResponse = LoginResponse.fromJson(response.data as Map<String, dynamic>);
    await _client.setTokens(
      loginResponse.tokens.accessToken,
      loginResponse.tokens.refreshToken,
    );
    return loginResponse;
  }

  Future<LoginResponse> register(RegisterRequest request) async {
    final response = await _client.dio.post('/auth/register', data: request.toJson());
    final loginResponse = LoginResponse.fromJson(response.data as Map<String, dynamic>);
    await _client.setTokens(
      loginResponse.tokens.accessToken,
      loginResponse.tokens.refreshToken,
    );
    return loginResponse;
  }

  Future<void> logout() async {
    try {
      await _client.dio.post('/auth/logout');
    } finally {
      await _client.clearTokens();
    }
  }

  Future<User> getProfile() async {
    final response = await _client.dio.get('/auth/profile');
    return User.fromJson(response.data['data'] as Map<String, dynamic>);
  }

  Future<User> updateProfile(Map<String, dynamic> data) async {
    final response = await _client.dio.put('/auth/profile', data: data);
    return User.fromJson(response.data['data'] as Map<String, dynamic>);
  }
}
