-- +goose Up
CREATE TABLE directions (
    id                uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    from_currency_id  uuid NOT NULL REFERENCES currencies (id) ON DELETE RESTRICT,
    to_currency_id    uuid NOT NULL REFERENCES currencies (id) ON DELETE RESTRICT,
    slug              text NOT NULL UNIQUE,
    is_popular        boolean NOT NULL DEFAULT false,
    is_active         boolean NOT NULL DEFAULT true,
    sort_order        integer NOT NULL DEFAULT 0,
    created_at        timestamptz NOT NULL DEFAULT now(),
    updated_at        timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT directions_pair_key UNIQUE (from_currency_id, to_currency_id),
    CONSTRAINT directions_distinct_currencies CHECK (from_currency_id <> to_currency_id)
);
CREATE INDEX directions_from_idx ON directions (from_currency_id);
CREATE INDEX directions_to_idx ON directions (to_currency_id);
CREATE INDEX directions_popular_idx ON directions (is_popular) WHERE is_popular;

CREATE TRIGGER directions_set_updated_at
    BEFORE UPDATE ON directions
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS directions;
