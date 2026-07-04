package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/ANTON-IVANOVICH/kurso-v0/kurso-api/internal/domain"
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
