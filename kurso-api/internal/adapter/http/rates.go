package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ANTON-IVANOVICH/kurso-v0/kurso-api/internal/adapter/http/openapi"
	"github.com/go-chi/chi/v5"
)

// getRates returns the current rates for a direction (best first), briefly cached.
func (a *api) getRates(w http.ResponseWriter, r *http.Request) {
	dir, ok, err := a.deps.Svc.DirectionBySlug(r.Context(), chi.URLParam(r, "direction"))
	if err != nil {
		a.serverError(w, "rates.direction", err)
		return
	}
	if !ok {
		writeError(w, http.StatusNotFound, "not_found", "направление не найдено")
		return
	}
	rows, err := a.deps.Svc.RatesByDirection(r.Context(), dir)
	if err != nil {
		a.serverError(w, "rates", err)
		return
	}
	writeJSON(w, http.StatusOK, openapi.RatesResponse{
		Direction: directionDTO(dir),
		Rates:     rateRowDTOs(rows),
	})
}

// streamRates is the SSE endpoint. It emits an initial snapshot, then every
// update the runner publishes for this direction, until the client disconnects.
func (a *api) streamRates(w http.ResponseWriter, r *http.Request) {
	dir, ok, err := a.deps.Svc.DirectionBySlug(r.Context(), chi.URLParam(r, "direction"))
	if err != nil {
		a.serverError(w, "stream.direction", err)
		return
	}
	if !ok {
		writeError(w, http.StatusNotFound, "not_found", "направление не найдено")
		return
	}
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeError(w, http.StatusInternalServerError, "internal", "streaming unsupported")
		return
	}
	// Clear the server WriteTimeout for this long-lived connection.
	_ = http.NewResponseController(w).SetWriteDeadline(time.Time{})

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	send := func(rows []openapi.RateRow) bool {
		b, err := json.Marshal(rows)
		if err != nil {
			return true
		}
		if _, err := fmt.Fprintf(w, "event: rates\ndata: %s\n\n", b); err != nil {
			return false
		}
		flusher.Flush()
		return true
	}

	// initial snapshot so a fresh subscriber renders immediately
	if rows, err := a.deps.Svc.RatesByDirection(r.Context(), dir); err == nil {
		if !send(rateRowDTOs(rows)) {
			return
		}
	}

	ch := a.deps.Svc.Hub().Subscribe(dir.Slug)
	defer a.deps.Svc.Hub().Unsubscribe(dir.Slug, ch)

	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			return
		case rows, open := <-ch:
			if !open {
				return
			}
			if !send(rateRowDTOs(rows)) {
				return
			}
		}
	}
}
