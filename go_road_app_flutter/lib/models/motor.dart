class Motor {
  final String id;
  final String name;
  final String brand;
  final String? model;
  final int year;
  final String? plateNumber;
  final String? color;
  final double? engineCc;
  final String? fuelType;
  final double? tankCapacity;
  final bool isPrimary;
  final DateTime createdAt;

  Motor({
    required this.id,
    required this.name,
    required this.brand,
    this.model,
    this.year = 2024,
    this.plateNumber,
    this.color,
    this.engineCc,
    this.fuelType,
    this.tankCapacity,
    this.isPrimary = false,
    DateTime? createdAt,
  }) : createdAt = createdAt ?? DateTime.now();

  factory Motor.fromJson(Map<String, dynamic> json) => Motor(
        id: json['id'] as String,
        name: json['name'] as String,
        brand: json['brand'] as String,
        model: json['model'] as String?,
        year: json['year'] as int? ?? 2024,
        plateNumber: json['plate_number'] as String?,
        color: json['color'] as String?,
        engineCc: (json['engine_cc'] as num?)?.toDouble(),
        fuelType: json['fuel_type'] as String?,
        tankCapacity: (json['tank_capacity'] as num?)?.toDouble(),
        isPrimary: json['is_primary'] as bool? ?? false,
        createdAt: json['created_at'] != null ? DateTime.parse(json['created_at'] as String) : DateTime.now(),
      );

  Map<String, dynamic> toJson() => {
        'id': id,
        'name': name,
        'brand': brand,
        'model': model,
        'year': year,
        'plate_number': plateNumber,
        'color': color,
        'engine_cc': engineCc,
        'fuel_type': fuelType,
        'tank_capacity': tankCapacity,
        'is_primary': isPrimary,
      };
}

class FuelLog {
  final String id;
  final String motorId;
  final double liters;
  final double price;
  final double? odometer;
  final String fuelType;
  final String? station;
  final DateTime filledAt;

  FuelLog({
    required this.id,
    required this.motorId,
    required this.liters,
    required this.price,
    this.odometer,
    this.fuelType = 'pertalite',
    this.station,
    DateTime? filledAt,
  }) : filledAt = filledAt ?? DateTime.now();

  factory FuelLog.fromJson(Map<String, dynamic> json) => FuelLog(
        id: json['id'] as String,
        motorId: json['motor_id'] as String,
        liters: (json['liters'] as num).toDouble(),
        price: (json['price'] as num).toDouble(),
        odometer: (json['odometer'] as num?)?.toDouble(),
        fuelType: json['fuel_type'] as String? ?? 'pertalite',
        station: json['station'] as String?,
        filledAt: json['filled_at'] != null ? DateTime.parse(json['filled_at'] as String) : DateTime.now(),
      );
}
