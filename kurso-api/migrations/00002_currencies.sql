-- +goose Up
CREATE TABLE currencies (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    code       text NOT NULL,
    name       text NOT NULL,
    kind       text NOT NULL CHECK (kind IN ('crypto', 'fiat', 'cash')),
    network    text,
    icon_url   text,
    is_active  boolean NOT NULL DEFAULT true,
    sort_order integer NOT NULL DEFAULT 0,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);
-- A currency is identified by its code + network (network is empty for fiat/cash).
CREATE UNIQUE INDEX currencies_code_network_key
    ON currencies (lower(code), coalesce(lower(network), ''));

CREATE TRIGGER currencies_set_updated_at
    BEFORE UPDATE ON currencies
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Maps parser/exchanger strings onto the canonical currency (Stage 1 normaliser).
CREATE TABLE currency_aliases (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    currency_id uuid NOT NULL REFERENCES currencies (id) ON DELETE CASCADE,
    alias       citext NOT NULL,
    source      text,
    created_at  timestamptz NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX currency_aliases_source_alias_key
    ON currency_aliases (coalesce(source, ''), lower(alias));
CREATE INDEX currency_aliases_alias_trgm
    ON currency_aliases USING gin (alias gin_trgm_ops);

-- +goose Down
DROP TABLE IF EXISTS currency_aliases;
DROP TABLE IF EXISTS currencies;
