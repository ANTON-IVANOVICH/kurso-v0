package store

import (
	"context"
	"fmt"

	"github.com/ANTON-IVANOVICH/kurso-v0/kurso-api/internal/domain"
)

// AdminByEmail loads an administrator by email (case-insensitive via citext).
func (s *Store) AdminByEmail(ctx context.Context, email string) (domain.Admin, error) {
	var a domain.Admin
	var role string
	err := s.db.QueryRow(ctx, `
		SELECT id::text, email, password_hash, role, totp_secret, totp_enabled, status
		FROM admins
		WHERE email = $1`, email).
		Scan(&a.ID, &a.Email, &a.PasswordHash, &role, &a.TOTPSecret, &a.TOTPEnabled, &a.Status)
	if err != nil {
		return domain.Admin{}, fmt.Errorf("query admin by email: %w", err)
	}
	a.Role = domain.AdminRole(role)
	return a, nil
}

// TouchAdminLogin stamps the admin's last successful login time.
func (s *Store) TouchAdminLogin(ctx context.Context, id string) error {
	_, err := s.db.Exec(ctx, `UPDATE admins SET last_login_at = now() WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("touch admin login: %w", err)
	}
	return nil
}

// EnsureAdmin inserts a seed administrator if the email is not already present.
// Existing accounts (and any changed password) are left untouched.
func (s *Store) EnsureAdmin(ctx context.Context, email, passwordHash, role string) error {
	_, err := s.db.Exec(ctx, `
		INSERT INTO admins (email, password_hash, role)
		VALUES ($1, $2, $3)
		ON CONFLICT (email) DO NOTHING`, email, passwordHash, role)
	if err != nil {
		return fmt.Errorf("ensure admin: %w", err)
	}
	return nil
}
