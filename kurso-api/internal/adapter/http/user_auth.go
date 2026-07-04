package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/ANTON-IVANOVICH/kurso-v0/kurso-api/internal/adapter/store"
	"github.com/ANTON-IVANOVICH/kurso-v0/kurso-api/internal/service/auth"
)

const userClaimsKey ctxKey = "userClaims"

// The end-user refresh cookie uses path "/" so it is sent with the Nuxt SSR
// page request (which forwards it to /session to resolve the user server-side),
// not just the auth endpoints. It stays httpOnly, so JS still can't read it.
const userRefreshCookieName = "kurso_rt"
const userRefreshCookiePath = "/"

type userDTO struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type userTokenResponse struct {
	Token string  `json:"token"`
	User  userDTO `json:"user"`
}

// nameFromEmail derives a display name when none is set.
func nameFromEmail(email string) string {
	if i := strings.IndexByte(email, '@'); i > 0 {
		return email[:i]
	}
	return email
}

func (a *api) setUserCookie(w http.ResponseWriter, token string) {
	a.setAuthCookie(w, userRefreshCookieName, userRefreshCookiePath, token, a.deps.UserAuth.RefreshTTL())
}
func (a *api) clearUserCookie(w http.ResponseWriter) {
	a.clearAuthCookie(w, userRefreshCookieName, userRefreshCookiePath)
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// userRegister creates a password account, logs it in, and returns an access token.
func (a *api) userRegister(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "invalid JSON body")
		return
	}
	email := strings.ToLower(strings.TrimSpace(req.Email))
	if !strings.Contains(email, "@") {
		writeError(w, http.StatusBadRequest, "invalid_email", "Введите корректный email")
		return
	}
	if len(req.Password) < 8 {
		writeError(w, http.StatusBadRequest, "weak_password", "Пароль минимум 8 символов")
		return
	}
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		a.serverError(w, "hash password", err)
		return
	}
	var name *string
	if n := strings.TrimSpace(req.Name); n != "" {
		name = &n
	}
	u, err := a.deps.Store.CreateUser(r.Context(), email, hash, name)
	if err != nil {
		if errors.Is(err, store.ErrEmailTaken) {
			writeError(w, http.StatusConflict, "email_taken", "Этот email уже зарегистрирован")
			return
		}
		a.serverError(w, "create user", err)
		return
	}

	// Affiliate attribution: if the visitor arrived via a partner link, credit the
	// referrer (base code before any ".tag" suffix). Best-effort.
	if c, cerr := r.Cookie("kurso_ref"); cerr == nil && c.Value != "" {
		baseCode := c.Value
		if i := strings.IndexByte(baseCode, '.'); i > 0 {
			baseCode = baseCode[:i]
		}
		if refID, ok, rerr := a.deps.Store.UserIDByReferralCode(r.Context(), baseCode); rerr == nil && ok {
			_ = a.deps.Store.SetReferrer(r.Context(), u.ID, refID)
		}
	}
	_, _ = a.deps.Store.EnsureReferralCode(r.Context(), u.ID)
	tokens, err := a.deps.UserAuth.IssueFor(auth.Account{ID: u.ID, Email: u.Email, Role: "user"})
	if err != nil {
		a.serverError(w, "issue tokens", err)
		return
	}
	a.setUserCookie(w, tokens.Refresh)
	writeJSON(w, http.StatusCreated, userTokenResponse{Token: tokens.Access, User: userFrom(u.ID, u.Email, u.DisplayName)})
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// userLogin verifies credentials and returns an access token (refresh in cookie).
func (a *api) userLogin(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "invalid JSON body")
		return
	}
	email := strings.ToLower(strings.TrimSpace(req.Email))
	tokens, acct, err := a.deps.UserAuth.Login(r.Context(), email, req.Password, "")
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid_credentials", "Неверный email или пароль")
		return
	}
	a.setUserCookie(w, tokens.Refresh)
	writeJSON(w, http.StatusOK, userTokenResponse{Token: tokens.Access, User: userFrom(acct.ID, acct.Email, nil)})
}

// userRefresh rotates the refresh cookie and returns a fresh access token.
func (a *api) userRefresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(userRefreshCookieName)
	if err != nil || cookie.Value == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized", "Сессия не найдена")
		return
	}
	tokens, claims, err := a.deps.UserAuth.Refresh(cookie.Value)
	if err != nil {
		a.clearUserCookie(w)
		writeError(w, http.StatusUnauthorized, "unauthorized", "Сессия истекла")
		return
	}
	a.setUserCookie(w, tokens.Refresh)
	writeJSON(w, http.StatusOK, userTokenResponse{Token: tokens.Access, User: userFrom(claims.Subject, claims.Email, nil)})
}

// userSession validates the refresh cookie WITHOUT rotating it and returns the
// current user. The Nuxt SSR server calls this (forwarding the cookie) so authed
// pages render server-side with the right identity.
func (a *api) userSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(userRefreshCookieName)
	if err != nil || cookie.Value == "" {
		writeJSON(w, http.StatusOK, map[string]any{"user": nil})
		return
	}
	claims, err := a.deps.UserAuth.ParseRefresh(cookie.Value)
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]any{"user": nil})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"user": userFrom(claims.Subject, claims.Email, nil)})
}

// userLogout clears the refresh cookie.
func (a *api) userLogout(w http.ResponseWriter, _ *http.Request) {
	a.clearUserCookie(w)
	w.WriteHeader(http.StatusNoContent)
}

// userMe returns the authenticated user (from the access token).
func (a *api) userMe(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value(userClaimsKey).(*auth.Claims)
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized", "Требуется вход")
		return
	}
	writeJSON(w, http.StatusOK, userFrom(claims.Subject, claims.Email, nil))
}

// requireUser gates endpoints behind a valid user access token.
func (a *api) requireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const prefix = "Bearer "
		authz := r.Header.Get("Authorization")
		if !strings.HasPrefix(authz, prefix) {
			writeError(w, http.StatusUnauthorized, "unauthorized", "Требуется вход")
			return
		}
		claims, err := a.deps.UserAuth.ParseAccess(strings.TrimPrefix(authz, prefix))
		if err != nil {
			writeError(w, http.StatusUnauthorized, "unauthorized", "Сессия истекла")
			return
		}
		ctx := context.WithValue(r.Context(), userClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func userFrom(id, email string, displayName *string) userDTO {
	name := nameFromEmail(email)
	if displayName != nil && strings.TrimSpace(*displayName) != "" {
		name = *displayName
	}
	return userDTO{ID: id, Email: email, Name: name}
}
