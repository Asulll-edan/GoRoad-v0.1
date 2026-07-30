CREATE TABLE rider_locations (
    time            TIMESTAMPTZ NOT NULL,
    room_id         UUID NOT NULL,
    user_id         UUID NOT NULL,
    location        GEOGRAPHY(POINT, 4326) NOT NULL,
    speed_kmh       DECIMAL(5,1),
    heading         DECIMAL(5,1),
    altitude        DECIMAL(7,1),
    accuracy        DECIMAL(5,1),
    battery_level   DECIMAL(4,1),
    is_charging     BOOLEAN DEFAULT false,
    gps_status      VARCHAR(20) DEFAULT 'fixed' CHECK (gps_status IN ('none', 'fixed', 'approximate', 'dead_reckoning')),
    source          VARCHAR(20) DEFAULT 'gps' CHECK (source IN ('gps', 'manual', 'simulated'))
);

SELECT create_hypertable('rider_locations', 'time', if_not_exists => TRUE);

CREATE INDEX idx_rider_locations_room_time ON rider_locations(room_id, time DESC);
CREATE INDEX idx_rider_locations_user_time ON rider_locations(user_id, time DESC);
CREATE INDEX idx_rider_locations_geom ON rider_locations USING GIST(location);

-- Compression: after 7 days
ALTER TABLE rider_locations SET (
    timescaledb.compress,
    timescaledb.compress_segmentby = 'room_id, user_id',
    timescaledb.compress_orderby = 'time DESC'
);

SELECT add_compression_policy('rider_locations', INTERVAL '7 days', if_not_exists => TRUE);

-- Retention: 90 days
SELECT add_retention_policy('rider_locations', INTERVAL '90 days', if_not_exists => TRUE);
