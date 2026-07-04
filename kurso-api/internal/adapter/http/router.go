// Package httpapi is the inbound HTTP adapter: it maps HTTP requests onto the
// application's use cases. Request/response types are shared with the
// frontends through the generated openapi package.
package httpapi

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// Deps are the collaborators the HTTP adapter needs.
type Deps struct {
	Log   *slog.Logger
	DB    *pgxpool.Pool
	Redis *redis.Client
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
	r.Use(middleware.Timeout(30 * time.Second))

	// System / probes.
	r.Get("/healthz", a.health)
	r.Get("/readyz", a.readiness)

	// Public API (v1). Endpoints are stubbed at Stage 0 and filled in at Stage 1.
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/currencies", a.listCurrencies)
	})

	return r
}
