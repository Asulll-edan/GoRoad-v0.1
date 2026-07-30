import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'package:go_road_app_flutter/main.dart';

void main() {
  testWidgets('App renders with ProviderScope', (WidgetTester tester) async {
    await tester.pumpWidget(const ProviderScope(child: GoRoadApp()));
    await tester.pumpAndSettle();

    expect(find.text('Go Road'), findsOneWidget);
  });
}
