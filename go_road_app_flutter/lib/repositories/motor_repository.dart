import '../models/motor.dart';
import 'api_client.dart';

class MotorRepository {
  final ApiClient _client;

  MotorRepository(this._client);

  Future<List<Motor>> getMotors() async {
    final response = await _client.dio.get('/motors');
    return (response.data['data'] as List<dynamic>)
        .map((e) => Motor.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  Future<Motor> getMotor(String id) async {
    final response = await _client.dio.get('/motors/$id');
    return Motor.fromJson(response.data['data'] as Map<String, dynamic>);
  }

  Future<Motor> createMotor(Map<String, dynamic> data) async {
    final response = await _client.dio.post('/motors', data: data);
    return Motor.fromJson(response.data['data'] as Map<String, dynamic>);
  }

  Future<Motor> updateMotor(String id, Map<String, dynamic> data) async {
    final response = await _client.dio.put('/motors/$id', data: data);
    return Motor.fromJson(response.data['data'] as Map<String, dynamic>);
  }

  Future<void> deleteMotor(String id) async {
    await _client.dio.delete('/motors/$id');
  }
}
