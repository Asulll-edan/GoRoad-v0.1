import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class ApiClient {
  late final Dio _dio;
  final FlutterSecureStorage _storage = const FlutterSecureStorage();

  ApiClient({required String baseUrl, String? token}) {
    _dio = Dio(BaseOptions(
      baseUrl: baseUrl,
      connectTimeout: const Duration(seconds: 15),
      receiveTimeout: const Duration(seconds: 15),
      headers: {'Content-Type': 'application/json'},
    ));
    if (token != null) {
      _dio.options.headers['Authorization'] = 'Bearer $token';
    }
    _dio.interceptors.addAll([
      _AuthInterceptor(_storage),
      _RetryInterceptor(),
    ]);
  }

  Future<Response> get(String path, {Map<String, dynamic>? queryParameters}) =>
      _dio.get(path, queryParameters: queryParameters);

  Future<Response> post(String path, {dynamic data}) =>
      _dio.post(path, data: data);

  Future<Response> put(String path, {dynamic data}) =>
      _dio.put(path, data: data);

  Future<Response> delete(String path) => _dio.delete(path);

  Future<Response> patch(String path, {dynamic data}) =>
      _dio.patch(path, data: data);
}

class _AuthInterceptor extends Interceptor {
  final FlutterSecureStorage _storage;
  _AuthInterceptor(this._storage);

  @override
  void onRequest(RequestOptions options, RequestInterceptorHandler handler) async {
    final token = await _storage.read(key: 'access_token');
    if (token != null) {
      options.headers['Authorization'] = 'Bearer $token';
    }
    handler.next(options);
  }

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) async {
    if (err.response?.statusCode == 401) {
      final refreshToken = await _storage.read(key: 'refresh_token');
      if (refreshToken != null) {
        try {
          final dio = Dio(BaseOptions(baseUrl: err.requestOptions.baseUrl));
          final resp = await dio.post('/v1/auth/refresh', data: {'refresh_token': refreshToken});
          if (resp.statusCode == 200) {
            final newToken = resp.data['tokens']['access_token'] as String;
            await _storage.write(key: 'access_token', value: newToken);
            err.requestOptions.headers['Authorization'] = 'Bearer $newToken';
            final retry = await dio.fetch(err.requestOptions);
            handler.resolve(retry);
            return;
          }
        } catch (_) {}
      }
    }
    handler.next(err);
  }
}

class _RetryInterceptor extends Interceptor {
  @override
  void onError(DioException err, ErrorInterceptorHandler handler) async {
    if (err.type == DioExceptionType.connectionTimeout || err.type == DioExceptionType.receiveTimeout) {
      await Future.delayed(const Duration(seconds: 1));
    }
    handler.next(err);
  }
}
