import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

class PaginatedState<T> {
  final List<T> items;
  final String? nextCursor;
  final bool hasMore;
  final bool isLoadingMore;
  final bool isLoading;
  final String? error;

  const PaginatedState({
    this.items = const [],
    this.nextCursor,
    this.hasMore = true,
    this.isLoadingMore = false,
    this.isLoading = false,
    this.error,
  });

  PaginatedState<T> copyWith({
    List<T>? items,
    String? nextCursor,
    bool? hasMore,
    bool? isLoadingMore,
    bool? isLoading,
    String? error,
  }) {
    return PaginatedState(
      items: items ?? this.items,
      nextCursor: nextCursor ?? this.nextCursor,
      hasMore: hasMore ?? this.hasMore,
      isLoadingMore: isLoadingMore ?? this.isLoadingMore,
      isLoading: isLoading ?? this.isLoading,
      error: error ?? this.error,
    );
  }
}

class PaginatedListView<T> extends ConsumerStatefulWidget {
  final AutoDisposeNotifierProvider<PaginatedNotifier<T>, PaginatedState<T>> provider;
  final Widget Function(BuildContext, T, int) itemBuilder;
  final Widget? emptyWidget;
  final Widget? loadingWidget;
  final Widget Function(String)? errorWidget;

  const PaginatedListView({
    super.key,
    required this.provider,
    required this.itemBuilder,
    this.emptyWidget,
    this.loadingWidget,
    this.errorWidget,
  });

  @override
  ConsumerState<PaginatedListView<T>> createState() => _PaginatedListViewState<T>();
}

class _PaginatedListViewState<T> extends ConsumerState<PaginatedListView<T>> {
  final _scrollController = ScrollController();

  @override
  void initState() {
    super.initState();
    _scrollController.addListener(_onScroll);
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(widget.provider.notifier).loadFirst();
    });
  }

  @override
  void dispose() {
    _scrollController.dispose();
    super.dispose();
  }

  void _onScroll() {
    if (_scrollController.position.pixels >= _scrollController.position.maxScrollExtent * 0.8) {
      ref.read(widget.provider.notifier).loadMore();
    }
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(widget.provider);

    if (state.isLoading) {
      return widget.loadingWidget ?? const Center(child: CircularProgressIndicator());
    }

    if (state.error != null && state.items.isEmpty) {
      return widget.errorWidget?.call(state.error!) ?? Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text(state.error!, style: const TextStyle(color: Colors.red)),
            const SizedBox(height: 8),
            ElevatedButton(onPressed: () => ref.read(widget.provider.notifier).refresh(), child: const Text('Retry')),
          ],
        ),
      );
    }

    if (state.items.isEmpty && !state.isLoadingMore) {
      return widget.emptyWidget ?? const Center(child: Text('No data'));
    }

    return RefreshIndicator(
      onRefresh: () => ref.read(widget.provider.notifier).refresh(),
      child: ListView.builder(
        controller: _scrollController,
        itemCount: state.items.length + (state.hasMore ? 1 : 0),
        itemBuilder: (context, index) {
          if (index >= state.items.length) {
            return const Padding(
              padding: EdgeInsets.all(16),
              child: Center(child: CircularProgressIndicator(strokeWidth: 2)),
            );
          }
          return widget.itemBuilder(context, state.items[index], index);
        },
      ),
    );
  }
}

class PaginatedNotifier<T> extends AutoDisposeNotifier<PaginatedState<T>> {
  Future<void> loadFirst() async {}
  Future<void> loadMore() async {}
  Future<void> refresh() async {
    state = state.copyWith(isLoading: true);
    await loadFirst();
  }

  @override
  PaginatedState<T> build() => const PaginatedState();
}
