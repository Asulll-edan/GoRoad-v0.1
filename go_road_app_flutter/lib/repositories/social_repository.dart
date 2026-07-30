import '../models/pagination.dart';
import 'api_client.dart';

class SocialPost {
  final String id;
  final String userId;
  final String? username;
  final String? avatarUrl;
  final String content;
  final List<String>? images;
  final int likeCount;
  final int commentCount;
  final bool isLiked;
  final String? roomId;
  final DateTime createdAt;

  SocialPost({
    required this.id,
    required this.userId,
    this.username,
    this.avatarUrl,
    required this.content,
    this.images,
    this.likeCount = 0,
    this.commentCount = 0,
    this.isLiked = false,
    this.roomId,
    DateTime? createdAt,
  }) : createdAt = createdAt ?? DateTime.now();

  factory SocialPost.fromJson(Map<String, dynamic> json) => SocialPost(
        id: json['id'] as String,
        userId: json['user_id'] as String,
        username: json['username'] as String?,
        avatarUrl: json['avatar_url'] as String?,
        content: json['content'] as String,
        images: (json['images'] as List<dynamic>?)?.cast<String>(),
        likeCount: json['like_count'] as int? ?? 0,
        commentCount: json['comment_count'] as int? ?? 0,
        isLiked: json['is_liked'] as bool? ?? false,
        roomId: json['room_id'] as String?,
        createdAt: json['created_at'] != null ? DateTime.parse(json['created_at'] as String) : DateTime.now(),
      );
}

class SocialRepository {
  final ApiClient _client;

  SocialRepository(this._client);

  Future<CursorPage<SocialPost>> getFeed({String? cursor, int limit = 20}) async {
    final params = <String, dynamic>{'limit': limit};
    if (cursor != null) params['cursor'] = cursor;

    final response = await _client.dio.get('/social/feed', queryParameters: params);
    return CursorPage.fromJson(
      response.data as Map<String, dynamic>,
      (json) => SocialPost.fromJson(json),
    );
  }

  Future<SocialPost> createPost(String content, {List<String>? images, String? roomId}) async {
    final response = await _client.dio.post('/social/posts', data: {
      'content': content,
      'images': images,
      'room_id': roomId,
    });
    return SocialPost.fromJson(response.data['data'] as Map<String, dynamic>);
  }

  Future<void> likePost(String postId) async {
    await _client.dio.post('/social/posts/$postId/like');
  }

  Future<void> unlikePost(String postId) async {
    await _client.dio.delete('/social/posts/$postId/like');
  }
}
