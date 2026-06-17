ALTER TABLE expenses ADD COLUMN place TEXT NOT NULL DEFAULT '';

CREATE INDEX idx_expenses_household_place ON expenses(household_id, place) WHERE place <> '';
