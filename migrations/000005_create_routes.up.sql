CREATE TABLE routes (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    room_id         UUID REFERENCES touring_rooms(id) ON DELETE SET NULL,
    created_by      UUID NOT NULL REFERENCES users(id),
    name            VARCHAR(200) NOT NULL,
    description     TEXT,
    
    -- Route geometry
    route_geom      GEOGRAPHY(LINESTRING, 4326),
    distance_km     DECIMAL(10,2),
    estimated_duration INTERVAL,
    elevation_gain  INTEGER,
    elevation_loss  INTEGER,
    max_elevation   INTEGER,
    min_elevation   INTEGER,
    
    -- Origin / Destination
    origin_lat      DECIMAL(10,7),
    origin_lng      DECIMAL(10,7),
    origin_name     VARCHAR(200),
    dest_lat        DECIMAL(10,7),
    dest_lng        DECIMAL(10,7),
    dest_name       VARCHAR(200),
    
    -- Polyline (simplified for mobile)
    polyline        TEXT,
    
    -- Status
    is_active       BOOLEAN DEFAULT false,
    is_public       BOOLEAN DEFAULT true,
    
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE TABLE waypoints (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    route_id        UUID NOT NULL REFERENCES routes(id) ON DELETE CASCADE,
    name            VARCHAR(200),
    description     TEXT,
    location        GEOGRAPHY(POINT, 4326) NOT NULL,
    order_index     INTEGER NOT NULL,
    waypoint_type   VARCHAR(30) DEFAULT 'stop' CHECK (waypoint_type IN ('start', 'stop', 'rest', 'fuel', 'eat', 'photo', 'pit_stop', 'finish', 'custom')),
    estimated_arrival TIMESTAMPTZ,
    estimated_departure TIMESTAMPTZ,
    is_poi          BOOLEAN DEFAULT false,
    
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_routes_room ON routes(room_id);
CREATE INDEX idx_routes_created_by ON routes(created_by);
CREATE INDEX idx_routes_active ON routes(room_id, is_active) WHERE is_active = true;
CREATE INDEX idx_routes_geom ON routes USING GIST(route_geom);
CREATE INDEX idx_waypoints_route ON waypoints(route_id);
CREATE INDEX idx_waypoints_location ON waypoints USING GIST(location);
