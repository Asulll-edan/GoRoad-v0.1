import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:sentry_flutter/sentry_flutter.dart';
import 'router/app_router.dart';

Future<void> main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await SentryFlutter.init(
    (options) {
      options.dsn = const String.fromEnvironment('SENTRY_DSN', defaultValue: '');
      options.tracesSampleRate = 0.1;
      options.environment = const String.fromEnvironment('APP_ENV', defaultValue: 'development');
    },
    appRunner: () => runApp(const ProviderScope(child: GoRoadApp())),
  );
}

class GoRoadApp extends ConsumerWidget {
  const GoRoadApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final router = ref.watch(routerProvider);

    return MaterialApp.router(
      title: 'Go Road',
      debugShowCheckedModeBanner: false,
      routerConfig: router,
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(
          seedColor: const Color(0xFF1A73E8),
          brightness: Brightness.light,
        ),
        useMaterial3: true,
      ),
    );
  }
}
