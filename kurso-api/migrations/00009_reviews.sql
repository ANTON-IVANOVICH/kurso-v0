-- +goose Up

-- User reviews of exchangers with a moderation state machine.
CREATE TABLE reviews (
    id                uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    exchanger_id      uuid NOT NULL REFERENCES exchangers (id) ON DELETE CASCADE,
    user_id           uuid REFERENCES users (id) ON DELETE SET NULL,
    author_name       text,
    author_email      citext,
    rating            smallint NOT NULL CHECK (rating BETWEEN 1 AND 5),
    title             text,
    body              text NOT NULL,
    status            text NOT NULL DEFAULT 'pending'
                          CHECK (status IN ('pending', 'published', 'rejected', 'needs_info')),
    moderation_reason text,
    ip                inet,
    created_at        timestamptz NOT NULL DEFAULT now(),
    updated_at        timestamptz NOT NULL DEFAULT now(),
    published_at      timestamptz
);
CREATE INDEX reviews_exchanger_status_idx ON reviews (exchanger_id, status);
CREATE INDEX reviews_status_idx ON reviews (status) WHERE status = 'pending';
CREATE TRIGGER reviews_set_updated_at
    BEFORE UPDATE ON reviews FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- One reply per review from an exchanger representative (Stage 5.5).
CREATE TABLE review_replies (
    id                uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    review_id         uuid NOT NULL UNIQUE REFERENCES reviews (id) ON DELETE CASCADE,
    exchanger_user_id uuid REFERENCES exchanger_users (id) ON DELETE SET NULL,
    body              text NOT NULL,
    status            text NOT NULL DEFAULT 'pending'
                          CHECK (status IN ('pending', 'published', 'rejected')),
    created_at        timestamptz NOT NULL DEFAULT now(),
    updated_at        timestamptz NOT NULL DEFAULT now()
);
CREATE TRIGGER review_replies_set_updated_at
    BEFORE UPDATE ON review_replies FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Abuse reports against reviews (Stage 4.4).
CREATE TABLE review_reports (
    id               uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    review_id        uuid NOT NULL REFERENCES reviews (id) ON DELETE CASCADE,
    reporter_user_id uuid REFERENCES users (id) ON DELETE SET NULL,
    reason           text NOT NULL,
    details          text,
    status           text NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'reviewed', 'dismissed')),
    created_at       timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX review_reports_review_idx ON review_reports (review_id);
CREATE INDEX review_reports_status_idx ON review_reports (status) WHERE status = 'open';

-- +goose Down
DROP TABLE IF EXISTS review_reports;
DROP TABLE IF EXISTS review_replies;
DROP TABLE IF EXISTS reviews;
