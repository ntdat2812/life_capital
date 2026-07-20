-- =============================================
-- BASE CURRENCY CONFIG
-- =============================================

ALTER TABLE users ADD COLUMN base_currency VARCHAR(3) DEFAULT 'VND';

-- =============================================
-- ASSETS & LIABILITIES (Module 2)
-- =============================================

CREATE TABLE assets (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID REFERENCES users(id) ON DELETE CASCADE,
    category        asset_category NOT NULL,
    name            VARCHAR(255) NOT NULL,
    ticker          VARCHAR(20),
    quantity        DECIMAL(18,6),
    avg_price       DECIMAL(18,2),
    current_price   DECIMAL(18,2),
    current_value   DECIMAL(18,2),
    cost_basis      DECIMAL(18,2),
    notes           TEXT,
    is_active       BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE liabilities (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID REFERENCES users(id) ON DELETE CASCADE,
    category        liability_category NOT NULL,
    name            VARCHAR(255) NOT NULL,
    remaining_balance DECIMAL(18,2) NOT NULL,
    interest_rate   DECIMAL(5,4),
    monthly_payment DECIMAL(18,2) DEFAULT 0,
    lender          VARCHAR(255),
    notes           TEXT,
    is_active       BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_assets_user ON assets(user_id);
CREATE INDEX idx_liabilities_user ON liabilities(user_id);
