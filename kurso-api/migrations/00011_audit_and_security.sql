-- +goose Up

-- Immutable log of admin actions (Stage 2.1).
CREATE TABLE admin_audit_log (
    id          bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    admin_id    uuid REFERENCES admins (id) ON DELETE SET NULL,
    action      text NOT NULL,
    entity_type text,
    entity_id   text,
    changes     jsonb,
    ip          inet,
    user_agent  text,
    created_at  timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX admin_audit_log_admin_idx ON admin_audit_log (admin_id, created_at DESC);
CREATE INDEX admin_audit_log_entity_idx ON admin_audit_log (entity_type, entity_id);

-- Security-relevant events across all actor types (logins, 2FA, lockouts, ...).
CREATE TABLE security_events (
    id                bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    event_type        text NOT NULL,
    severity          text NOT NULL DEFAULT 'info' CHECK (severity IN ('info', 'warning', 'critical')),
    user_id           uuid REFERENCES users (id) ON DELETE SET NULL,
    admin_id          uuid REFERENCES admins (id) ON DELETE SET NULL,
    exchanger_user_id uuid REFERENCES exchanger_users (id) ON DELETE SET NULL,
    ip                inet,
    user_agent        text,
    metadata          jsonb,
    created_at        timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX security_events_type_created_idx ON security_events (event_type, created_at DESC);
CREATE INDEX security_events_severity_idx ON security_events (severity, created_at DESC);

-- +goose Down
DROP TABLE IF EXISTS security_events;
DROP TABLE IF EXISTS admin_audit_log;
