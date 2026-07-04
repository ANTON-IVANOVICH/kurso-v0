-- +goose Up
CREATE TABLE exchangers (
    id                    uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    slug                  text NOT NULL UNIQUE,
    name                  text NOT NULL,
    status                text NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'paused', 'banned')),
    website_url           text,
    referral_url_template text,
    logo_url              text,
    description           text,
    rating_avg            numeric(3, 2),
    reviews_count         integer NOT NULL DEFAULT 0,
    is_verified           boolean NOT NULL DEFAULT false,
    created_at            timestamptz NOT NULL DEFAULT now(),
    updated_at            timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX exchangers_status_idx ON exchangers (status);

CREATE TRIGGER exchangers_set_updated_at
    BEFORE UPDATE ON exchangers
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS exchangers;
