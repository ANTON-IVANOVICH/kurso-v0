package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ANTON-IVANOVICH/kurso-v0/kurso-api/internal/domain"
)

const merchantKey ctxKey = "merchant"

// The merchant refresh cookie is scoped to the partner auth path — it is only
// ever sent to refresh/logout, never to the data endpoints (which use the
// in-memory access token). httpOnly, so JS can't read it.
const merchantRefreshCookieName = "kurso_partner_rt"
const merchantRefreshCookiePath = "/partner/auth"

type merchantDTO struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Role          string `json:"role"`
	ExchangerSlug string `json:"exchangerSlug"`
	ExchangerName string `json:"exchangerName"`
}

type merchantTokenResponse struct {
	Token    string      `json:"token"`
	Merchant merchantDTO `json:"merchant"`
}

func (a *api) setMerchantCookie(w http.ResponseWriter, token string) {
	a.setAuthCookie(w, merchantRefreshCookieName, merchantRefreshCookiePath, token, a.deps.PartnerAuth.RefreshTTL())
}
func (a *api) clearMerchantCookie(w http.ResponseWriter) {
	a.clearAuthCookie(w, merchantRefreshCookieName, merchantRefreshCookiePath)
}

// merchantFromID resolves the exchanger a representative belongs to and builds
// the response DTO. Returns ok=false if the account has vanished.
func (a *api) merchantFromID(ctx context.Context, userID string) (domain.MerchantIdentity, merchantDTO, bool) {
	id, err := a.deps.Store.MerchantIdentity(ctx, userID)
	if err != nil {
		return domain.MerchantIdentity{}, merchantDTO{}, false
	}
	return id, merchantDTO{
		ID: id.UserID, Email: id.Email, Role: id.Role,
		ExchangerSlug: id.ExchangerSlug, ExchangerName: id.ExchangerName,
	}, true
}

// partnerLogin authenticates a merchant representative and returns an access
// token (refresh in the httpOnly cookie).
func (a *api) partnerLogin(w http.ResponseWriter, r *http.Request) {
	var req adminLoginRequest // {email, password, otp} — same shape as admin
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "invalid JSON body")
		return
	}
	email := strings.ToLower(strings.TrimSpace(req.Email))
	tokens, acct, err := a.deps.PartnerAuth.Login(r.Context(), email, req.Password, req.OTP)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid_credentials", "Неверный email или пароль")
		return
	}
	_, dto, ok := a.merchantFromID(r.Context(), acct.ID)
	if !ok {
		writeError(w, http.StatusUnauthorized, "invalid_credentials", "Кабинет не найден")
		return
	}
	a.setMerchantCookie(w, tokens.Refresh)
	writeJSON(w, http.StatusOK, merchantTokenResponse{Token: tokens.Access, Merchant: dto})
}

// partnerRefresh exchanges the refresh cookie for a new access token.
func (a *api) partnerRefresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(merchantRefreshCookieName)
	if err != nil || cookie.Value == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized", "Сессия не найдена")
		return
	}
	tokens, claims, err := a.deps.PartnerAuth.Refresh(cookie.Value)
	if err != nil {
		a.clearMerchantCookie(w)
		writeError(w, http.StatusUnauthorized, "unauthorized", "Сессия истекла, войдите снова")
		return
	}
	_, dto, ok := a.merchantFromID(r.Context(), claims.Subject)
	if !ok {
		a.clearMerchantCookie(w)
		writeError(w, http.StatusUnauthorized, "unauthorized", "Кабинет не найден")
		return
	}
	a.setMerchantCookie(w, tokens.Refresh)
	writeJSON(w, http.StatusOK, merchantTokenResponse{Token: tokens.Access, Merchant: dto})
}

// partnerLogout clears the refresh cookie.
func (a *api) partnerLogout(w http.ResponseWriter, _ *http.Request) {
	a.clearMerchantCookie(w)
	w.WriteHeader(http.StatusNoContent)
}

// partnerMe returns the authenticated merchant.
func (a *api) partnerMe(w http.ResponseWriter, r *http.Request) {
	m := merchantFromContext(r.Context())
	if m == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized", "Требуется вход")
		return
	}
	writeJSON(w, http.StatusOK, merchantDTO{
		ID: m.UserID, Email: m.Email, Role: m.Role,
		ExchangerSlug: m.ExchangerSlug, ExchangerName: m.ExchangerName,
	})
}

// requireMerchant gates the cabinet endpoints: it validates the access token,
// resolves the representative's exchanger, and stashes the identity in context.
func (a *api) requireMerchant(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const prefix = "Bearer "
		authz := r.Header.Get("Authorization")
		if !strings.HasPrefix(authz, prefix) {
			writeError(w, http.StatusUnauthorized, "unauthorized", "Требуется вход")
			return
		}
		claims, err := a.deps.PartnerAuth.ParseAccess(strings.TrimPrefix(authz, prefix))
		if err != nil {
			writeError(w, http.StatusUnauthorized, "unauthorized", "Сессия истекла, войдите снова")
			return
		}
		id, err := a.deps.Store.MerchantIdentity(r.Context(), claims.Subject)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "unauthorized", "Кабинет не найден")
			return
		}
		ctx := context.WithValue(r.Context(), merchantKey, &id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func merchantFromContext(ctx context.Context) *domain.MerchantIdentity {
	m, _ := ctx.Value(merchantKey).(*domain.MerchantIdentity)
	return m
}
