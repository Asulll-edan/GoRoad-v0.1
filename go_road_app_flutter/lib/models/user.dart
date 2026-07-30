class User {
  final String id;
  final String username;
  final String email;
  final String? displayName;
  final String? avatarUrl;
  final String? bio;
  final String? phone;
  final String role;
  final bool isVerified;
  final String? fcmToken;
  final String? motorName;
  final String? motorPlate;
  final DateTime createdAt;

  User({
    required this.id,
    required this.username,
    required this.email,
    this.displayName,
    this.avatarUrl,
    this.bio,
    this.phone,
    this.role = 'user',
    this.isVerified = false,
    this.fcmToken,
    this.motorName,
    this.motorPlate,
    DateTime? createdAt,
  }) : createdAt = createdAt ?? DateTime.now();

  factory User.fromJson(Map<String, dynamic> json) => User(
        id: json['id'] as String,
        username: json['username'] as String,
        email: json['email'] as String,
        displayName: json['display_name'] as String?,
        avatarUrl: json['avatar_url'] as String?,
        bio: json['bio'] as String?,
        phone: json['phone'] as String?,
        role: json['role'] as String? ?? 'user',
        isVerified: json['is_verified'] as bool? ?? false,
        fcmToken: json['fcm_token'] as String?,
        motorName: json['motor_name'] as String?,
        motorPlate: json['motor_plate'] as String?,
        createdAt: json['created_at'] != null ? DateTime.parse(json['created_at'] as String) : null,
      );

  Map<String, dynamic> toJson() => {
        'id': id,
        'username': username,
        'email': email,
        'display_name': displayName,
        'avatar_url': avatarUrl,
        'bio': bio,
        'phone': phone,
        'role': role,
        'is_verified': isVerified,
        'fcm_token': fcmToken,
        'motor_name': motorName,
        'motor_plate': motorPlate,
      };
}

class AuthTokens {
  final String accessToken;
  final String refreshToken;

  AuthTokens({required this.accessToken, required this.refreshToken});

  factory AuthTokens.fromJson(Map<String, dynamic> json) => AuthTokens(
        accessToken: json['access_token'] as String,
        refreshToken: json['refresh_token'] as String,
      );

  Map<String, dynamic> toJson() => {
        'access_token': accessToken,
        'refresh_token': refreshToken,
      };
}

class LoginRequest {
  final String email;
  final String password;

  LoginRequest({required this.email, required this.password});

  Map<String, dynamic> toJson() => {'email': email, 'password': password};
}

class RegisterRequest {
  final String username;
  final String email;
  final String password;
  final String? displayName;
  final String? phone;

  RegisterRequest({
    required this.username,
    required this.email,
    required this.password,
    this.displayName,
    this.phone,
  });

  Map<String, dynamic> toJson() => {
        'username': username,
        'email': email,
        'password': password,
        'display_name': displayName,
        'phone': phone,
      };
}

class LoginResponse {
  final User user;
  final AuthTokens tokens;

  LoginResponse({required this.user, required this.tokens});

  factory LoginResponse.fromJson(Map<String, dynamic> json) => LoginResponse(
        user: User.fromJson(json['user'] as Map<String, dynamic>),
        tokens: AuthTokens.fromJson(json['tokens'] as Map<String, dynamic>),
      );
}
