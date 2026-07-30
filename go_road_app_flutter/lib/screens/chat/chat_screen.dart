import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../providers/chat_provider.dart';

class ChatScreen extends ConsumerStatefulWidget {
  final String roomId;

  const ChatScreen({super.key, required this.roomId});

  @override
  ConsumerState<ChatScreen> createState() => _ChatScreenState();
}

class _ChatScreenState extends ConsumerState<ChatScreen> {
  final _messageController = TextEditingController();
  final _scrollController = ScrollController();

  @override
  void initState() {
    super.initState();
    Future.microtask(() => ref.read(messagesProvider(widget.roomId).notifier).loadMessages());
  }

  @override
  void dispose() {
    _messageController.dispose();
    _scrollController.dispose();
    super.dispose();
  }

  void _sendMessage() {
    final text = _messageController.text.trim();
    if (text.isEmpty) return;
    ref.read(messagesProvider(widget.roomId).notifier).sendMessage(text);
    _messageController.clear();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _scrollController.animateTo(
        _scrollController.position.maxScrollExtent,
        duration: const Duration(milliseconds: 300),
        curve: Curves.easeOut,
      );
    });
  }

  @override
  Widget build(BuildContext context) {
    final messagesState = ref.watch(messagesProvider(widget.roomId));
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(title: const Text('Chat')),
      body: Column(
        children: [
          Expanded(
            child: messagesState.isLoading && messagesState.messages.isEmpty
                ? const Center(child: CircularProgressIndicator())
                : messagesState.messages.isEmpty
                    ? const Center(child: Text('Belum ada pesan'))
                    : ListView.builder(
                        controller: _scrollController,
                        padding: const EdgeInsets.all(16),
                        itemCount: messagesState.messages.length,
                        itemBuilder: (context, index) {
                          final msg = messagesState.messages[index];
                          final isMe = msg.userId == 'me';
                          return Align(
                            alignment: isMe ? Alignment.centerRight : Alignment.centerLeft,
                            child: Container(
                              margin: const EdgeInsets.only(bottom: 8),
                              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 10),
                              decoration: BoxDecoration(
                                color: isMe ? theme.colorScheme.primary : Colors.grey[200],
                                borderRadius: BorderRadius.only(
                                  topLeft: const Radius.circular(16),
                                  topRight: const Radius.circular(16),
                                  bottomLeft: Radius.circular(isMe ? 16 : 4),
                                  bottomRight: Radius.circular(isMe ? 4 : 16),
                                ),
                              ),
                              constraints: BoxConstraints(maxWidth: MediaQuery.of(context).size.width * 0.75),
                              child: Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  if (!isMe && msg.username != null)
                                    Text(msg.username!, style: const TextStyle(fontSize: 12, fontWeight: FontWeight.bold, color: Colors.black54)),
                                  const SizedBox(height: 4),
                                  Text(
                                    msg.message,
                                    style: TextStyle(color: isMe ? Colors.white : Colors.black87),
                                  ),
                                  const SizedBox(height: 4),
                                  Text(
                                    _formatTime(msg.createdAt),
                                    style: TextStyle(fontSize: 10, color: isMe ? Colors.white60 : Colors.grey),
                                  ),
                                ],
                              ),
                            ),
                          );
                        },
                      ),
          ),
          Container(
            padding: const EdgeInsets.all(16),
            decoration: BoxDecoration(
              color: theme.colorScheme.surface,
              boxShadow: [BoxShadow(color: Colors.black12, blurRadius: 4, offset: const Offset(0, -2))],
            ),
            child: SafeArea(
              child: Row(
                children: [
                  Expanded(
                    child: TextField(
                      controller: _messageController,
                      decoration: const InputDecoration(
                        hintText: 'Ketik pesan...',
                        border: OutlineInputBorder(),
                        contentPadding: EdgeInsets.symmetric(horizontal: 16, vertical: 12),
                      ),
                      onSubmitted: (_) => _sendMessage(),
                    ),
                  ),
                  const SizedBox(width: 8),
                  IconButton.filled(
                    onPressed: _sendMessage,
                    icon: const Icon(Icons.send),
                  ),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }

  String _formatTime(DateTime dt) {
    return '${dt.hour.toString().padLeft(2, '0')}:${dt.minute.toString().padLeft(2, '0')}';
  }
}
