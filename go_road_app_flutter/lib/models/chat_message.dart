class ChatMessage {
  final String id;
  final String roomId;
  final String userId;
  final String? username;
  final String? avatarUrl;
  final String message;
  final String messageType;
  final DateTime createdAt;

  ChatMessage({
    required this.id,
    required this.roomId,
    required this.userId,
    this.username,
    this.avatarUrl,
    required this.message,
    this.messageType = 'text',
    DateTime? createdAt,
  }) : createdAt = createdAt ?? DateTime.now();

  factory ChatMessage.fromJson(Map<String, dynamic> json) => ChatMessage(
        id: json['id'] as String,
        roomId: json['room_id'] as String,
        userId: json['user_id'] as String,
        username: json['username'] as String?,
        avatarUrl: json['avatar_url'] as String?,
        message: json['message'] as String,
        messageType: json['message_type'] as String? ?? 'text',
        createdAt: json['created_at'] != null ? DateTime.parse(json['created_at'] as String) : DateTime.now(),
      );

  Map<String, dynamic> toJson() => {
        'id': id,
        'room_id': roomId,
        'user_id': userId,
        'message': message,
        'message_type': messageType,
      };
}
