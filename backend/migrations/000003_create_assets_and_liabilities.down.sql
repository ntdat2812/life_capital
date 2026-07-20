DROP INDEX IF EXISTS idx_liabilities_user;
DROP INDEX IF EXISTS idx_assets_user;

DROP TABLE IF EXISTS liabilities;
DROP TABLE IF EXISTS assets;

ALTER TABLE users DROP COLUMN IF EXISTS base_currency;
