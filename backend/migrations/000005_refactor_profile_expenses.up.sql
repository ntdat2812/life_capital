ALTER TABLE investor_profiles DROP COLUMN IF EXISTS total_monthly_income;
ALTER TABLE investor_profiles RENAME COLUMN total_monthly_expense TO essential_monthly_expense;
ALTER TABLE investor_profiles ADD COLUMN discretionary_monthly_expense DECIMAL(15, 2) DEFAULT 0;
