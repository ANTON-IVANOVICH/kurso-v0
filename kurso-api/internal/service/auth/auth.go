// Package auth is the admin authentication use case: password verification
// (bcrypt) and a short-lived access / long-lived refresh JWT pair (HS256). The
// access token is returned to the client (kept in memory); the refresh token is
// delivered as an httpOnly cookie by the HTTP adapter and exchanged for fresh
// access tokens here. It depends only on an AdminRepo port.
package auth

import (
	"context"
	"errors"
	"time"

	"github.com/ANTON-IVANOVICH/kurso-v0/kurso-api/internal/domain"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// ErrInvalidCredentials is returned for any login failure — wrong email,
// wrong password, or a disabled account — so callers can't distinguish them.
var ErrInvalidCredentials = errors.New("invalid credentials")

// ErrInvalidToken is returned when a refresh token is missing, malformed,
// expired, or of the wrong kind.
var ErrInvalidToken = errors.New("invalid token")

const (
	typeAccess  = "access"
	typeRefresh = "refresh"
)

// AdminRepo is the storage port the auth service needs.
type AdminRepo interface {
	AdminByEmail(ctx context.Context, email string) (domain.Admin, error)
	TouchAdminLogin(ctx context.Context, id string) error
}

// Claims is the JWT payload for an admin token.
type Claims struct {
	jwt.RegisteredClaims
	Typ   string `json:"typ"`
	Role  string `json:"role"`
	Email string `json:"email"`
}

// Tokens is an issued access/refresh pair.
type Tokens struct {
	Access  string
	Refresh string
}

// Service issues and validates admin tokens.
type Service struct {
	repo       AdminRepo
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

// NewService builds the auth service.
func NewService(repo AdminRepo, secret string, accessTTL, refreshTTL time.Duration) *Service {
	return &Service{repo: repo, secret: []byte(secret), accessTTL: accessTTL, refreshTTL: refreshTTL}
}

// RefreshTTL exposes the refresh lifetime so the HTTP adapter can size the cookie.
func (s *Service) RefreshTTL() time.Duration { return s.refreshTTL }

// HashPassword returns a bcrypt hash suitable for admins.password_hash.
func HashPassword(password string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(h), err
}

// Login verifies credentials and returns a fresh token pair plus the admin.
// `otp` is accepted for forward-compatibility with TOTP 2FA (seeded admins have
// it disabled).
func (s *Service) Login(ctx context.Context, email, password, _ string) (Tokens, domain.Admin, error) {
	a, err := s.repo.AdminByEmail(ctx, email)
	if err != nil {
		_ = bcrypt.CompareHashAndPassword(
			[]byte("$2a$10$C6UzMDM.H6dfI/f/IKcEeO3fL7lHnJh5m9Q3qYcU8b6b3W2t3k2eK"),
			[]byte(password),
		)
		return Tokens{}, domain.Admin{}, ErrInvalidCredentials
	}
	if a.Status != "active" {
		return Tokens{}, domain.Admin{}, ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(a.PasswordHash), []byte(password)); err != nil {
		return Tokens{}, domain.Admin{}, ErrInvalidCredentials
	}

	tokens, err := s.issue(a)
	if err != nil {
		return Tokens{}, domain.Admin{}, err
	}
	_ = s.repo.TouchAdminLogin(ctx, a.ID)
	return tokens, a, nil
}

// Refresh validates a refresh token and mints a rotated token pair. The returned
// claims describe the admin (id/email/role) for a /me-style response.
func (s *Service) Refresh(refreshToken string) (Tokens, *Claims, error) {
	claims, err := s.parse(refreshToken)
	if err != nil || claims.Typ != typeRefresh {
		return Tokens{}, nil, ErrInvalidToken
	}
	a := domain.Admin{ID: claims.Subject, Email: claims.Email, Role: domain.AdminRole(claims.Role)}
	tokens, err := s.issue(a)
	if err != nil {
		return Tokens{}, nil, err
	}
	return tokens, claims, nil
}

// ParseAccess validates an access token (rejecting refresh tokens).
func (s *Service) ParseAccess(token string) (*Claims, error) {
	claims, err := s.parse(token)
	if err != nil {
		return nil, err
	}
	if claims.Typ != typeAccess {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

func (s *Service) issue(a domain.Admin) (Tokens, error) {
	access, err := s.sign(a, typeAccess, s.accessTTL)
	if err != nil {
		return Tokens{}, err
	}
	refresh, err := s.sign(a, typeRefresh, s.refreshTTL)
	if err != nil {
		return Tokens{}, err
	}
	return Tokens{Access: access, Refresh: refresh}, nil
}

func (s *Service) sign(a domain.Admin, typ string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   a.ID,
			Issuer:    "kurso-admin",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
		Typ:   typ,
		Role:  string(a.Role),
		Email: a.Email,
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secret)
}

func (s *Service) parse(token string) (*Claims, error) {
	parsed, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := parsed.Claims.(*Claims)
	if !ok || !parsed.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
