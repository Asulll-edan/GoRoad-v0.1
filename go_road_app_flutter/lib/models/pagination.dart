class CursorPage<T> {
  final List<T> items;
  final String? nextCursor;
  final bool hasMore;

  CursorPage({
    required this.items,
    this.nextCursor,
    this.hasMore = false,
  });

  factory CursorPage.fromJson(
    Map<String, dynamic> json,
    T Function(Map<String, dynamic>) fromItem,
  ) {
    final data = json['data'] as List<dynamic>? ?? [];
    return CursorPage(
      items: data.map((e) => fromItem(e as Map<String, dynamic>)).toList(),
      nextCursor: json['next_cursor'] as String?,
      hasMore: json['has_more'] as bool? ?? false,
    );
  }
}
