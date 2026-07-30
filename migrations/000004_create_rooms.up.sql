CREATE TABLE touring_rooms (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name            VARCHAR(100) NOT NULL,
    description     TEXT,
    cover_photo_url VARCHAR(500),
    
    -- Status machine
    status          VARCHAR(20) NOT NULL DEFAULT 'planning'
                    CHECK (status IN ('planning', 'ready', 'touring', 'paused', 'completed', 'cancelled')),
    
    -- Dates
    start_date      TIMESTAMPTZ,
    end_date        TIMESTAMPTZ,
    gathering_point GEOGRAPHY(POINT, 4326),
    gathering_address TEXT,
    gathering_time  TIMESTAMPTZ,
    
    -- Location names
    start_location  VARCHAR(200),
    end_location    VARCHAR(200),
    route_names     TEXT[],
    
    -- Settings
    max_members     INTEGER DEFAULT 20,
    is_public       BOOLEAN DEFAULT true,
    requires_approval BOOLEAN DEFAULT false,
    allow_guests    BOOLEAN DEFAULT false,
    
    -- Touring type
    touring_type    VARCHAR(30) DEFAULT 'fun_tour' CHECK (touring_type IN ('fun_tour', 'long_distance', 'adventure', 'competition', 'charity', 'education', 'other')),
    difficulty      VARCHAR(20) DEFAULT 'easy' CHECK (difficulty IN ('easy', 'moderate', 'hard', 'extreme')),
    
    -- Routes (simplified, full in routes table)
    distance_km     DECIMAL(10,2),
    estimated_duration INTERVAL,
    
    created_by      UUID NOT NULL REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE TABLE room_members (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    room_id         UUID NOT NULL REFERENCES touring_rooms(id) ON DELETE CASCADE,
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role            VARCHAR(20) NOT NULL DEFAULT 'member'
                    CHECK (role IN ('owner', 'co_owner', 'lead', 'sweep', 'marshal', 'member', 'guest')),
    position_in_formation INTEGER DEFAULT 0,
    joined_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_online       BOOLEAN DEFAULT false,
    UNIQUE(room_id, user_id)
);

CREATE TABLE room_role_history (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    room_id         UUID NOT NULL REFERENCES touring_rooms(id) ON DELETE CASCADE,
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    old_role        VARCHAR(20),
    new_role        VARCHAR(20) NOT NULL,
    changed_by      UUID NOT NULL REFERENCES users(id),
    reason          TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Spatial indexes
CREATE INDEX idx_rooms_status ON touring_rooms(status);
CREATE INDEX idx_rooms_start_date ON touring_rooms(start_date);
CREATE INDEX idx_rooms_public ON touring_rooms(is_public) WHERE is_public = true;
CREATE INDEX idx_rooms_created_by ON touring_rooms(created_by);
CREATE INDEX idx_rooms_gathering ON touring_rooms USING GIST(gathering_point);
CREATE INDEX idx_room_members_room ON room_members(room_id);
CREATE INDEX idx_room_members_user ON room_members(user_id);
CREATE INDEX idx_room_members_role ON room_members(room_id, role);
CREATE INDEX idx_room_role_history_room ON room_role_history(room_id);
