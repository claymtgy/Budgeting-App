ALTER TABLE income_receipts ADD COLUMN income_id UUID REFERENCES incomes(id) ON DELETE SET NULL;
ALTER TABLE income_receipts ADD COLUMN auto_generated BOOLEAN NOT NULL DEFAULT FALSE;

CREATE UNIQUE INDEX idx_income_receipts_recurring_unique
  ON income_receipts (household_id, income_id, (date_trunc('month', income_date::timestamp)::date))
  WHERE income_id IS NOT NULL AND auto_generated AND NOT voided;

-- Seed current month for existing recurring sources.
INSERT INTO income_receipts (household_id, income_id, amount_cents, description, income_date, auto_generated)
SELECT
  i.household_id,
  i.id,
  i.amount_cents,
  i.name,
  date_trunc('month', CURRENT_DATE)::date,
  TRUE
FROM incomes i
WHERE NOT EXISTS (
  SELECT 1 FROM income_receipts r
  WHERE r.household_id = i.household_id
    AND r.income_id = i.id
    AND r.auto_generated
    AND NOT r.voided
    AND date_trunc('month', r.income_date) = date_trunc('month', CURRENT_DATE)
);
