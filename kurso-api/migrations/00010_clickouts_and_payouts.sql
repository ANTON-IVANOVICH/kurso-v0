-- +goose Up

-- Outbound clicks to exchangers — the core monetisation event (Stage 1.9).
CREATE TABLE clickouts (
    id           bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    exchanger_id uuid NOT NULL REFERENCES exchangers (id) ON DELETE CASCADE,
    direction_id uuid REFERENCES directions (id) ON DELETE SET NULL,
    user_id      uuid REFERENCES users (id) ON DELETE SET NULL,
    ref_code     text,
    ip           inet,
    user_agent   text,
    referer      text,
    created_at   timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX clickouts_exchanger_created_idx ON clickouts (exchanger_id, created_at DESC);
CREATE INDEX clickouts_direction_created_idx ON clickouts (direction_id, created_at DESC);

-- Partner payouts owed to / reconciled with exchangers per period.
-- (Stage 9 extends the payout model to blogger referrals.)
CREATE TABLE partner_payouts (
    id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    exchanger_id uuid NOT NULL REFERENCES exchangers (id) ON DELETE CASCADE,
    period_start date NOT NULL,
    period_end   date NOT NULL,
    clicks_count integer NOT NULL DEFAULT 0,
    amount       numeric(18, 2) NOT NULL DEFAULT 0,
    currency     text NOT NULL DEFAULT 'USD',
    status       text NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'paid', 'cancelled')),
    note         text,
    created_at   timestamptz NOT NULL DEFAULT now(),
    updated_at   timestamptz NOT NULL DEFAULT now(),
    paid_at      timestamptz,
    CONSTRAINT partner_payouts_period_valid CHECK (period_end >= period_start)
);
CREATE INDEX partner_payouts_exchanger_idx ON partner_payouts (exchanger_id, period_start DESC);
CREATE TRIGGER partner_payouts_set_updated_at
    BEFORE UPDATE ON partner_payouts FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS partner_payouts;
DROP TABLE IF EXISTS clickouts;
