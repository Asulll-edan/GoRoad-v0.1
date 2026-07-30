import 'package:dio/dio.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class ApiClient {
  static const String _tokenKey = 'access_token';
  static const String _refreshKey = 'refresh_token';

  late final Dio _dio;
  final FlutterSecureStorage _storage;
  String? _baseUrl;

  ApiClient({String? baseUrl, Dio? dio}) : _storage = const FlutterSecureStorage() {
    _dio = dio ?? Dio(BaseOptions(
      connectTimeout: const Duration(seconds: 15),
      receiveTimeout: const Duration(seconds: 15),
      headers: {'Content-Type': 'application/json'},
    ));
    _baseUrl = baseUrl;
    _setupInterceptors();
  }

  Dio get dio => _dio;
  String get baseUrl => _baseUrl ?? 'http://10.0.2.2:8080/api/v1';

  void _setupInterceptors() {
    _dio.interceptors.add(InterceptorsWrapper(
      onRequest: (options, handler) async {
        options.baseUrl = baseUrl;
        final token = await _storage.read(key: _tokenKey);
        if (token != null) {
          options.headers['Authorization'] = 'Bearer $token';
        }
        handler.next(options);
      },
      onError: (error, handler) async {
        if (error.response?.statusCode == 401) {
          final refreshed = await _refreshToken();
          if (refreshed) {
            final retryResponse = await _retry(error.requestOptions);
            handler.resolve(retryResponse);
            return;
          }
        }
        handler.next(error);
      },
    ));
  }

  Future<bool> _refreshToken() async {
    try {
      final refreshToken = await _storage.read(key: _refreshKey);
      if (refreshToken == null) return false;

      final response = await Dio(BaseOptions(baseUrl: baseUrl)).post(
        '/auth/refresh',
        data: {'refresh_token': refreshToken},
      );

      if (response.statusCode == 200) {
        final tokens = response.data['tokens'];
        await _storage.write(key: _tokenKey, value: tokens['access_token']);
        await _storage.write(key: _refreshKey, value: tokens['refresh_token']);
        return true;
      }
      return false;
    } catch (_) {
      return false;
    }
  }

  Future<Response> _retry(RequestOptions requestOptions) async {
    final options = Options(
      method: requestOptions.method,
      headers: requestOptions.headers,
    );
    return _dio.request(requestOptions.path, options: options);
  }

  Future<void> setTokens(String accessToken, String refreshToken) async {
    await _storage.write(key: _tokenKey, value: accessToken);
    await _storage.write(key: _refreshKey, value: refreshToken);
  }

  Future<void> clearTokens() async {
    await _storage.delete(key: _tokenKey);
    await _storage.delete(key: _refreshKey);
  }

  Future<String?> getAccessToken() => _storage.read(key: _tokenKey);
}
