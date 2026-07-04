package httpapi

import (
	"net/http"
	"time"
)

// setAuthCookie writes an httpOnly refresh cookie scoped to a path so it is only
// sent to that auth group's endpoints (admin vs user), never to normal requests.
func (a *api) setAuthCookie(w http.ResponseWriter, name, path, token string, ttl time.Duration) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    token,
		Path:     path,
		Domain:   a.deps.CookieDomain,
		HttpOnly: true,
		Secure:   a.deps.CookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(ttl.Seconds()),
	})
}

// clearAuthCookie expires a refresh cookie.
func (a *api) clearAuthCookie(w http.ResponseWriter, name, path string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     path,
		Domain:   a.deps.CookieDomain,
		HttpOnly: true,
		Secure:   a.deps.CookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}
