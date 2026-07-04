package httpapi

import (
	"context"
	"net/http"
	"time"

	"github.com/ANTON-IVANOVICH/kurso-v0/kurso-api/internal/adapter/http/openapi"
)

// health is the liveness probe: it only proves the process is serving.
func (a *api) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// readiness is the readiness probe: it verifies that dependencies (Postgres,
// Redis) are reachable and returns 503 if any is down.
func (a *api) readiness(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	checks := map[string]string{}
	ready := true

	if err := a.deps.DB.Ping(ctx); err != nil {
		checks["postgres"] = "error: " + err.Error()
		ready = false
	} else {
		checks["postgres"] = "ok"
	}

	if err := a.deps.Redis.Ping(ctx).Err(); err != nil {
		checks["redis"] = "error: " + err.Error()
		ready = false
	} else {
		checks["redis"] = "ok"
	}

	status := "ok"
	code := http.StatusOK
	if !ready {
		status = "degraded"
		code = http.StatusServiceUnavailable
	}

	writeJSON(w, code, map[string]any{"status": status, "checks": checks})
}

// listCurrencies returns the currency catalogue. Stage 1 backs this with the
// database; for now it returns an empty, contract-typed list.
func (a *api) listCurrencies(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, []openapi.Currency{})
}
