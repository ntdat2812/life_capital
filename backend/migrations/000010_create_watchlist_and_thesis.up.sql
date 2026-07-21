-- =============================================
-- INVESTMENT THESIS (Module 5)
-- =============================================

CREATE TABLE investment_theses (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID REFERENCES users(id) ON DELETE CASCADE,
    ticker          VARCHAR(20) NOT NULL,
    company_name    VARCHAR(255) NOT NULL,
    status          VARCHAR(20) DEFAULT 'active',
    why_i_own           TEXT NOT NULL,
    thesis_summary      TEXT,
    thesis_detail       TEXT,
    moat                JSONB,
    catalysts           JSONB,
    risks               JSONB,
    key_metrics         JSONB,
    sell_conditions     JSONB,
    conviction_score    INTEGER CHECK (conviction_score BETWEEN 1 AND 10),
    quality_score       INTEGER CHECK (quality_score BETWEEN 1 AND 10),
    valuation_score     INTEGER CHECK (valuation_score BETWEEN 1 AND 10),
    fair_value          DECIMAL(18,2),
    margin_of_safety    DECIMAL(5,4),
    initial_date        DATE,
    last_reviewed       DATE,
    version             INTEGER DEFAULT 1,
    notes               TEXT,
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW()
);

-- =============================================
-- WATCHLIST (Module 3 Extension)
-- =============================================

CREATE TABLE watchlist (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID REFERENCES users(id) ON DELETE CASCADE,
    ticker          VARCHAR(20) NOT NULL,
    company_name    VARCHAR(255) NOT NULL,
    thesis_id       UUID REFERENCES investment_theses(id) ON DELETE SET NULL,
    added_date      DATE DEFAULT CURRENT_DATE,
    target_price    DECIMAL(18,2),
    current_price   DECIMAL(18,2),
    fair_value      DECIMAL(18,2),
    quality_score   INTEGER CHECK (quality_score BETWEEN 1 AND 10),
    status          VARCHAR(20) DEFAULT 'watching',
    priority        INTEGER DEFAULT 5,
    notes           TEXT,
    ai_alert        TEXT,
    last_ai_check   TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_theses_user ON investment_theses(user_id);
CREATE INDEX idx_theses_ticker ON investment_theses(user_id, ticker);
CREATE INDEX idx_watchlist_user ON watchlist(user_id);
