-- +goose Up

-- Affiliate program: every user gets a short referral code and may be attributed
-- to the user who referred them (via the ?ref= link / kurso_ref cookie).
ALTER TABLE users
    ADD COLUMN referral_code text UNIQUE,
    ADD COLUMN referred_by   uuid REFERENCES users (id) ON DELETE SET NULL;

CREATE INDEX users_referred_by_idx ON users (referred_by) WHERE referred_by IS NOT NULL;

-- Backfill a stable code (first 8 hex of the id) for existing users.
UPDATE users SET referral_code = substr(replace(id::text, '-', ''), 1, 8) WHERE referral_code IS NULL;

-- +goose Down
DROP INDEX IF EXISTS users_referred_by_idx;
ALTER TABLE users
    DROP COLUMN referral_code,
    DROP COLUMN referred_by;
