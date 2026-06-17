ALTER TABLE envelopes ADD COLUMN balance_cents BIGINT NOT NULL DEFAULT 0;

ALTER TABLE households ADD COLUMN last_budget_month DATE;

-- Seed balance from current monthly allocation minus all-time spending.
UPDATE envelopes e
SET balance_cents = e.allocated_cents - COALESCE((
    SELECT SUM(ex.amount_cents)
    FROM expenses ex
    WHERE ex.envelope_id = e.id AND NOT ex.voided
), 0);

-- Current month is already reflected in the seeded balance; skip re-funding it.
UPDATE households SET last_budget_month = date_trunc('month', CURRENT_DATE)::date;
