CREATE TABLE income_receipts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    household_id UUID NOT NULL REFERENCES households(id) ON DELETE CASCADE,
    amount_cents BIGINT NOT NULL CHECK (amount_cents > 0),
    description TEXT NOT NULL DEFAULT '',
    income_date DATE NOT NULL DEFAULT CURRENT_DATE,
    voided BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_income_receipts_household_id ON income_receipts(household_id);
CREATE INDEX idx_income_receipts_household_date ON income_receipts(household_id, income_date);
