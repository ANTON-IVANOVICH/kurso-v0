-- +goose Up

-- Current rate per (exchanger, direction) — hot table the public API reads.
CREATE TABLE rates (
    id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    exchanger_id uuid NOT NULL REFERENCES exchangers (id) ON DELETE CASCADE,
    direction_id uuid NOT NULL REFERENCES directions (id) ON DELETE CASCADE,
    rate         numeric(38, 18) NOT NULL,
    reserve      numeric(38, 8),
    min_amount   numeric(38, 8),
    max_amount   numeric(38, 8),
    is_active    boolean NOT NULL DEFAULT true,
    fetched_at   timestamptz NOT NULL DEFAULT now(),
    created_at   timestamptz NOT NULL DEFAULT now(),
    updated_at   timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT rates_exchanger_direction_key UNIQUE (exchanger_id, direction_id)
);
CREATE INDEX rates_direction_idx ON rates (direction_id) WHERE is_active;
CREATE TRIGGER rates_set_updated_at
    BEFORE UPDATE ON rates FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Append-only history of every observed rate, partitioned by month on recorded_at.
-- No FKs: this is a high-volume log; the partition key must be part of the PK.
CREATE TABLE rates_history (
    id           uuid NOT NULL DEFAULT gen_random_uuid(),
    exchanger_id uuid NOT NULL,
    direction_id uuid NOT NULL,
    rate         numeric(38, 18) NOT NULL,
    reserve      numeric(38, 8),
    recorded_at  timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (id, recorded_at)
) PARTITION BY RANGE (recorded_at);

CREATE INDEX rates_history_direction_idx ON rates_history (direction_id, recorded_at DESC);
CREATE INDEX rates_history_exchanger_idx ON rates_history (exchanger_id, recorded_at DESC);

-- DEFAULT partition guarantees inserts always land somewhere. A partition
-- manager (pg_partman / cron job) creates monthly partitions in later stages.
CREATE TABLE rates_history_default PARTITION OF rates_history DEFAULT;

-- +goose Down
DROP TABLE IF EXISTS rates_history;
DROP TABLE IF EXISTS rates;
