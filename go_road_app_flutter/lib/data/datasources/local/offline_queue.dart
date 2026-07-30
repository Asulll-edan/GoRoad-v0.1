import 'dart:convert';
import 'package:flutter/foundation.dart';

enum QueuePriority { sos, emergency, location, chat, others }

class OfflineOperation {
  final String id;
  final QueuePriority priority;
  final String entityType;
  final String entityId;
  final String action;
  final Map<String, dynamic> data;
  final DateTime createdAt;
  bool isSynced;

  OfflineOperation({
    required this.id,
    this.priority = QueuePriority.others,
    required this.entityType,
    required this.entityId,
    required this.action,
    required this.data,
    DateTime? createdAt,
    this.isSynced = false,
  }) : createdAt = createdAt ?? DateTime.now();

  Map<String, dynamic> toJson() => {
    'id': id,
    'priority': priority.index,
    'entity_type': entityType,
    'entity_id': entityId,
    'action': action,
    'data': data,
    'created_at': createdAt.toIso8601String(),
    'is_synced': isSynced,
  };

  factory OfflineOperation.fromJson(Map<String, dynamic> json) => OfflineOperation(
    id: json['id'] as String,
    priority: QueuePriority.values[json['priority'] as int],
    entityType: json['entity_type'] as String,
    entityId: json['entity_id'] as String,
    action: json['action'] as String,
    data: Map<String, dynamic>.from(json['data'] as Map),
    createdAt: DateTime.parse(json['created_at'] as String),
    isSynced: json['is_synced'] as bool? ?? false,
  );
}

class OfflineQueueManager {
  final List<OfflineOperation> _queue = [];
  final ValueNotifier<int> pendingCount = ValueNotifier(0);

  void enqueue(OfflineOperation operation) {
    _queue.add(operation);
    _queue.sort((a, b) => a.priority.index.compareTo(b.priority.index));
    pendingCount.value = _queue.length;
  }

  OfflineOperation? dequeue() {
    if (_queue.isEmpty) return null;
    final op = _queue.removeAt(0);
    pendingCount.value = _queue.length;
    return op;
  }

  List<OfflineOperation> get pending => List.unmodifiable(_queue);

  void markSynced(String id) {
    final idx = _queue.indexWhere((op) => op.id == id);
    if (idx != -1) {
      _queue[idx].isSynced = true;
      _queue.removeAt(idx);
      pendingCount.value = _queue.length;
    }
  }

  int get count => _queue.length;

  void clear() {
    _queue.clear();
    pendingCount.value = 0;
  }

  Future<void> processQueue(Future<bool> Function(OfflineOperation op) syncFn) async {
    final sorted = List<OfflineOperation>.from(_queue)
      ..sort((a, b) => a.priority.index.compareTo(b.priority.index));
    for (final op in sorted) {
      final success = await syncFn(op);
      if (success) {
        markSynced(op.id);
      }
    }
  }
}
