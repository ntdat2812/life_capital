-- =============================================
-- MONTHLY REVIEW (Module 8)
-- =============================================

CREATE TABLE monthly_reviews (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID REFERENCES users(id) ON DELETE CASCADE,
    review_month    DATE NOT NULL, -- Stored as YYYY-MM-01
    status          VARCHAR(20) DEFAULT 'draft',
    new_investment_amount   DECIMAL(18,2) DEFAULT 0,
    
    -- Context snapshot at the time of review
    portfolio_snapshot      JSONB,
    net_worth_at_review     DECIMAL(18,2),
    
    -- AI Generated content
    ai_recommendations          JSONB,
    ai_overall_summary          TEXT,
    
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, review_month)
);

CREATE INDEX idx_monthly_reviews_user ON monthly_reviews(user_id);
