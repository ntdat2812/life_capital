CREATE TABLE investor_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    version INT NOT NULL DEFAULT 1,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    date_of_birth DATE,
    marital_status VARCHAR(50),
    risk_tolerance VARCHAR(50),
    risk_score INT,
    total_monthly_income DECIMAL(15, 2),
    total_monthly_expense DECIMAL(15, 2),
    fi_target_amount DECIMAL(15, 2),
    life_constraints JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_investor_profiles_user_id ON investor_profiles(user_id);
