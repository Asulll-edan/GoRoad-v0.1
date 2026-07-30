class TouringRoute {
  final String id;
  final String name;
  final String? description;
  final double distanceKm;
  final double? durationHours;
  final double? elevationGain;
  final String? polyline;
  final String? originLat;
  final String? originLon;
  final String? destLat;
  final String? destLon;
  final String status;
  final int? rating;
  final DateTime createdAt;

  TouringRoute({
    required this.id,
    required this.name,
    this.description,
    this.distanceKm = 0,
    this.durationHours,
    this.elevationGain,
    this.polyline,
    this.originLat,
    this.originLon,
    this.destLat,
    this.destLon,
    this.status = 'draft',
    this.rating,
    DateTime? createdAt,
  }) : createdAt = createdAt ?? DateTime.now();

  factory TouringRoute.fromJson(Map<String, dynamic> json) => TouringRoute(
        id: json['id'] as String,
        name: json['name'] as String,
        description: json['description'] as String?,
        distanceKm: (json['distance_km'] as num?)?.toDouble() ?? 0,
        durationHours: (json['duration_hours'] as num?)?.toDouble(),
        elevationGain: (json['elevation_gain'] as num?)?.toDouble(),
        polyline: json['polyline'] as String?,
        originLat: json['origin_lat'] as String?,
        originLon: json['origin_lon'] as String?,
        destLat: json['dest_lat'] as String?,
        destLon: json['dest_lon'] as String?,
        status: json['status'] as String? ?? 'draft',
        rating: json['rating'] as int?,
        createdAt: json['created_at'] != null ? DateTime.parse(json['created_at'] as String) : DateTime.now(),
      );

  Map<String, dynamic> toJson() => {
        'id': id,
        'name': name,
        'description': description,
        'distance_km': distanceKm,
        'duration_hours': durationHours,
        'elevation_gain': elevationGain,
        'polyline': polyline,
        'origin_lat': originLat,
        'origin_lon': originLon,
        'dest_lat': destLat,
        'dest_lon': destLon,
        'status': status,
        'rating': rating,
      };
}

class CreateRouteRequest {
  final String name;
  final String? description;
  final double distanceKm;
  final String? polyline;
  final String? originLat;
  final String? originLon;
  final String? destLat;
  final String? destLon;

  CreateRouteRequest({
    required this.name,
    this.description,
    this.distanceKm = 0,
    this.polyline,
    this.originLat,
    this.originLon,
    this.destLat,
    this.destLon,
  });

  Map<String, dynamic> toJson() => {
        'name': name,
        'description': description,
        'distance_km': distanceKm,
        'polyline': polyline,
        'origin_lat': originLat,
        'origin_lon': originLon,
        'dest_lat': destLat,
        'dest_lon': destLon,
      };
}
