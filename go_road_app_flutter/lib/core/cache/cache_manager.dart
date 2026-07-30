import 'dart:collection';
import 'package:flutter/foundation.dart';

class CacheEntry<T> {
  final T data;
  final DateTime cachedAt;
  final Duration ttl;

  CacheEntry(this.data, {this.ttl = const Duration(minutes: 5)})
      : cachedAt = DateTime.now();

  bool get isExpired => DateTime.now().difference(cachedAt) > ttl;

  Duration get age => DateTime.now().difference(cachedAt);
}

class CacheManager<T> {
  final LinkedHashMap<String, CacheEntry<T>> _cache = LinkedHashMap();
  final int maxEntries;
  final Duration defaultTtl;

  CacheManager({this.maxEntries = 100, this.defaultTtl = const Duration(minutes: 5)});

  T? get(String key) {
    final entry = _cache[key];
    if (entry == null) return null;
    if (entry.isExpired) {
      _cache.remove(key);
      return null;
    }
    return entry.data;
  }

  void set(String key, T data, {Duration? ttl}) {
    if (_cache.length >= maxEntries) {
      _cache.remove(_cache.keys.first);
    }
    _cache[key] = CacheEntry(data, ttl: ttl ?? defaultTtl);
  }

  void remove(String key) => _cache.remove(key);

  void clear() => _cache.clear();

  bool has(String key) => _cache.containsKey(key) && !_cache[key]!.isExpired;

  int get size => _cache.length;

  List<MapEntry<String, T>> getEntries() {
    return _cache.entries
        .where((e) => !e.value.isExpired)
        .map((e) => MapEntry(e.key, e.value.data))
        .toList();
  }
}
