CREATE TABLE income_streams (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID REFERENCES users(id) ON DELETE CASCADE,
    name            VARCHAR(255) NOT NULL,
    type            VARCHAR(50) NOT NULL,
    is_passive      BOOLEAN DEFAULT FALSE,
    amount          DECIMAL(18,2) NOT NULL,
    frequency       VARCHAR(20) DEFAULT 'monthly',
    is_active       BOOLEAN DEFAULT TRUE,
    start_date      DATE,
    end_date        DATE,
    notes           TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE dependents (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID REFERENCES users(id) ON DELETE CASCADE,
    name            VARCHAR(255) NOT NULL,
    relationship    VARCHAR(50) NOT NULL,
    date_of_birth   DATE,
    is_active       BOOLEAN DEFAULT TRUE,
    monthly_cost    DECIMAL(18,2),
    notes           TEXT,
    added_date      DATE DEFAULT CURRENT_DATE,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_income_streams_user ON income_streams(user_id);
CREATE INDEX idx_dependents_user ON dependents(user_id);
