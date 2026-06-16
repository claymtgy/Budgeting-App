DROP INDEX IF EXISTS idx_expenses_household_id;
DROP INDEX IF EXISTS idx_envelopes_household_id;
DROP INDEX IF EXISTS idx_incomes_household_id;

ALTER TABLE expenses ADD COLUMN user_id UUID;
ALTER TABLE envelopes ADD COLUMN user_id UUID;
ALTER TABLE incomes ADD COLUMN user_id UUID;

UPDATE expenses e SET user_id = (
    SELECT u.id FROM users u WHERE u.household_id = e.household_id LIMIT 1
);
UPDATE envelopes e SET user_id = (
    SELECT u.id FROM users u WHERE u.household_id = e.household_id LIMIT 1
);
UPDATE incomes i SET user_id = (
    SELECT u.id FROM users u WHERE u.household_id = i.household_id LIMIT 1
);

ALTER TABLE expenses DROP COLUMN household_id;
ALTER TABLE envelopes DROP COLUMN household_id;
ALTER TABLE incomes DROP COLUMN household_id;

ALTER TABLE users DROP COLUMN household_id;

DROP TABLE households;

CREATE INDEX idx_incomes_user_id ON incomes(user_id);
CREATE INDEX idx_envelopes_user_id ON envelopes(user_id);
CREATE INDEX idx_expenses_user_id ON expenses(user_id);
