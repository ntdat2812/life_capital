-- =============================================
-- USERS & ENUMS
-- =============================================

CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email           VARCHAR(255) UNIQUE NOT NULL,
    name            VARCHAR(255) NOT NULL,
    password_hash   VARCHAR(255) NULL,
    auth_provider   VARCHAR(50) DEFAULT 'local',
    google_id       VARCHAR(255) UNIQUE,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TYPE life_event_category AS ENUM (
    'income_change',
    'family_change',
    'housing',
    'dependent_change',
    'inheritance',
    'health',
    'education',
    'career',
    'windfall',
    'major_expense',
    'other'
);

CREATE TYPE asset_category AS ENUM (
    'cash',
    'deposit',
    'gold',
    'stock',
    'fund',
    'crypto',
    'real_estate'
);

CREATE TYPE liability_category AS ENUM (
    'mortgage',
    'auto_loan',
    'student_loan',
    'credit_card',
    'personal_loan',
    'other'
);

CREATE TYPE decision_type AS ENUM ('buy', 'sell', 'hold', 'add', 'reduce', 'skip');
