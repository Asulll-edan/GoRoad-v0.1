import 'package:go_router/go_router.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../providers/auth_provider.dart';
import '../screens/auth/login_screen.dart';
import '../screens/auth/register_screen.dart';
import '../screens/rooms/rooms_screen.dart';
import '../screens/rooms/room_detail_screen.dart';
import '../screens/rooms/create_room_screen.dart';
import '../screens/chat/chat_screen.dart';
import '../screens/profile/profile_screen.dart';
import '../screens/profile/edit_profile_screen.dart';
import '../screens/touring/route_list_screen.dart';
import '../screens/touring/create_route_screen.dart';
import '../screens/explore/explore_screen.dart';

final routerProvider = Provider<GoRouter>((ref) {
  final authState = ref.watch(authStateProvider);

  return GoRouter(
    initialLocation: '/rooms',
    redirect: (context, state) {
      final isAuthenticated = authState.isAuthenticated;
      final isAuthRoute = state.matchedLocation.startsWith('/auth');

      if (!isAuthenticated && !isAuthRoute) return '/auth/login';
      if (isAuthenticated && isAuthRoute) return '/rooms';
      return null;
    },
    routes: [
      GoRoute(path: '/auth/login', builder: (_, __) => const LoginScreen()),
      GoRoute(path: '/auth/register', builder: (_, __) => const RegisterScreen()),
      GoRoute(path: '/rooms', builder: (_, __) => const RoomsScreen()),
      GoRoute(path: '/rooms/create', builder: (_, __) => const CreateRoomScreen()),
      GoRoute(path: '/rooms/:id', builder: (_, state) => RoomDetailScreen(roomId: state.pathParameters['id']!)),
      GoRoute(path: '/rooms/:id/chat', builder: (_, state) => ChatScreen(roomId: state.pathParameters['id']!)),
      GoRoute(path: '/profile', builder: (_, __) => const ProfileScreen()),
      GoRoute(path: '/profile/edit', builder: (_, __) => const EditProfileScreen()),
      GoRoute(path: '/routes', builder: (_, __) => const RouteListScreen()),
      GoRoute(path: '/routes/create', builder: (_, __) => const CreateRouteScreen()),
      GoRoute(path: '/explore', builder: (_, __) => const ExploreScreen()),
    ],
  );
});
