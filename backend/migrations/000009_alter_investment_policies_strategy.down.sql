ALTER TABLE investment_policies 
DROP COLUMN IF EXISTS detailed_strategy,
ADD COLUMN buy_rules JSONB,
ADD COLUMN sell_rules JSONB,
ADD COLUMN ai_rationale TEXT;
