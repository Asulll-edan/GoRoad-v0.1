import '../models/route_model.dart';
import '../models/pagination.dart';
import 'api_client.dart';

class RouteRepository {
  final ApiClient _client;

  RouteRepository(this._client);

  Future<CursorPage<TouringRoute>> getRoutes({String? cursor, int limit = 20}) async {
    final params = <String, dynamic>{'limit': limit};
    if (cursor != null) params['cursor'] = cursor;

    final response = await _client.dio.get('/routes', queryParameters: params);
    return CursorPage.fromJson(
      response.data as Map<String, dynamic>,
      (json) => TouringRoute.fromJson(json),
    );
  }

  Future<TouringRoute> getRoute(String id) async {
    final response = await _client.dio.get('/routes/$id');
    return TouringRoute.fromJson(response.data['data'] as Map<String, dynamic>);
  }

  Future<TouringRoute> createRoute(CreateRouteRequest request) async {
    final response = await _client.dio.post('/routes', data: request.toJson());
    return TouringRoute.fromJson(response.data['data'] as Map<String, dynamic>);
  }
}
