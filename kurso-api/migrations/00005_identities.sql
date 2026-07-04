-- +goose Up

-- End users of kurso-web (alerts, reviews, partner program).
CREATE TABLE users (
    id                uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email             citext UNIQUE,
    email_verified_at timestamptz,
    password_hash     text,
    telegram_id       bigint UNIQUE,
    google_id         text UNIQUE,
    display_name      text,
    status            text NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'blocked', 'deleted')),
    created_at        timestamptz NOT NULL DEFAULT now(),
    updated_at        timestamptz NOT NULL DEFAULT now(),
    -- A user must be reachable through at least one identity.
    CONSTRAINT users_has_identity CHECK (email IS NOT NULL OR telegram_id IS NOT NULL OR google_id IS NOT NULL)
);
CREATE TRIGGER users_set_updated_at
    BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Platform administrators (admin.kurso.io).
CREATE TABLE admins (
    id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email         citext NOT NULL UNIQUE,
    password_hash text NOT NULL,
    role          text NOT NULL DEFAULT 'moderator' CHECK (role IN ('superadmin', 'moderator')),
    totp_secret   text,
    totp_enabled  boolean NOT NULL DEFAULT false,
    last_login_at timestamptz,
    status        text NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'disabled')),
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now()
);
CREATE TRIGGER admins_set_updated_at
    BEFORE UPDATE ON admins FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Representatives of an exchanger (partner.kurso.io).
CREATE TABLE exchanger_users (
    id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    exchanger_id  uuid NOT NULL REFERENCES exchangers (id) ON DELETE CASCADE,
    email         citext NOT NULL,
    password_hash text,
    role          text NOT NULL DEFAULT 'viewer' CHECK (role IN ('owner', 'manager', 'viewer')),
    totp_secret   text,
    totp_enabled  boolean NOT NULL DEFAULT false,
    status        text NOT NULL DEFAULT 'invited' CHECK (status IN ('invited', 'active', 'disabled')),
    invited_at    timestamptz NOT NULL DEFAULT now(),
    accepted_at   timestamptz,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT exchanger_users_email_key UNIQUE (exchanger_id, email)
);
CREATE INDEX exchanger_users_exchanger_idx ON exchanger_users (exchanger_id);
CREATE TRIGGER exchanger_users_set_updated_at
    BEFORE UPDATE ON exchanger_users FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS exchanger_users;
DROP TABLE IF EXISTS admins;
DROP TABLE IF EXISTS users;
