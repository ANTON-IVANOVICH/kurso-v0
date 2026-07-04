package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/ANTON-IVANOVICH/kurso-v0/kurso-api/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// ErrEmailTaken is returned by CreateUser when the email already exists.
var ErrEmailTaken = errors.New("email already registered")

// UserByEmail loads a user by email (case-insensitive via citext).
func (s *Store) UserByEmail(ctx context.Context, email string) (domain.User, error) {
	var u domain.User
	err := s.db.QueryRow(ctx, `
		SELECT id::text, email, COALESCE(password_hash, ''), display_name, status
		FROM users
		WHERE email = $1`, email).
		Scan(&u.ID, &u.Email, &u.PasswordHash, &u.DisplayName, &u.Status)
	if err != nil {
		return domain.User{}, fmt.Errorf("query user by email: %w", err)
	}
	return u, nil
}

// CreateUser inserts a password user and returns it. A duplicate email yields
// ErrEmailTaken.
func (s *Store) CreateUser(ctx context.Context, email, passwordHash string, displayName *string) (domain.User, error) {
	var u domain.User
	err := s.db.QueryRow(ctx, `
		INSERT INTO users (email, password_hash, display_name)
		VALUES ($1, $2, $3)
		RETURNING id::text, email, COALESCE(password_hash, ''), display_name, status`,
		email, passwordHash, displayName).
		Scan(&u.ID, &u.Email, &u.PasswordHash, &u.DisplayName, &u.Status)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.User{}, ErrEmailTaken
		}
		return domain.User{}, fmt.Errorf("insert user: %w", err)
	}
	return u, nil
}

// TouchUserLogin bumps the user's updated_at on a successful login.
func (s *Store) TouchUserLogin(ctx context.Context, id string) error {
	_, err := s.db.Exec(ctx, `UPDATE users SET updated_at = now() WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("touch user login: %w", err)
	}
	return nil
}

// EnsureReferralCode returns the user's affiliate code, generating a stable one
// (first 8 hex of the id) if absent.
func (s *Store) EnsureReferralCode(ctx context.Context, userID string) (string, error) {
	var code string
	err := s.db.QueryRow(ctx, `
		UPDATE users
		SET referral_code = COALESCE(referral_code, substr(replace(id::text, '-', ''), 1, 8))
		WHERE id = $1
		RETURNING referral_code`, userID).Scan(&code)
	if err != nil {
		return "", fmt.Errorf("ensure referral code: %w", err)
	}
	return code, nil
}

// UserIDByReferralCode resolves an affiliate base code to the owning user id.
func (s *Store) UserIDByReferralCode(ctx context.Context, code string) (string, bool, error) {
	var id string
	err := s.db.QueryRow(ctx, `SELECT id::text FROM users WHERE referral_code = $1`, code).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", false, nil
	}
	if err != nil {
		return "", false, fmt.Errorf("user by referral code: %w", err)
	}
	return id, true, nil
}

// SetReferrer attributes a newly-registered user to their referrer (no self-refs).
func (s *Store) SetReferrer(ctx context.Context, userID, referrerID string) error {
	if userID == referrerID {
		return nil
	}
	_, err := s.db.Exec(ctx, `UPDATE users SET referred_by = $2 WHERE id = $1 AND referred_by IS NULL`, userID, referrerID)
	if err != nil {
		return fmt.Errorf("set referrer: %w", err)
	}
	return nil
}

// ReferralStats aggregates real affiliate performance for a base code over the
// last `days` days: total clicks (clickouts whose ref_code is code or code.tag),
// a daily series, per-tag breakdown, and referred registrations.
func (s *Store) ReferralStats(ctx context.Context, userID, code string, days int) (domain.ReferralStats, error) {
	var st domain.ReferralStats
	// like-pattern matches the base code and any ".tag" suffix.
	pattern := code + ".%"

	if err := s.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM clickouts
		WHERE (ref_code = $1 OR ref_code LIKE $2)
		  AND created_at >= date_trunc('day', now()) - make_interval(days => $3 - 1)`,
		code, pattern, days).Scan(&st.Clicks); err != nil {
		return st, fmt.Errorf("referral clicks: %w", err)
	}

	if err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE referred_by = $1`, userID).
		Scan(&st.Registrations); err != nil {
		return st, fmt.Errorf("referral registrations: %w", err)
	}

	rows, err := s.db.Query(ctx, `
		SELECT gs::date, COALESCE(c.n, 0)
		FROM generate_series(date_trunc('day', now()) - make_interval(days => $3 - 1),
		                     date_trunc('day', now()), interval '1 day') gs
		LEFT JOIN (
			SELECT date_trunc('day', created_at) AS day, COUNT(*) AS n
			FROM clickouts
			WHERE (ref_code = $1 OR ref_code LIKE $2)
			  AND created_at >= date_trunc('day', now()) - make_interval(days => $3 - 1)
			GROUP BY 1
		) c ON c.day = gs
		ORDER BY gs`, code, pattern, days)
	if err != nil {
		return st, fmt.Errorf("referral series: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var p domain.TrafficPoint
		if err := rows.Scan(&p.Day, &p.Clicks); err != nil {
			return st, fmt.Errorf("scan referral series: %w", err)
		}
		st.Series = append(st.Series, p)
	}
	if err := rows.Err(); err != nil {
		return st, err
	}

	// Per-tag breakdown: the part after "code." (empty tag = the bare link).
	tagRows, err := s.db.Query(ctx, `
		SELECT CASE WHEN ref_code = $1 THEN '' ELSE substr(ref_code, length($1) + 2) END AS tag,
		       COUNT(*)
		FROM clickouts
		WHERE (ref_code = $1 OR ref_code LIKE $2)
		  AND created_at >= date_trunc('day', now()) - make_interval(days => $3 - 1)
		GROUP BY 1
		ORDER BY COUNT(*) DESC`, code, pattern, days)
	if err != nil {
		return st, fmt.Errorf("referral tags: %w", err)
	}
	defer tagRows.Close()
	for tagRows.Next() {
		var t domain.ReferralTag
		if err := tagRows.Scan(&t.Tag, &t.Clicks); err != nil {
			return st, fmt.Errorf("scan referral tag: %w", err)
		}
		st.ByTag = append(st.ByTag, t)
	}
	return st, tagRows.Err()
}
