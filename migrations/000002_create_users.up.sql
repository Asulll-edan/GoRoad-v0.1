CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email           VARCHAR(255) NOT NULL UNIQUE,
    phone           VARCHAR(20) UNIQUE,
    username        VARCHAR(50) NOT NULL UNIQUE,
    full_name       VARCHAR(100) NOT NULL,
    password_hash   VARCHAR(255) NOT NULL,
    photo_url       VARCHAR(500),
    bio             TEXT,
    date_of_birth   DATE,
    gender          VARCHAR(10) CHECK (gender IN ('male', 'female', 'other')),
    
    -- Rider profile
    riding_skill    VARCHAR(20) DEFAULT 'beginner' CHECK (riding_skill IN ('beginner', 'intermediate', 'advanced', 'expert')),
    riding_role     VARCHAR(20) DEFAULT 'member' CHECK (riding_role IN ('lead', 'sweep', 'member', 'marshal')),
    total_points    BIGINT DEFAULT 0,
    experience_years DECIMAL(3,1) DEFAULT 0,
    emergency_contact_name VARCHAR(100),
    emergency_contact_phone VARCHAR(20),
    medical_notes TEXT,
    blood_type VARCHAR(5),
    
    -- Preferences
    preferences     JSONB DEFAULT '{}',
    privacy_settings JSONB DEFAULT '{"show_location": true, "show_phone": false, "show_email": false}',
    
    -- Status
    is_verified     BOOLEAN DEFAULT false,
    is_active       BOOLEAN DEFAULT true,
    is_banned       BOOLEAN DEFAULT false,
    last_active_at  TIMESTAMPTZ,
    
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE TABLE device_tokens (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token           VARCHAR(500) NOT NULL,
    platform        VARCHAR(20) NOT NULL CHECK (platform IN ('ios', 'android', 'web')),
    is_active       BOOLEAN DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE refresh_tokens (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash      VARCHAR(255) NOT NULL UNIQUE,
    expires_at      TIMESTAMPTZ NOT NULL,
    is_revoked      BOOLEAN DEFAULT false,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at      TIMESTAMPTZ
);

-- Indexes
CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_phone ON users(phone) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_total_points ON users(total_points DESC);
CREATE INDEX idx_users_riding_skill ON users(riding_skill);
CREATE INDEX idx_users_active ON users(is_active) WHERE is_active = true;
CREATE INDEX idx_device_tokens_user ON device_tokens(user_id);
CREATE INDEX idx_refresh_tokens_hash ON refresh_tokens(token_hash);
CREATE INDEX idx_refresh_tokens_user ON refresh_tokens(user_id);
