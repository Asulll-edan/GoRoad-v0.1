CREATE TABLE motors (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    brand           VARCHAR(50) NOT NULL,
    model           VARCHAR(100) NOT NULL,
    year            INTEGER NOT NULL,
    license_plate   VARCHAR(20),
    
    -- Encrypted fields (pgcrypto)
    vin             BYTEA,
    insurance_info  BYTEA,
    stnk_number     BYTEA,
    
    -- Specifications
    engine_cc       INTEGER,
    fuel_type       VARCHAR(20) CHECK (fuel_type IN ('bensin', 'solar', 'electric', 'hybrid')),
    tank_capacity   DECIMAL(5,1),
    fuel_consumption DECIMAL(5,2),
    tire_pressure_front DECIMAL(3,1),
    tire_pressure_rear DECIMAL(3,1),
    
    -- Media
    photo_url       VARCHAR(500),
    
    -- Status
    is_primary      BOOLEAN DEFAULT false,
    is_active       BOOLEAN DEFAULT true,
    
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_motors_user ON motors(user_id);
CREATE INDEX idx_motors_primary ON motors(user_id, is_primary) WHERE is_primary = true;
