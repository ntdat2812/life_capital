CREATE TABLE IF NOT EXISTS life_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    event_date DATE NOT NULL DEFAULT CURRENT_DATE,
    category life_event_category NOT NULL,
    title VARCHAR(255) NOT NULL,
    income_impact DECIMAL(15,2) DEFAULT 0.00,
    expense_impact DECIMAL(15,2) DEFAULT 0.00,
    ai_impact_analysis TEXT,
    triggered_profile_version INT,
    triggered_ips_version INT,
    requires_ips_update BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_life_events_user_id ON life_events(user_id);
CREATE INDEX idx_life_events_event_date ON life_events(event_date);
