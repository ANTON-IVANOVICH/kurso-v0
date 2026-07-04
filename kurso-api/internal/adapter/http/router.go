// Package httpapi is the inbound HTTP adapter: it maps HTTP requests onto the
// application's use cases. Request/response types are shared with the
// frontends through the generated openapi package.
package httpapi

import (
	"log/slog"
	"net/http"
	"slices"
	"time"

	"github.com/ANTON-IVANOVICH/kurso-v0/kurso-api/internal/adapter/store"
	"github.com/ANTON-IVANOVICH/kurso-v0/kurso-api/internal/service/auth"
	"github.com/ANTON-IVANOVICH/kurso-v0/kurso-api/internal/service/rates"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// Deps are the collaborators the HTTP adapter needs.
type Deps struct {
	Log            *slog.Logger
	DB             *pgxpool.Pool
	Redis          *redis.Client
	Svc            *rates.Service
	Auth           *auth.Service // admin identity
	UserAuth       *auth.Service // end-user identity
	Store          *store.Store
	AllowedOrigins []string
	// CookieSecure marks the refresh cookie Secure (production/HTTPS).
	CookieSecure bool
	// CookieDomain scopes the refresh cookie (empty = host-only for dev; set to
	// ".kurso.io" in prod so the Nuxt SSR origin can read it).
	CookieDomain string
}

// api binds handlers to their dependencies.
type api struct {
	deps Deps
}

// NewRouter builds the fully-wired chi router.
func NewRouter(d Deps) http.Handler {
	a := &api{deps: d}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(a.cors)

	// System / probes.
	r.Get("/healthz", a.health)
	r.Get("/readyz", a.readiness)

	// Public API (v1).
	r.Route("/api/v1", func(r chi.Router) {
		// Regular request/response endpoints get a bounded timeout.
		r.Group(func(r chi.Router) {
			r.Use(middleware.Timeout(15 * time.Second))
			r.Get("/currencies", a.listCurrencies)
			r.Get("/directions", a.listDirections)
			r.Get("/exchangers", a.listExchangers)
			r.Get("/exchangers/{slug}", a.getExchanger)
			r.Get("/exchangers/{slug}/reviews", a.listReviews)
			r.Post("/exchangers/{slug}/reviews", a.createReview)
			r.Post("/reviews/{id}/report", a.reportReview)
			r.Get("/rates/{direction}", a.getRates)

			// End-user auth. session is a non-rotating check the Nuxt SSR server
			// calls (forwarding the refresh cookie) to resolve the current user.
			r.Route("/auth", func(r chi.Router) {
				r.Post("/register", a.userRegister)
				r.Post("/login", a.userLogin)
				r.Post("/refresh", a.userRefresh)
				r.Post("/logout", a.userLogout)
				r.Get("/session", a.userSession)
				r.Group(func(r chi.Router) {
					r.Use(a.requireUser)
					r.Get("/me", a.userMe)
				})
			})
		})
		// SSE stream is long-lived — it must not be wrapped by a request timeout.
		r.Get("/rates/{direction}/stream", a.streamRates)
	})

	// Admin API (admin.kurso.io). Login is public; everything else needs a valid
	// admin JWT via RequireAdmin.
	r.Route("/admin", func(r chi.Router) {
		r.Use(middleware.Timeout(15 * time.Second))
		r.Post("/auth/login", a.adminLogin)
		r.Post("/auth/refresh", a.adminRefresh) // reads the httpOnly refresh cookie
		r.Post("/auth/logout", a.adminLogout)
		r.Group(func(r chi.Router) {
			r.Use(a.requireAdmin)
			r.Get("/auth/me", a.adminMe)
			r.Get("/reviews", a.adminListReviews)
			r.Patch("/reviews/{id}", a.adminModerateReview)
			r.Get("/reports", a.adminListReports)
			r.Patch("/reports/{id}", a.adminResolveReport)
		})
	})

	// Outbound click + referral redirect.
	r.Get("/go/{slug}", a.clickout)

	return r
}

// cors reflects allowed origins so the browser frontends can call the API in
// development and production. Empty AllowedOrigins means "reflect any origin".
func (a *api) cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		allowed := len(a.deps.AllowedOrigins) == 0 ||
			slices.Contains(a.deps.AllowedOrigins, "*") ||
			slices.Contains(a.deps.AllowedOrigins, origin)
		if origin != "" && allowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Add("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			// Credentialed so the admin/partner refresh cookie is accepted; the
			// reflected specific origin (never "*") keeps this valid.
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
