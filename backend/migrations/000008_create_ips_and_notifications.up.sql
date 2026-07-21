CREATE TABLE investment_policies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    version INT NOT NULL,
    profile_version INT NOT NULL,
    trigger_event_id UUID REFERENCES life_events(id) ON DELETE SET NULL,
    target_allocation JSONB NOT NULL,
    buy_rules JSONB,
    sell_rules JSONB,
    ai_rationale TEXT,
    is_ai_recommended BOOLEAN DEFAULT FALSE,
    status VARCHAR(50) DEFAULT 'active', -- 'draft', 'active', 'superseded'
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_investment_policies_user_id ON investment_policies(user_id);
CREATE INDEX idx_investment_policies_status ON investment_policies(status);

CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    action_link VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_is_read ON notifications(is_read);
