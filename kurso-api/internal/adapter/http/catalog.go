package httpapi

import (
	"net/http"

	"github.com/ANTON-IVANOVICH/kurso-v0/kurso-api/internal/adapter/http/openapi"
	"github.com/go-chi/chi/v5"
)

// listCurrencies returns the currency catalogue.
func (a *api) listCurrencies(w http.ResponseWriter, r *http.Request) {
	items, err := a.deps.Svc.Currencies(r.Context())
	if err != nil {
		a.serverError(w, "currencies", err)
		return
	}
	out := make([]openapi.Currency, 0, len(items))
	for _, c := range items {
		out = append(out, currencyDTO(c))
	}
	writeJSON(w, http.StatusOK, out)
}

// listDirections returns comparable currency pairs, popular first.
func (a *api) listDirections(w http.ResponseWriter, r *http.Request) {
	items, err := a.deps.Svc.Directions(r.Context())
	if err != nil {
		a.serverError(w, "directions", err)
		return
	}
	out := make([]openapi.Direction, 0, len(items))
	for _, d := range items {
		out = append(out, directionDTO(d))
	}
	writeJSON(w, http.StatusOK, out)
}

// listExchangers returns the active exchanger catalogue.
func (a *api) listExchangers(w http.ResponseWriter, r *http.Request) {
	items, err := a.deps.Svc.Exchangers(r.Context())
	if err != nil {
		a.serverError(w, "exchangers", err)
		return
	}
	out := make([]openapi.Exchanger, 0, len(items))
	for _, e := range items {
		out = append(out, exchangerDTO(e))
	}
	writeJSON(w, http.StatusOK, out)
}

// getExchanger returns a single exchanger card by slug.
func (a *api) getExchanger(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	e, ok, err := a.deps.Svc.ExchangerBySlug(r.Context(), slug)
	if err != nil {
		a.serverError(w, "exchanger", err)
		return
	}
	if !ok {
		writeError(w, http.StatusNotFound, "not_found", "обменник не найден")
		return
	}
	writeJSON(w, http.StatusOK, exchangerDTO(e))
}

// serverError logs and returns a generic 500 without leaking internals.
func (a *api) serverError(w http.ResponseWriter, op string, err error) {
	a.deps.Log.Error("request failed", "op", op, "err", err)
	writeError(w, http.StatusInternalServerError, "internal", "внутренняя ошибка сервиса")
}
