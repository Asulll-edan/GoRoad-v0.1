CREATE TABLE votings (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    room_id         UUID NOT NULL REFERENCES touring_rooms(id) ON DELETE CASCADE,
    created_by      UUID NOT NULL REFERENCES users(id),
    title           VARCHAR(200) NOT NULL,
    description     TEXT,
    voting_type     VARCHAR(20) NOT NULL DEFAULT 'single' CHECK (voting_type IN ('single', 'multiple', 'ranked')),
    status          VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'closed', 'cancelled')),
    starts_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ends_at         TIMESTAMPTZ,
    is_anonymous    BOOLEAN DEFAULT false,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE voting_answers (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    voting_id       UUID NOT NULL REFERENCES votings(id) ON DELETE CASCADE,
    label           VARCHAR(200) NOT NULL,
    order_index     INTEGER NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE voting_votes (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    voting_id       UUID NOT NULL REFERENCES votings(id) ON DELETE CASCADE,
    answer_id       UUID NOT NULL REFERENCES voting_answers(id) ON DELETE CASCADE,
    user_id         UUID NOT NULL REFERENCES users(id),
    rank            INTEGER DEFAULT 1,
    voted_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(voting_id, answer_id, user_id)
);

CREATE INDEX idx_votings_room ON votings(room_id);
CREATE INDEX idx_votings_status ON votings(status);
CREATE INDEX idx_votings_ends ON votings(ends_at) WHERE status = 'active';
CREATE INDEX idx_voting_answers_voting ON voting_answers(voting_id);
CREATE INDEX idx_voting_votes_voting ON voting_votes(voting_id);
