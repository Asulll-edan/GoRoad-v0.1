class Expense {
  final String id;
  final String roomId;
  final String? userId;
  final String? username;
  final double amount;
  final String category;
  final String? description;
  final DateTime createdAt;

  Expense({
    required this.id,
    required this.roomId,
    this.userId,
    this.username,
    this.amount = 0,
    this.category = 'other',
    this.description,
    DateTime? createdAt,
  }) : createdAt = createdAt ?? DateTime.now();

  factory Expense.fromJson(Map<String, dynamic> json) => Expense(
        id: json['id'] as String,
        roomId: json['room_id'] as String,
        userId: json['user_id'] as String?,
        username: json['username'] as String?,
        amount: (json['amount'] as num).toDouble(),
        category: json['category'] as String? ?? 'other',
        description: json['description'] as String?,
        createdAt: json['created_at'] != null ? DateTime.parse(json['created_at'] as String) : DateTime.now(),
      );

  Map<String, dynamic> toJson() => {
        'id': id,
        'room_id': roomId,
        'amount': amount,
        'category': category,
        'description': description,
      };
}

class Poi {
  final String id;
  final String name;
  final String? category;
  final double lat;
  final double lon;
  final String? address;
  final String? phone;
  final double? rating;
  final String? photoUrl;

  Poi({
    required this.id,
    required this.name,
    this.category,
    required this.lat,
    required this.lon,
    this.address,
    this.phone,
    this.rating,
    this.photoUrl,
  });

  factory Poi.fromJson(Map<String, dynamic> json) => Poi(
        id: json['id'] as String,
        name: json['name'] as String,
        category: json['category'] as String?,
        lat: (json['lat'] as num).toDouble(),
        lon: (json['lon'] as num).toDouble(),
        address: json['address'] as String?,
        phone: json['phone'] as String?,
        rating: (json['rating'] as num?)?.toDouble(),
        photoUrl: json['photo_url'] as String?,
      );
}
