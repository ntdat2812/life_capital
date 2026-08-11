CREATE TABLE IF NOT EXISTS asset_alerts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    alert_type VARCHAR(50) NOT NULL,
    target_value DECIMAL(18,2) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    is_triggered BOOLEAN DEFAULT FALSE,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_asset_alerts_user ON asset_alerts(user_id);
CREATE INDEX idx_asset_alerts_asset ON asset_alerts(asset_id);
