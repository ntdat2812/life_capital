ALTER TABLE investor_profiles ADD COLUMN total_monthly_income DECIMAL(15, 2);
ALTER TABLE investor_profiles DROP COLUMN IF EXISTS discretionary_monthly_expense;
ALTER TABLE investor_profiles RENAME COLUMN essential_monthly_expense TO total_monthly_expense;
