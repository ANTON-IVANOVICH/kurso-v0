package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ANTON-IVANOVICH/kurso-v0/kurso-api/internal/service/auth"
)

type ctxKey string

const adminClaimsKey ctxKey = "adminClaims"

// refreshCookieName is the httpOnly cookie carrying the refresh token. It is
// scoped to the auth path so it is only ever sent to the refresh/logout
// endpoints, never to normal API calls.
const refreshCookieName = "kurso_admin_rt"
const refreshCookiePath = "/admin/auth"

type adminLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	OTP      string `json:"otp"`
}

type adminDTO struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// tokenResponse returns the ACCESS token only. The refresh token never reaches
// JavaScript — it rides in the httpOnly cookie.
type tokenResponse struct {
	Token string   `json:"token"`
	Admin adminDTO `json:"admin"`
}

func (a *api) setRefreshCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     refreshCookieName,
		Value:    token,
		Path:     refreshCookiePath,
		HttpOnly: true,
		Secure:   a.deps.CookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(a.deps.Auth.RefreshTTL().Seconds()),
	})
}

func (a *api) clearRefreshCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     refreshCookieName,
		Value:    "",
		Path:     refreshCookiePath,
		HttpOnly: true,
		Secure:   a.deps.CookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}

// adminLogin authenticates an administrator, sets the refresh cookie, and
// returns a short-lived access token in the body.
func (a *api) adminLogin(w http.ResponseWriter, r *http.Request) {
	var req adminLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "invalid JSON body")
		return
	}
	tokens, admin, err := a.deps.Auth.Login(r.Context(), strings.TrimSpace(req.Email), req.Password, req.OTP)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid_credentials", "Неверный email или пароль")
		return
	}
	a.setRefreshCookie(w, tokens.Refresh)
	writeJSON(w, http.StatusOK, tokenResponse{
		Token: tokens.Access,
		Admin: adminDTO{ID: admin.ID, Email: admin.Email, Role: string(admin.Role)},
	})
}

// adminRefresh exchanges the refresh cookie for a new access token (and rotates
// the refresh cookie). Public — authentication is the cookie itself.
func (a *api) adminRefresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(refreshCookieName)
	if err != nil || cookie.Value == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized", "Сессия не найдена")
		return
	}
	tokens, claims, err := a.deps.Auth.Refresh(cookie.Value)
	if err != nil {
		a.clearRefreshCookie(w)
		writeError(w, http.StatusUnauthorized, "unauthorized", "Сессия истекла, войдите снова")
		return
	}
	a.setRefreshCookie(w, tokens.Refresh)
	writeJSON(w, http.StatusOK, tokenResponse{
		Token: tokens.Access,
		Admin: adminDTO{ID: claims.Subject, Email: claims.Email, Role: claims.Role},
	})
}

// adminLogout clears the refresh cookie server-side.
func (a *api) adminLogout(w http.ResponseWriter, _ *http.Request) {
	a.clearRefreshCookie(w)
	w.WriteHeader(http.StatusNoContent)
}

// adminMe returns the currently-authenticated administrator (from the access JWT).
func (a *api) adminMe(w http.ResponseWriter, r *http.Request) {
	claims := adminFromContext(r.Context())
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized", "Требуется вход")
		return
	}
	writeJSON(w, http.StatusOK, adminDTO{ID: claims.Subject, Email: claims.Email, Role: claims.Role})
}

// requireAdmin rejects requests without a valid admin access token and stashes
// the parsed claims in the request context.
func (a *api) requireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const prefix = "Bearer "
		authz := r.Header.Get("Authorization")
		if !strings.HasPrefix(authz, prefix) {
			writeError(w, http.StatusUnauthorized, "unauthorized", "Требуется вход")
			return
		}
		claims, err := a.deps.Auth.ParseAccess(strings.TrimPrefix(authz, prefix))
		if err != nil {
			writeError(w, http.StatusUnauthorized, "unauthorized", "Сессия истекла, войдите снова")
			return
		}
		ctx := context.WithValue(r.Context(), adminClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func adminFromContext(ctx context.Context) *auth.Claims {
	c, _ := ctx.Value(adminClaimsKey).(*auth.Claims)
	return c
}
