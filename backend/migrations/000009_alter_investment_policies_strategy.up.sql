ALTER TABLE investment_policies 
DROP COLUMN IF EXISTS buy_rules,
DROP COLUMN IF EXISTS sell_rules,
DROP COLUMN IF EXISTS ai_rationale,
ADD COLUMN detailed_strategy TEXT;
