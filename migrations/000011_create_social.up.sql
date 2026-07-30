-- ============================================
-- BADGES
-- ============================================
CREATE TABLE badges (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code            VARCHAR(50) NOT NULL UNIQUE,
    name            VARCHAR(100) NOT NULL,
    description     TEXT,
    icon_url        VARCHAR(500),
    category        VARCHAR(30) NOT NULL CHECK (category IN (
                        'distance', 'frequency', 'elevation', 'speed',
                        'endurance', 'social', 'achievement', 'special'
                    )),
    tier            VARCHAR(10) NOT NULL CHECK (tier IN ('bronze', 'silver', 'gold', 'platinum', 'diamond')),
    criteria        JSONB NOT NULL,
    is_hidden       BOOLEAN DEFAULT false,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Seed: 14 badges
INSERT INTO badges (code, name, description, category, tier, criteria) VALUES
    ('first_tour', 'First Tour', 'Complete your first touring', 'achievement', 'bronze',
     '{"type": "touring_count", "operator": "gte", "value": 1}'),
    ('road_warrior', 'Road Warrior', 'Complete 10 tourings', 'frequency', 'silver',
     '{"type": "touring_count", "operator": "gte", "value": 10}'),
    ('century_rider', 'Century Rider', 'Ride 100 km in a single touring', 'distance', 'bronze',
     '{"type": "single_distance_km", "operator": "gte", "value": 100}'),
    ('iron_butt', 'Iron Butt', 'Ride 500 km in a single touring', 'distance', 'gold',
     '{"type": "single_distance_km", "operator": "gte", "value": 500}'),
    ('thousand_milestone', 'Thousand Milestone', 'Accumulate 1000 km total', 'distance', 'silver',
     '{"type": "total_distance_km", "operator": "gte", "value": 1000}'),
    ('eagle_eye', 'Eagle Eye', 'Climb 1000m elevation in a single touring', 'elevation', 'silver',
     '{"type": "single_elevation_m", "operator": "gte", "value": 1000}'),
    ('speed_demon', 'Speed Demon', 'Maintain average speed > 80 km/h in a touring', 'speed', 'bronze',
     '{"type": "avg_speed_kmh", "operator": "gte", "value": 80}'),
    ('night_owl', 'Night Owl', 'Complete a night touring (8 PM - 5 AM)', 'endurance', 'silver',
     '{"type": "night_touring", "operator": "eq", "value": true}'),
    ('rain_rider', 'Rain Rider', 'Complete a touring in heavy rain', 'endurance', 'gold',
     '{"type": "rain_touring", "operator": "eq", "value": true}'),
    ('convoy_leader', 'Convoy Leader', 'Lead 5 tourings as Lead or Marshal', 'social', 'silver',
     '{"type": "lead_count", "operator": "gte", "value": 5}'),
    ('social_butterfly', 'Social Butterfly', 'Join 20 different tourings', 'social', 'gold',
     '{"type": "unique_rooms", "operator": "gte", "value": 20}'),
    ('early_adopter', 'Early Adopter', 'Join Go Road in the first month', 'achievement', 'diamond',
     '{"type": "registration_date", "operator": "lte", "value": "2026-08-31"}'),
    ('perfect_attendance', 'Perfect Attendance', 'Complete a multi-day touring without missing a day', 'endurance', 'platinum',
     '{"type": "perfect_attendance", "operator": "eq", "value": true}'),
    ('globetrotter', 'Globetrotter', 'Tour in 5 different provinces/cities', 'distance', 'platinum',
     '{"type": "unique_locations", "operator": "gte", "value": 5}');

CREATE TABLE user_badges (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    badge_id        UUID NOT NULL REFERENCES badges(id) ON DELETE CASCADE,
    awarded_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    touring_id      UUID REFERENCES touring_rooms(id) ON DELETE SET NULL,
    UNIQUE(user_id, badge_id)
);

-- ============================================
-- SOCIAL
-- ============================================
CREATE TABLE user_follows (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    follower_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    following_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(follower_id, following_id)
);

CREATE TABLE user_blocks (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    blocker_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    blocked_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reason          TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(blocker_id, blocked_id)
);

CREATE TABLE reports (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    reporter_id     UUID NOT NULL REFERENCES users(id),
    reported_type   VARCHAR(20) NOT NULL CHECK (reported_type IN ('user', 'room', 'message', 'post', 'comment')),
    reported_id     UUID NOT NULL,
    reason          VARCHAR(30) NOT NULL CHECK (reason IN (
                        'spam', 'harassment', 'inappropriate', 'fake', 'violence',
                        'hate_speech', 'nudity', 'copyright', 'other'
                    )),
    description     TEXT,
    status          VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'reviewed', 'resolved', 'dismissed')),
    reviewed_by     UUID REFERENCES users(id),
    reviewed_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================
