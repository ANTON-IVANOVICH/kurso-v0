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

type mapPointDTO struct {
	Slug         string   `json:"slug"`
	Name         string   `json:"name"`
	Lat          float64  `json:"lat"`
	Lng          float64  `json:"lng"`
	Address      *string  `json:"address,omitempty"`
	City         *string  `json:"city,omitempty"`
	Hours        *string  `json:"hours,omitempty"`
	RatingAvg    *float64 `json:"ratingAvg,omitempty"`
	ReviewsCount int      `json:"reviewsCount"`
	Partner      bool     `json:"partner"`
	Rate         *string  `json:"rate,omitempty"`
}

type mapResponseDTO struct {
	Direction openapi.Direction `json:"direction"`
	Points    []mapPointDTO     `json:"points"`
}

// listMapPoints returns located exchangers (cash desks) with their rate for the
// selected direction (query `direction`, default usdt-tinkoff). Powers the map.
func (a *api) listMapPoints(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get("direction")
	if slug == "" {
		slug = "usdt-tinkoff"
	}
	dir, ok, err := a.deps.Svc.DirectionBySlug(r.Context(), slug)
	if err != nil {
		a.serverError(w, "map.direction", err)
		return
	}
	if !ok {
		writeError(w, http.StatusNotFound, "not_found", "направление не найдено")
		return
	}
	points, err := a.deps.Store.MapExchangers(r.Context(), dir.ID)
	if err != nil {
		a.serverError(w, "map points", err)
		return
	}
	out := mapResponseDTO{Direction: directionDTO(dir), Points: make([]mapPointDTO, 0, len(points))}
	for _, p := range points {
		out.Points = append(out.Points, mapPointDTO{
			Slug: p.Slug, Name: p.Name, Lat: p.Latitude, Lng: p.Longitude,
			Address: p.Address, City: p.City, Hours: p.Hours, RatingAvg: p.RatingAvg,
			ReviewsCount: p.ReviewsCount, Partner: p.Partner, Rate: p.Rate,
		})
	}
	writeJSON(w, http.StatusOK, out)
}

// serverError logs and returns a generic 500 without leaking internals.
func (a *api) serverError(w http.ResponseWriter, op string, err error) {
	a.deps.Log.Error("request failed", "op", op, "err", err)
	writeError(w, http.StatusInternalServerError, "internal", "внутренняя ошибка сервиса")
}
