class Room {
  final String id;
  final String name;
  final String? description;
  final String? coverUrl;
  final String? category;
  final bool isPrivate;
  final int memberCount;
  final int maxMembers;
  final String? region;
  final String? originCity;
  final String? destinationCity;
  final String status;
  final String? departureTime;
  final String createdBy;
  final String? createdByName;
  final DateTime createdAt;
  final List<RoomMember>? members;

  Room({
    required this.id,
    required this.name,
    this.description,
    this.coverUrl,
    this.category,
    this.isPrivate = false,
    this.memberCount = 0,
    this.maxMembers = 50,
    this.region,
    this.originCity,
    this.destinationCity,
    this.status = 'active',
    this.departureTime,
    required this.createdBy,
    this.createdByName,
    DateTime? createdAt,
    this.members,
  }) : createdAt = createdAt ?? DateTime.now();

  factory Room.fromJson(Map<String, dynamic> json) => Room(
        id: json['id'] as String,
        name: json['name'] as String,
        description: json['description'] as String?,
        coverUrl: json['cover_url'] as String?,
        category: json['category'] as String?,
        isPrivate: json['is_private'] as bool? ?? false,
        memberCount: json['member_count'] as int? ?? 0,
        maxMembers: json['max_members'] as int? ?? 50,
        region: json['region'] as String?,
        originCity: json['origin_city'] as String?,
        destinationCity: json['destination_city'] as String?,
        status: json['status'] as String? ?? 'active',
        departureTime: json['departure_time'] as String?,
        createdBy: json['created_by'] as String,
        createdByName: json['created_by_name'] as String?,
        createdAt: json['created_at'] != null ? DateTime.parse(json['created_at'] as String) : DateTime.now(),
        members: (json['members'] as List<dynamic>?)
            ?.map((e) => RoomMember.fromJson(e as Map<String, dynamic>))
            .toList(),
      );

  Map<String, dynamic> toJson() => {
        'id': id,
        'name': name,
        'description': description,
        'cover_url': coverUrl,
        'category': category,
        'is_private': isPrivate,
        'member_count': memberCount,
        'max_members': maxMembers,
        'region': region,
        'origin_city': originCity,
        'destination_city': destinationCity,
        'status': status,
        'departure_time': departureTime,
        'created_by': createdBy,
      };
}

class RoomMember {
  final String id;
  final String userId;
  final String? username;
  final String? avatarUrl;
  final String role;
  final String status;
  final DateTime joinedAt;

  RoomMember({
    required this.id,
    required this.userId,
    this.username,
    this.avatarUrl,
    this.role = 'member',
    this.status = 'active',
    DateTime? joinedAt,
  }) : joinedAt = joinedAt ?? DateTime.now();

  factory RoomMember.fromJson(Map<String, dynamic> json) => RoomMember(
        id: json['id'] as String,
        userId: json['user_id'] as String,
        username: json['username'] as String?,
        avatarUrl: json['avatar_url'] as String?,
        role: json['role'] as String? ?? 'member',
        status: json['status'] as String? ?? 'active',
        joinedAt: json['joined_at'] != null ? DateTime.parse(json['joined_at'] as String) : DateTime.now(),
      );
}

class CreateRoomRequest {
  final String name;
  final String? description;
  final String? category;
  final bool isPrivate;
  final int maxMembers;
  final String? originCity;
  final String? destinationCity;
  final String? departureTime;

  CreateRoomRequest({
    required this.name,
    this.description,
    this.category,
    this.isPrivate = false,
    this.maxMembers = 50,
    this.originCity,
    this.destinationCity,
    this.departureTime,
  });

  Map<String, dynamic> toJson() => {
        'name': name,
        'description': description,
        'category': category,
        'is_private': isPrivate,
        'max_members': maxMembers,
        'origin_city': originCity,
        'destination_city': destinationCity,
        'departure_time': departureTime,
      };
}
