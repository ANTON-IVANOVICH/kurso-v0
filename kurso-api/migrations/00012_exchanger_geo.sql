-- +goose Up

-- Physical location for exchangers with a cash desk, so the map can plot them.
-- All nullable: online-only exchangers simply have no coordinates and never
-- appear on the map.
ALTER TABLE exchangers
    ADD COLUMN latitude  numeric(9, 6),
    ADD COLUMN longitude numeric(9, 6),
    ADD COLUMN address   text,
    ADD COLUMN city      text,
    ADD COLUMN hours     text;

-- Only located exchangers are queried for the map.
CREATE INDEX exchangers_geo_idx ON exchangers (latitude, longitude)
    WHERE latitude IS NOT NULL AND longitude IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS exchangers_geo_idx;
ALTER TABLE exchangers
    DROP COLUMN latitude,
    DROP COLUMN longitude,
    DROP COLUMN address,
    DROP COLUMN city,
    DROP COLUMN hours;
