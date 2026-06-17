CREATE TABLE households (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    join_code TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_households_join_code ON households(join_code);

ALTER TABLE users ADD COLUMN household_id UUID REFERENCES households(id) ON DELETE CASCADE;

ALTER TABLE incomes ADD COLUMN household_id UUID REFERENCES households(id) ON DELETE CASCADE;
ALTER TABLE envelopes ADD COLUMN household_id UUID REFERENCES households(id) ON DELETE CASCADE;
ALTER TABLE expenses ADD COLUMN household_id UUID REFERENCES households(id) ON DELETE CASCADE;

-- Create a household for each existing user and link their budget data.
DO $$
DECLARE
    u RECORD;
    h_id UUID;
    code TEXT;
BEGIN
    FOR u IN SELECT id FROM users WHERE household_id IS NULL LOOP
        LOOP
            code := upper(substr(replace(gen_random_uuid()::text, '-', ''), 1, 8));
            BEGIN
                INSERT INTO households (join_code) VALUES (code) RETURNING id INTO h_id;
                EXIT;
            EXCEPTION WHEN unique_violation THEN
                CONTINUE;
            END;
        END LOOP;

        UPDATE users SET household_id = h_id WHERE id = u.id;
        UPDATE incomes SET household_id = h_id WHERE user_id = u.id AND household_id IS NULL;
        UPDATE envelopes SET household_id = h_id WHERE user_id = u.id AND household_id IS NULL;
        UPDATE expenses SET household_id = h_id WHERE user_id = u.id AND household_id IS NULL;
    END LOOP;
END $$;

ALTER TABLE users ALTER COLUMN household_id SET NOT NULL;

ALTER TABLE incomes ALTER COLUMN household_id SET NOT NULL;
ALTER TABLE envelopes ALTER COLUMN household_id SET NOT NULL;
ALTER TABLE expenses ALTER COLUMN household_id SET NOT NULL;

DROP INDEX IF EXISTS idx_incomes_user_id;
DROP INDEX IF EXISTS idx_envelopes_user_id;
DROP INDEX IF EXISTS idx_expenses_user_id;

ALTER TABLE incomes DROP COLUMN user_id;
ALTER TABLE envelopes DROP COLUMN user_id;
ALTER TABLE expenses DROP COLUMN user_id;

CREATE INDEX idx_incomes_household_id ON incomes(household_id);
CREATE INDEX idx_envelopes_household_id ON envelopes(household_id);
CREATE INDEX idx_expenses_household_id ON expenses(household_id);
