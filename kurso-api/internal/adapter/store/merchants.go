package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/ANTON-IVANOVICH/kurso-v0/kurso-api/internal/domain"
	"github.com/jackc/pgx/v5"
)

// staleAfter marks a rate's feed "delayed" once its last fetch is older than
// this. The dev ticker refreshes every few seconds, so healthy feeds stay green.
const staleAfterSQL = "now() - interval '3 minutes'"

// ExchangerUserByEmail loads a merchant representative by email (any exchanger).
func (s *Store) ExchangerUserByEmail(ctx context.Context, email string) (domain.ExchangerUser, error) {
	var u domain.ExchangerUser
	err := s.db.QueryRow(ctx, `
		SELECT id::text, exchanger_id::text, email, COALESCE(password_hash, ''), role, status
		FROM exchanger_users
		WHERE email = $1
		ORDER BY created_at
		LIMIT 1`, email).
		Scan(&u.ID, &u.ExchangerID, &u.Email, &u.PasswordHash, &u.Role, &u.Status)
	if err != nil {
		return domain.ExchangerUser{}, fmt.Errorf("query exchanger user by email: %w", err)
	}
	return u, nil
}

// TouchExchangerUserLogin stamps the representative's accepted_at on login.
func (s *Store) TouchExchangerUserLogin(ctx context.Context, id string) error {
	_, err := s.db.Exec(ctx, `
		UPDATE exchanger_users
		SET accepted_at = COALESCE(accepted_at, now())
		WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("touch exchanger user login: %w", err)
	}
	return nil
}

// MerchantIdentity resolves the exchanger an authenticated representative owns.
func (s *Store) MerchantIdentity(ctx context.Context, userID string) (domain.MerchantIdentity, error) {
	var m domain.MerchantIdentity
	err := s.db.QueryRow(ctx, `
		SELECT eu.id::text, eu.email, eu.role, e.id::text, e.slug, e.name
		FROM exchanger_users eu
		JOIN exchangers e ON e.id = eu.exchanger_id
		WHERE eu.id = $1`, userID).
		Scan(&m.UserID, &m.Email, &m.Role, &m.ExchangerID, &m.ExchangerSlug, &m.ExchangerName)
	if err != nil {
		return domain.MerchantIdentity{}, fmt.Errorf("merchant identity: %w", err)
	}
	return m, nil
}

// EnsureExchangerUser seeds a representative for an exchanger (by slug) if the
// email is absent, so a fresh database has a working merchant login. No-op when
// the exchanger does not exist yet or the email already belongs to it.
func (s *Store) EnsureExchangerUser(ctx context.Context, exchangerSlug, email, passwordHash, role string) error {
	_, err := s.db.Exec(ctx, `
		INSERT INTO exchanger_users (exchanger_id, email, password_hash, role, status, accepted_at)
		SELECT e.id, $2, $3, $4, 'active', now()
		FROM exchangers e
		WHERE e.slug = $1
		ON CONFLICT (exchanger_id, email) DO NOTHING`, exchangerSlug, email, passwordHash, role)
	if err != nil {
		return fmt.Errorf("ensure exchanger user: %w", err)
	}
	return nil
}

// MerchantMetrics computes the dashboard headline numbers for one exchanger.
func (s *Store) MerchantMetrics(ctx context.Context, exchangerID string) (domain.MerchantMetrics, error) {
	var m domain.MerchantMetrics
	err := s.db.QueryRow(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE is_active AND fetched_at >= `+staleAfterSQL+`)                AS active_fresh,
			COUNT(*) FILTER (WHERE is_active)                                                     AS active_total,
			COUNT(*) FILTER (WHERE is_active AND fetched_at < `+staleAfterSQL+`)                  AS stale
		FROM rates WHERE exchanger_id = $1`, exchangerID).
		Scan(&m.RatesActive, &m.RatesTotal, &m.RatesStale)
	if err != nil {
		return domain.MerchantMetrics{}, fmt.Errorf("merchant rate metrics: %w", err)
	}

	err = s.db.QueryRow(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE created_at >= date_trunc('day', now())),
			COUNT(*) FILTER (WHERE created_at >= date_trunc('day', now()) - interval '1 day'
			                 AND created_at < date_trunc('day', now()))
		FROM clickouts WHERE exchanger_id = $1`, exchangerID).
		Scan(&m.ClicksToday, &m.ClicksYesterday)
	if err != nil {
		return domain.MerchantMetrics{}, fmt.Errorf("merchant click metrics: %w", err)
	}

	err = s.db.QueryRow(ctx, `
		SELECT
			COALESCE(round(avg(rating)::numeric, 2), 0)::float8,
			COUNT(*),
			COUNT(*) FILTER (WHERE NOT EXISTS (
				SELECT 1 FROM review_replies rr WHERE rr.review_id = r.id))
		FROM reviews r
		WHERE r.exchanger_id = $1 AND r.status = 'published'`, exchangerID).
		Scan(&m.RatingAvg, &m.ReviewsCount, &m.ReviewsUnanswered)
	if err != nil {
		return domain.MerchantMetrics{}, fmt.Errorf("merchant review metrics: %w", err)
	}
	return m, nil
}

// MerchantRates lists the exchanger's directions with the current rate + feed
// freshness, popular directions first.
func (s *Store) MerchantRates(ctx context.Context, exchangerID string) ([]domain.MerchantRate, error) {
	rows, err := s.db.Query(ctx, `
		SELECT d.id::text, d.slug, fc.code, tc.code,
		       r.rate::text, r.reserve::text, r.fetched_at, COALESCE(r.is_active, false)
		FROM rates r
		JOIN directions d ON d.id = r.direction_id
		JOIN currencies fc ON fc.id = d.from_currency_id
		JOIN currencies tc ON tc.id = d.to_currency_id
		WHERE r.exchanger_id = $1
		ORDER BY d.is_popular DESC, d.sort_order, d.slug`, exchangerID)
	if err != nil {
		return nil, fmt.Errorf("query merchant rates: %w", err)
	}
	defer rows.Close()

	var out []domain.MerchantRate
	for rows.Next() {
		var mr domain.MerchantRate
		if err := rows.Scan(&mr.DirectionID, &mr.DirectionSlug, &mr.FromCode, &mr.ToCode,
			&mr.Rate, &mr.Reserve, &mr.FetchedAt, &mr.IsActive); err != nil {
			return nil, fmt.Errorf("scan merchant rate: %w", err)
		}
		out = append(out, mr)
	}
	return out, rows.Err()
}

// DirectionByID looks up a single direction by id. ok is false when not found.
func (s *Store) DirectionByID(ctx context.Context, id string) (domain.Direction, bool, error) {
	d, err := scanDirection(s.db.QueryRow(ctx, directionSelect+` WHERE d.id = $1`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Direction{}, false, nil
	}
	if err != nil {
		return domain.Direction{}, false, fmt.Errorf("query direction by id: %w", err)
	}
	return d, true, nil
}

// TouchRate re-stamps a rate's fetched_at (a successful manual feed poll). ok is
// false when the exchanger has no rate for that direction.
func (s *Store) TouchRate(ctx context.Context, exchangerID, directionID string) (bool, error) {
	tag, err := s.db.Exec(ctx, `
		UPDATE rates SET fetched_at = now(), is_active = true
		WHERE exchanger_id = $1 AND direction_id = $2`, exchangerID, directionID)
	if err != nil {
		return false, fmt.Errorf("touch rate: %w", err)
	}
	return tag.RowsAffected() > 0, nil
}

// MerchantReviews lists an exchanger's published reviews with any reply.
func (s *Store) MerchantReviews(ctx context.Context, exchangerID string) ([]domain.MerchantReview, error) {
	rows, err := s.db.Query(ctx, `
		SELECT r.id::text, COALESCE(r.author_name, 'Аноним'), r.rating, r.title, r.body, r.created_at,
		       rr.body, rr.created_at
		FROM reviews r
		LEFT JOIN review_replies rr ON rr.review_id = r.id
		WHERE r.exchanger_id = $1 AND r.status = 'published'
		ORDER BY (rr.id IS NOT NULL), r.created_at DESC`, exchangerID)
	if err != nil {
		return nil, fmt.Errorf("query merchant reviews: %w", err)
	}
	defer rows.Close()

	var out []domain.MerchantReview
	for rows.Next() {
		var mr domain.MerchantReview
		if err := rows.Scan(&mr.ID, &mr.Author, &mr.Rating, &mr.Title, &mr.Body, &mr.CreatedAt,
			&mr.Reply, &mr.ReplyAt); err != nil {
			return nil, fmt.Errorf("scan merchant review: %w", err)
		}
		out = append(out, mr)
	}
	return out, rows.Err()
}

// ReviewBelongsToExchanger reports whether a review is about the given exchanger
// (guards reply writes so a merchant can only answer their own reviews).
func (s *Store) ReviewBelongsToExchanger(ctx context.Context, reviewID, exchangerID string) (bool, error) {
	var ok bool
	err := s.db.QueryRow(ctx, `
		SELECT EXISTS (SELECT 1 FROM reviews WHERE id = $1 AND exchanger_id = $2)`,
		reviewID, exchangerID).Scan(&ok)
	if err != nil {
		return false, fmt.Errorf("review ownership: %w", err)
	}
	return ok, nil
}

// UpsertReviewReply creates or replaces the exchanger's single reply to a review.
func (s *Store) UpsertReviewReply(ctx context.Context, reviewID, exchangerUserID, body string) error {
	_, err := s.db.Exec(ctx, `
		INSERT INTO review_replies (review_id, exchanger_user_id, body, status)
		VALUES ($1, $2, $3, 'published')
		ON CONFLICT (review_id) DO UPDATE
		SET body = EXCLUDED.body, exchanger_user_id = EXCLUDED.exchanger_user_id,
		    status = 'published', updated_at = now()`, reviewID, exchangerUserID, body)
	if err != nil {
		return fmt.Errorf("upsert review reply: %w", err)
	}
	return nil
}

// TrafficSeries returns daily clickout totals for the last `days` days
// (oldest → newest), zero-filled so the chart has one bar per day.
func (s *Store) TrafficSeries(ctx context.Context, exchangerID string, days int) ([]domain.TrafficPoint, error) {
	rows, err := s.db.Query(ctx, `
		SELECT gs::date AS day, COALESCE(c.n, 0)
		FROM generate_series(date_trunc('day', now()) - make_interval(days => $2 - 1),
		                     date_trunc('day', now()), interval '1 day') gs
		LEFT JOIN (
			SELECT date_trunc('day', created_at) AS day, COUNT(*) AS n
			FROM clickouts
			WHERE exchanger_id = $1 AND created_at >= date_trunc('day', now()) - make_interval(days => $2 - 1)
			GROUP BY 1
		) c ON c.day = gs
		ORDER BY gs`, exchangerID, days)
	if err != nil {
		return nil, fmt.Errorf("query traffic series: %w", err)
	}
	defer rows.Close()

	var out []domain.TrafficPoint
	for rows.Next() {
		var p domain.TrafficPoint
		if err := rows.Scan(&p.Day, &p.Clicks); err != nil {
			return nil, fmt.Errorf("scan traffic point: %w", err)
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// TrafficByDirection returns clickout counts per direction over the last `days`.
func (s *Store) TrafficByDirection(ctx context.Context, exchangerID string, days int) ([]domain.TrafficDirection, error) {
	rows, err := s.db.Query(ctx, `
		SELECT d.slug, fc.code, tc.code, COUNT(c.id)
		FROM clickouts c
		JOIN directions d ON d.id = c.direction_id
		JOIN currencies fc ON fc.id = d.from_currency_id
		JOIN currencies tc ON tc.id = d.to_currency_id
		WHERE c.exchanger_id = $1
		  AND c.created_at >= date_trunc('day', now()) - make_interval(days => $2 - 1)
		GROUP BY d.slug, fc.code, tc.code
		ORDER BY COUNT(c.id) DESC`, exchangerID, days)
	if err != nil {
		return nil, fmt.Errorf("query traffic by direction: %w", err)
	}
	defer rows.Close()

	var out []domain.TrafficDirection
	for rows.Next() {
		var td domain.TrafficDirection
		if err := rows.Scan(&td.DirectionSlug, &td.FromCode, &td.ToCode, &td.Clicks); err != nil {
			return nil, fmt.Errorf("scan traffic direction: %w", err)
		}
		out = append(out, td)
	}
	return out, rows.Err()
}

// UpdateExchangerProfile writes the merchant-editable profile fields and sends
// the card back to moderation (is_verified = false).
func (s *Store) UpdateExchangerProfile(ctx context.Context, exchangerID, name, description, website string) error {
	_, err := s.db.Exec(ctx, `
		UPDATE exchangers
		SET name = $2,
		    description = NULLIF($3, ''),
		    website_url = NULLIF($4, ''),
		    is_verified = false
		WHERE id = $1`, exchangerID, name, description, website)
	if err != nil {
		return fmt.Errorf("update exchanger profile: %w", err)
	}
	return nil
}

// MerchantComplaints lists reports filed against the exchanger's reviews.
func (s *Store) MerchantComplaints(ctx context.Context, exchangerID string) ([]domain.MerchantComplaint, error) {
	rows, err := s.db.Query(ctx, `
		SELECT rp.id::text, rp.review_id::text, rp.reason, rp.details, rp.status, rp.created_at,
		       r.body, COALESCE(r.author_name, 'Аноним'), r.rating
		FROM review_reports rp
		JOIN reviews r ON r.id = rp.review_id
		WHERE r.exchanger_id = $1
		ORDER BY (rp.status = 'open') DESC, rp.created_at DESC`, exchangerID)
	if err != nil {
		return nil, fmt.Errorf("query merchant complaints: %w", err)
	}
	defer rows.Close()

	var out []domain.MerchantComplaint
	for rows.Next() {
		var c domain.MerchantComplaint
		if err := rows.Scan(&c.ID, &c.ReviewID, &c.Reason, &c.Details, &c.Status, &c.CreatedAt,
			&c.ReviewBody, &c.Author, &c.Rating); err != nil {
			return nil, fmt.Errorf("scan merchant complaint: %w", err)
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// MerchantPayouts lists the exchanger's settlement history, newest first.
func (s *Store) MerchantPayouts(ctx context.Context, exchangerID string) ([]domain.Payout, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id::text, period_start, period_end, clicks_count, amount::text, currency, status, created_at, paid_at
		FROM partner_payouts
		WHERE exchanger_id = $1
		ORDER BY period_start DESC`, exchangerID)
	if err != nil {
		return nil, fmt.Errorf("query merchant payouts: %w", err)
	}
	defer rows.Close()

	var out []domain.Payout
	for rows.Next() {
		var p domain.Payout
		if err := rows.Scan(&p.ID, &p.PeriodStart, &p.PeriodEnd, &p.ClicksCount, &p.Amount,
			&p.Currency, &p.Status, &p.CreatedAt, &p.PaidAt); err != nil {
			return nil, fmt.Errorf("scan payout: %w", err)
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// ClicksThisMonth counts clickouts since the start of the current month — used
// to estimate the running (unbilled) balance in the cabinet.
func (s *Store) ClicksThisMonth(ctx context.Context, exchangerID string) (int, error) {
	var n int
	err := s.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM clickouts
		WHERE exchanger_id = $1 AND created_at >= date_trunc('month', now())`, exchangerID).Scan(&n)
	if err != nil {
		return 0, fmt.Errorf("clicks this month: %w", err)
	}
	return n, nil
}
