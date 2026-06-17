DROP INDEX IF EXISTS idx_income_receipts_recurring_unique;
ALTER TABLE income_receipts DROP COLUMN IF EXISTS auto_generated;
ALTER TABLE income_receipts DROP COLUMN IF EXISTS income_id;
