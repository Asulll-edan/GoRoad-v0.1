CREATE TABLE fuel_logs (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id         UUID NOT NULL REFERENCES users(id),
    motor_id        UUID REFERENCES motors(id),
    room_id         UUID REFERENCES touring_rooms(id) ON DELETE SET NULL,
    
    -- Fuel details
    fuel_type       VARCHAR(20) NOT NULL CHECK (fuel_type IN ('pertalite', 'pertamax', 'pertamax_turbo', 'solar', 'dexlite', 'pertamina_dex', 'other')),
    amount_liters   DECIMAL(8,2) NOT NULL,
    price_per_liter DECIMAL(10,2) NOT NULL,
    total_cost      DECIMAL(12,2) NOT NULL,
    station_name    VARCHAR(200),
    location        GEOGRAPHY(POINT, 4326),
    odometer_km     DECIMAL(10,1),
    is_full_tank    BOOLEAN DEFAULT false,
    
    logged_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE expenses (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id         UUID NOT NULL REFERENCES users(id),
    room_id         UUID REFERENCES touring_rooms(id) ON DELETE SET NULL,
    
    -- Expense details
    category        VARCHAR(30) NOT NULL CHECK (category IN (
                        'fuel', 'food', 'accommodation', 'toll', 'parking',
                        'ferry', 'entrance', 'equipment', 'medical', 'other'
                    )),
    amount          DECIMAL(12,2) NOT NULL,
    description     TEXT,
    location        GEOGRAPHY(POINT, 4326),
    receipt_url     VARCHAR(500),
    is_split_bill   BOOLEAN DEFAULT false,
    split_with      UUID[],
    
    logged_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE service_reminders (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id         UUID NOT NULL REFERENCES users(id),
    motor_id        UUID NOT NULL REFERENCES motors(id),
    service_type    VARCHAR(30) NOT NULL CHECK (service_type IN (
                        'oil_change', 'tire', 'brake', 'chain', 'spark_plug',
                        'coolant', 'battery', 'general', 'custom'
                    )),
    title           VARCHAR(200) NOT NULL,
    description     TEXT,
    due_date        DATE NOT NULL,
    due_odometer    DECIMAL(10,1),
    completed_at    TIMESTAMPTZ,
    is_recurring    BOOLEAN DEFAULT false,
    recurring_interval_days INTEGER,
    recurring_interval_km DECIMAL(10,1),
    notified_h7     BOOLEAN DEFAULT false,
    notified_h1     BOOLEAN DEFAULT false,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE checklist_templates (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_by      UUID NOT NULL REFERENCES users(id),
    name            VARCHAR(200) NOT NULL,
    description     TEXT,
    is_public       BOOLEAN DEFAULT true,
    category        VARCHAR(30) DEFAULT 'general' CHECK (category IN (
                        'general', 'safety', 'document', 'equipment', 'camping', 'long_distance', 'custom'
                    )),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE checklist_items (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    template_id     UUID NOT NULL REFERENCES checklist_templates(id) ON DELETE CASCADE,
    label           VARCHAR(200) NOT NULL,
    order_index     INTEGER NOT NULL,
    is_required     BOOLEAN DEFAULT false,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE touring_checklists (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    room_id         UUID NOT NULL REFERENCES touring_rooms(id) ON DELETE CASCADE,
    template_id     UUID REFERENCES checklist_templates(id),
    user_id         UUID NOT NULL REFERENCES users(id),
    item_id         UUID NOT NULL REFERENCES checklist_items(id),
    is_checked      BOOLEAN DEFAULT false,
    checked_at      TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(room_id, user_id, item_id)
);

CREATE INDEX idx_fuel_logs_user ON fuel_logs(user_id);
CREATE INDEX idx_fuel_logs_room ON fuel_logs(room_id);
CREATE INDEX idx_fuel_logs_logged ON fuel_logs(logged_at DESC);
CREATE INDEX idx_expenses_user ON expenses(user_id);
CREATE INDEX idx_expenses_room ON expenses(room_id);
CREATE INDEX idx_expenses_category ON expenses(category);
CREATE INDEX idx_expenses_logged ON expenses(logged_at DESC);
CREATE INDEX idx_service_reminders_user ON service_reminders(user_id);
CREATE INDEX idx_service_reminders_due ON service_reminders(due_date) WHERE completed_at IS NULL;
CREATE INDEX idx_checklist_templates_created_by ON checklist_templates(created_by);
CREATE INDEX idx_checklist_items_template ON checklist_items(template_id);
CREATE INDEX idx_touring_checklists_room ON touring_checklists(room_id);
CREATE INDEX idx_touring_checklists_user ON touring_checklists(user_id);
