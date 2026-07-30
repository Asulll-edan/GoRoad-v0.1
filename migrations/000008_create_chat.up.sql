CREATE TABLE chat_messages (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    room_id         UUID NOT NULL REFERENCES touring_rooms(id) ON DELETE CASCADE,
    sender_id       UUID NOT NULL REFERENCES users(id),
    message_type    VARCHAR(20) NOT NULL DEFAULT 'text' CHECK (message_type IN (
                        'text', 'image', 'file', 'location', 'system',
                        'alert', 'emergency', 'voting', 'route'
                    )),
    content         TEXT NOT NULL,
    reply_to_id     UUID REFERENCES chat_messages(id),
    is_pinned       BOOLEAN DEFAULT false,
    is_deleted      BOOLEAN DEFAULT false,
    metadata        JSONB DEFAULT '{}',
    sent_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    edited_at       TIMESTAMPTZ
);

CREATE TABLE message_reads (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    message_id      UUID NOT NULL REFERENCES chat_messages(id) ON DELETE CASCADE,
    user_id         UUID NOT NULL REFERENCES users(id),
    read_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(message_id, user_id)
);

CREATE INDEX idx_chat_messages_room ON chat_messages(room_id, sent_at DESC);
CREATE INDEX idx_chat_messages_sender ON chat_messages(sender_id);
CREATE INDEX idx_chat_messages_type ON chat_messages(room_id, message_type);
CREATE INDEX idx_chat_messages_pinned ON chat_messages(room_id, is_pinned) WHERE is_pinned = true;
CREATE INDEX idx_message_reads_message ON message_reads(message_id);
CREATE INDEX idx_message_reads_user ON message_reads(user_id);
