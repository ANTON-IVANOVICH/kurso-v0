// Package store is the driven adapter that persists and reads the catalogue and
// rates in Postgres. It maps rows onto domain types and hides all SQL.
package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/ANTON-IVANOVICH/kurso-v0/kurso-api/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store reads and writes the Kurso catalogue and rate tables.
type Store struct {
	db *pgxpool.Pool
}

// New builds a Store over an existing pgx pool.
func New(db *pgxpool.Pool) *Store { return &Store{db: db} }

// Currencies returns the active currency catalogue in display order.
func (s *Store) Currencies(ctx context.Context) ([]domain.Currency, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id::text, code, name, kind, network, icon_url
		FROM currencies
		WHERE is_active
		ORDER BY sort_order, code`)
	if err != nil {
		return nil, fmt.Errorf("query currencies: %w", err)
	}
	defer rows.Close()

	var out []domain.Currency
	for rows.Next() {
		var c domain.Currency
		if err := rows.Scan(&c.ID, &c.Code, &c.Name, &c.Kind, &c.Network, &c.IconURL); err != nil {
			return nil, fmt.Errorf("scan currency: %w", err)
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

const directionSelect = `
	SELECT d.id::text, d.slug, d.from_currency_id::text, d.to_currency_id::text,
	       fc.code, fc.name, tc.code, tc.name, d.is_popular
	FROM directions d
	JOIN currencies fc ON fc.id = d.from_currency_id
	JOIN currencies tc ON tc.id = d.to_currency_id`

func scanDirection(row pgx.Row) (domain.Direction, error) {
	var d domain.Direction
	err := row.Scan(&d.ID, &d.Slug, &d.FromID, &d.ToID,
		&d.FromCode, &d.FromName, &d.ToCode, &d.ToName, &d.IsPopular)
	return d, err
}

// Directions returns active directions, popular ones first.
func (s *Store) Directions(ctx context.Context) ([]domain.Direction, error) {
	rows, err := s.db.Query(ctx, directionSelect+`
		WHERE d.is_active
		ORDER BY d.is_popular DESC, d.sort_order, d.slug`)
	if err != nil {
		return nil, fmt.Errorf("query directions: %w", err)
	}
	defer rows.Close()

	var out []domain.Direction
	for rows.Next() {
		d, err := scanDirection(rows)
		if err != nil {
			return nil, fmt.Errorf("scan direction: %w", err)
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

// DirectionBySlug looks up a single direction. ok is false when not found.
func (s *Store) DirectionBySlug(ctx context.Context, slug string) (domain.Direction, bool, error) {
	d, err := scanDirection(s.db.QueryRow(ctx, directionSelect+` WHERE d.slug = $1`, slug))
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Direction{}, false, nil
	}
	if err != nil {
		return domain.Direction{}, false, fmt.Errorf("query direction %q: %w", slug, err)
	}
	return d, true, nil
}

const exchangerSelect = `
	SELECT e.id::text, e.slug, e.name, e.status, e.website_url, e.referral_url_template, e.logo_url,
	       e.description, e.rating_avg::float8, e.reviews_count, e.is_verified,
	       (e.referral_url_template IS NOT NULL) AS partner,
	       (SELECT SUM(r.reserve) FROM rates r WHERE r.exchanger_id = e.id AND r.is_active)::text AS reserve_total,
	       (SELECT COUNT(DISTINCT r.direction_id) FROM rates r WHERE r.exchanger_id = e.id AND r.is_active) AS directions_count,
	       COALESCE((SELECT array_agg(DISTINCT c.code ORDER BY c.code)
	                 FROM rates r
	                 JOIN directions d ON d.id = r.direction_id
	                 JOIN currencies c ON c.id = d.from_currency_id
	                 WHERE r.exchanger_id = e.id AND r.is_active), ARRAY[]::text[]) AS assets,
	       EXTRACT(YEAR FROM e.created_at)::int AS on_since
	FROM exchangers e`

func scanExchanger(row pgx.Row) (domain.Exchanger, error) {
	var e domain.Exchanger
	err := row.Scan(&e.ID, &e.Slug, &e.Name, &e.Status, &e.WebsiteURL, &e.ReferralTmpl,
		&e.LogoURL, &e.Description, &e.RatingAvg, &e.ReviewsCount, &e.IsVerified, &e.Partner,
		&e.ReserveTotal, &e.DirectionsCount, &e.Assets, &e.OnSinceYear)
	return e, err
}

// Exchangers returns active exchangers ordered by name.
func (s *Store) Exchangers(ctx context.Context) ([]domain.Exchanger, error) {
	rows, err := s.db.Query(ctx, exchangerSelect+` WHERE status = 'active' ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("query exchangers: %w", err)
	}
	defer rows.Close()

	var out []domain.Exchanger
	for rows.Next() {
		e, err := scanExchanger(rows)
		if err != nil {
			return nil, fmt.Errorf("scan exchanger: %w", err)
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

// ExchangerBySlug looks up a single exchanger. ok is false when not found.
func (s *Store) ExchangerBySlug(ctx context.Context, slug string) (domain.Exchanger, bool, error) {
	e, err := scanExchanger(s.db.QueryRow(ctx, exchangerSelect+` WHERE slug = $1`, slug))
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Exchanger{}, false, nil
	}
	if err != nil {
		return domain.Exchanger{}, false, fmt.Errorf("query exchanger %q: %w", slug, err)
	}
	return e, true, nil
}

// RatesByDirection returns active rates for a direction, best (highest) first.
func (s *Store) RatesByDirection(ctx context.Context, directionID string) ([]domain.RateRow, error) {
	rows, err := s.db.Query(ctx, `
		SELECT e.id::text, e.slug, e.name, (e.referral_url_template IS NOT NULL),
		       e.rating_avg::float8, e.reviews_count,
		       r.rate::text, r.reserve::text, r.min_amount::text, r.max_amount::text, r.fetched_at
		FROM rates r
		JOIN exchangers e ON e.id = r.exchanger_id
		WHERE r.direction_id = $1 AND r.is_active AND e.status = 'active'
		ORDER BY r.rate DESC`, directionID)
	if err != nil {
		return nil, fmt.Errorf("query rates: %w", err)
	}
	defer rows.Close()

	var out []domain.RateRow
	for rows.Next() {
		var r domain.RateRow
		if err := rows.Scan(&r.ExchangerID, &r.ExchangerSlug, &r.ExchangerName, &r.Partner,
			&r.Rating, &r.ReviewsCount, &r.Rate, &r.Reserve, &r.MinAmount, &r.MaxAmount,
			&r.FetchedAt); err != nil {
			return nil, fmt.Errorf("scan rate: %w", err)
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// UpsertRate writes the current rate for an (exchanger, direction) and appends a
// history row. Used by the rate runner.
func (s *Store) UpsertRate(ctx context.Context, exchangerID, directionID, rate, reserve string) error {
	if _, err := s.db.Exec(ctx, `
		INSERT INTO rates (exchanger_id, direction_id, rate, reserve, fetched_at, is_active)
		VALUES ($1, $2, $3::numeric, NULLIF($4,'')::numeric, now(), true)
		ON CONFLICT (exchanger_id, direction_id) DO UPDATE
		SET rate = EXCLUDED.rate, reserve = EXCLUDED.reserve, fetched_at = now(), is_active = true`,
		exchangerID, directionID, rate, reserve); err != nil {
		return fmt.Errorf("upsert rate: %w", err)
	}
	if _, err := s.db.Exec(ctx, `
		INSERT INTO rates_history (exchanger_id, direction_id, rate, reserve)
		VALUES ($1, $2, $3::numeric, NULLIF($4,'')::numeric)`,
		exchangerID, directionID, rate, reserve); err != nil {
		return fmt.Errorf("append rate history: %w", err)
	}
	return nil
}