-- NOTIFICATIONS
-- ============================================
CREATE TABLE notifications (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type            VARCHAR(30) NOT NULL CHECK (type IN (
                        'system', 'room_invite', 'room_update', 'role_change',
                        'new_follower', 'badge_awarded', 'mention', 'comment',
                        'like', 'emergency', 'sos', 'weather_alert',
                        'service_reminder', 'touring_reminder', 'friend_request'
                    )),
    title           VARCHAR(200) NOT NULL,
    body            TEXT,
    data            JSONB DEFAULT '{}',
    is_read         BOOLEAN DEFAULT false,
    read_at         TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================
-- FEED / POSTS
-- ============================================
CREATE TABLE touring_posts (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    room_id         UUID REFERENCES touring_rooms(id) ON DELETE SET NULL,
    author_id       UUID NOT NULL REFERENCES users(id),
    caption         TEXT,
    photos          TEXT[],
    stats_snapshot  JSONB DEFAULT '{}',
    route_snapshot  JSONB DEFAULT '{}',
    is_public       BOOLEAN DEFAULT true,
    likes_count     INTEGER DEFAULT 0,
    comments_count  INTEGER DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE TABLE post_likes (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    post_id         UUID NOT NULL REFERENCES touring_posts(id) ON DELETE CASCADE,
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(post_id, user_id)
);

CREATE TABLE post_comments (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    post_id         UUID NOT NULL REFERENCES touring_posts(id) ON DELETE CASCADE,
    author_id       UUID NOT NULL REFERENCES users(id),
    content         TEXT NOT NULL,
    reply_to_id     UUID REFERENCES post_comments(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

-- Indexes
CREATE INDEX idx_user_badges_user ON user_badges(user_id);
CREATE INDEX idx_user_badges_badge ON user_badges(badge_id);
CREATE INDEX idx_user_follows_follower ON user_follows(follower_id);
CREATE INDEX idx_user_follows_following ON user_follows(following_id);
CREATE INDEX idx_user_blocks_blocker ON user_blocks(blocker_id);
CREATE INDEX idx_user_blocks_blocked ON user_blocks(blocked_id);
CREATE INDEX idx_reports_status ON reports(status);
CREATE INDEX idx_reports_reporter ON reports(reporter_id);
CREATE INDEX idx_notifications_user ON notifications(user_id, created_at DESC);
CREATE INDEX idx_notifications_unread ON notifications(user_id, is_read) WHERE is_read = false;
CREATE INDEX idx_touring_posts_author ON touring_posts(author_id);
CREATE INDEX idx_touring_posts_room ON touring_posts(room_id);
CREATE INDEX idx_touring_posts_created ON touring_posts(created_at DESC);
CREATE INDEX idx_touring_posts_public ON touring_posts(is_public) WHERE is_public = true;
CREATE INDEX idx_post_likes_post ON post_likes(post_id);
CREATE INDEX idx_post_comments_post ON post_comments(post_id, created_at);
