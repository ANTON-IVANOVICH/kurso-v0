-- +goose Up
CREATE TABLE parser_configs (
    id               uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    exchanger_id     uuid NOT NULL REFERENCES exchangers (id) ON DELETE CASCADE,
    type             text NOT NULL CHECK (type IN ('bestchange_xml', 'rest_api', 'manual', 'push')),
    name             text,
    source_url       text,
    config           jsonb NOT NULL DEFAULT '{}'::jsonb,
    interval_seconds integer NOT NULL DEFAULT 60 CHECK (interval_seconds > 0),
    is_enabled       boolean NOT NULL DEFAULT true,
    last_run_at      timestamptz,
    last_success_at  timestamptz,
    last_error       text,
    success_rate     numeric(5, 2),
    created_at       timestamptz NOT NULL DEFAULT now(),
    updated_at       timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX parser_configs_exchanger_idx ON parser_configs (exchanger_id);
CREATE INDEX parser_configs_enabled_idx ON parser_configs (is_enabled) WHERE is_enabled;

CREATE TRIGGER parser_configs_set_updated_at
    BEFORE UPDATE ON parser_configs FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS parser_configs;
