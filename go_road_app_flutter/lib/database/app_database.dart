import 'dart:convert';
import 'dart:io';
import 'package:path/path.dart' as p;
import 'package:path_provider/path_provider.dart';

class LocalDatabase {
  Future<File> _getFile(String name) async {
    final dir = await getApplicationDocumentsDirectory();
    return File(p.join(dir.path, 'goroad_cache', '$name.json'));
  }

  Future<void> _ensureDir() async {
    final dir = await getApplicationDocumentsDirectory();
    final cacheDir = Directory(p.join(dir.path, 'goroad_cache'));
    if (!await cacheDir.exists()) {
      await cacheDir.create(recursive: true);
    }
  }

  Future<void> put(String collection, String id, Map<String, dynamic> data) async {
    await _ensureDir();
    final file = await _getFile('$collection\_$id');
    await file.writeAsString(jsonEncode(data));
  }

  Future<Map<String, dynamic>?> get(String collection, String id) async {
    final file = await _getFile('$collection\_$id');
    if (await file.exists()) {
      return jsonDecode(await file.readAsString()) as Map<String, dynamic>;
    }
    return null;
  }

  Future<List<Map<String, dynamic>>> getAll(String collection) async {
    await _ensureDir();
    final dir = await getApplicationDocumentsDirectory();
    final cacheDir = Directory(p.join(dir.path, 'goroad_cache'));
    if (!await cacheDir.exists()) return [];
    final files = await cacheDir.list().where((e) => e.path.contains(collection)).toList();
    final results = <Map<String, dynamic>>[];
    for (final f in files) {
      if (f is File) {
        results.add(jsonDecode(await f.readAsString()) as Map<String, dynamic>);
      }
    }
    return results;
  }

  Future<void> delete(String collection, String id) async {
    final file = await _getFile('$collection\_$id');
    if (await file.exists()) {
      await file.delete();
    }
  }

  Future<void> clear(String collection) async {
    await _ensureDir();
    final dir = await getApplicationDocumentsDirectory();
    final cacheDir = Directory(p.join(dir.path, 'goroad_cache'));
    if (!await cacheDir.exists()) return;
    await for (final e in cacheDir.list()) {
      if (e is File && e.path.contains(collection)) {
        await e.delete();
      }
    }
  }
}

final localDb = LocalDatabase();

Future<void> cacheRoom(Map<String, dynamic> room) async {
  await localDb.put('rooms', room['id'], room);
}

Future<List<Map<String, dynamic>>> getCachedRooms() async {
  return await localDb.getAll('rooms');
}

Future<void> cacheMessage(Map<String, dynamic> msg) async {
  await localDb.put('messages', msg['id'], msg);
}

Future<void> enqueueOffline(Map<String, dynamic> op) async {
  final id = 'offline\_${DateTime.now().millisecondsSinceEpoch}\_${op['id']}';
  await localDb.put('offline', id, op);
}

Future<List<Map<String, dynamic>>> getPendingSync() async {
  return await localDb.getAll('offline');
}

Future<void> markSynced(String opId) async {
  await localDb.delete('offline', opId);
}
