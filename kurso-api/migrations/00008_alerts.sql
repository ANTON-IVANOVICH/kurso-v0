-- +goose Up

-- User-configured rate alerts.
CREATE TABLE alerts (
    id                uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id           uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    direction_id      uuid NOT NULL REFERENCES directions (id) ON DELETE CASCADE,
    condition         text NOT NULL CHECK (condition IN ('above', 'below')),
    threshold         numeric(38, 18) NOT NULL,
    channels          text[] NOT NULL DEFAULT ARRAY['telegram']::text[],
    cooldown_seconds  integer NOT NULL DEFAULT 3600 CHECK (cooldown_seconds >= 0),
    is_active         boolean NOT NULL DEFAULT true,
    last_triggered_at timestamptz,
    created_at        timestamptz NOT NULL DEFAULT now(),
    updated_at        timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX alerts_user_idx ON alerts (user_id);
CREATE INDEX alerts_active_direction_idx ON alerts (direction_id) WHERE is_active;
CREATE TRIGGER alerts_set_updated_at
    BEFORE UPDATE ON alerts FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Log of alert firings and their delivery outcome.
CREATE TABLE triggered_alerts (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    alert_id        uuid NOT NULL REFERENCES alerts (id) ON DELETE CASCADE,
    exchanger_id    uuid REFERENCES exchangers (id) ON DELETE SET NULL,
    rate            numeric(38, 18) NOT NULL,
    delivery_status text NOT NULL DEFAULT 'pending' CHECK (delivery_status IN ('pending', 'sent', 'failed')),
    triggered_at    timestamptz NOT NULL DEFAULT now(),
    delivered_at    timestamptz
);
CREATE INDEX triggered_alerts_alert_idx ON triggered_alerts (alert_id, triggered_at DESC);

-- Web Push subscriptions (Stage 6).
CREATE TABLE push_subscriptions (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    endpoint   text NOT NULL,
    p256dh     text NOT NULL,
    auth       text NOT NULL,
    user_agent text,
    created_at timestamptz NOT NULL DEFAULT now()
);
-- endpoint URLs can exceed btree index limits, so key on its hash.
CREATE UNIQUE INDEX push_subscriptions_user_endpoint_key
    ON push_subscriptions (user_id, md5(endpoint));

-- +goose Down
DROP TABLE IF EXISTS push_subscriptions;
DROP TABLE IF EXISTS triggered_alerts;
DROP TABLE IF EXISTS alerts;
