CREATE TABLE emergency_events (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    room_id         UUID REFERENCES touring_rooms(id) ON DELETE SET NULL,
    reported_by     UUID NOT NULL REFERENCES users(id),
    event_type      VARCHAR(30) NOT NULL CHECK (event_type IN (
                        'accident', 'breakdown', 'lost_rider', 'medical', 'weather',
                        'road_hazard', 'crime', 'fire', 'natural_disaster', 'other'
                    )),
    severity        VARCHAR(10) NOT NULL CHECK (severity IN ('low', 'medium', 'high', 'critical')),
    location        GEOGRAPHY(POINT, 4326),
    description     TEXT,
    status          VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'acknowledged', 'resolved', 'false_alarm')),
    resolved_by     UUID REFERENCES users(id),
    resolved_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE sos_events (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id         UUID NOT NULL REFERENCES users(id),
    room_id         UUID REFERENCES touring_rooms(id) ON DELETE SET NULL,
    location        GEOGRAPHY(POINT, 4326) NOT NULL,
    triggered_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    status          VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'acknowledged', 'resolved', 'false_alarm')),
    acknowledged_by UUID REFERENCES users(id),
    acknowledged_at TIMESTAMPTZ,
    resolved_at     TIMESTAMPTZ,
    notes           TEXT
);

CREATE INDEX idx_emergency_room ON emergency_events(room_id);
CREATE INDEX idx_emergency_status ON emergency_events(status);
CREATE INDEX idx_emergency_type ON emergency_events(event_type);
CREATE INDEX idx_emergency_created ON emergency_events(created_at DESC);
CREATE INDEX idx_emergency_location ON emergency_events USING GIST(location);
CREATE INDEX idx_sos_user ON sos_events(user_id);
CREATE INDEX idx_sos_room ON sos_events(room_id);
CREATE INDEX idx_sos_status ON sos_events(status);
CREATE INDEX idx_sos_location ON sos_events USING GIST(location);
