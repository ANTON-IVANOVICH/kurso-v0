package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ANTON-IVANOVICH/kurso-v0/kurso-api/internal/domain"
	"github.com/go-chi/chi/v5"
)

// perClickRUB is the CPA rate used to estimate the running (unbilled) balance
// from real clickouts until a payout is finalised by the platform.
const perClickRUB = 25

// --- dashboard ---

type merchantMetricsDTO struct {
	RatesActive       int     `json:"ratesActive"`
	RatesTotal        int     `json:"ratesTotal"`
	RatesStale        int     `json:"ratesStale"`
	ClicksToday       int     `json:"clicksToday"`
	ClicksYesterday   int     `json:"clicksYesterday"`
	RatingAvg         float64 `json:"ratingAvg"`
	ReviewsCount      int     `json:"reviewsCount"`
	ReviewsUnanswered int     `json:"reviewsUnanswered"`
}

type dashboardDTO struct {
	Metrics merchantMetricsDTO `json:"metrics"`
	Traffic []trafficPointDTO  `json:"traffic"`
}

type trafficPointDTO struct {
	Day    string `json:"day"`
	Clicks int    `json:"clicks"`
}

func (a *api) merchantDashboard(w http.ResponseWriter, r *http.Request) {
	m := merchantFromContext(r.Context())
	metrics, err := a.deps.Store.MerchantMetrics(r.Context(), m.ExchangerID)
	if err != nil {
		a.serverError(w, "merchant metrics", err)
		return
	}
	series, err := a.deps.Store.TrafficSeries(r.Context(), m.ExchangerID, 7)
	if err != nil {
		a.serverError(w, "merchant traffic series", err)
		return
	}
	writeJSON(w, http.StatusOK, dashboardDTO{
		Metrics: merchantMetricsDTO(metrics),
		Traffic: trafficPointDTOs(series),
	})
}

func trafficPointDTOs(series []domain.TrafficPoint) []trafficPointDTO {
	out := make([]trafficPointDTO, 0, len(series))
	for _, p := range series {
		out = append(out, trafficPointDTO{Day: p.Day.Format("2006-01-02"), Clicks: p.Clicks})
	}
	return out
}

// --- rates / courses ---

type merchantRateDTO struct {
	DirectionID   string  `json:"directionId"`
	DirectionSlug string  `json:"directionSlug"`
	FromCode      string  `json:"fromCode"`
	ToCode        string  `json:"toCode"`
	Rate          *string `json:"rate"`
	Reserve       *string `json:"reserve"`
	FetchedAt     *string `json:"fetchedAt"`
	Feed          string  `json:"feed"` // ok | delayed | down
}

// feedStatus classifies a rate row by freshness for the cabinet UI.
func feedStatus(mr domain.MerchantRate) string {
	if !mr.IsActive || mr.Rate == nil || mr.FetchedAt == nil {
		return "down"
	}
	if time.Since(*mr.FetchedAt) > 3*time.Minute {
		return "delayed"
	}
	return "ok"
}

func merchantRateDTOOf(mr domain.MerchantRate) merchantRateDTO {
	var fetched *string
	if mr.FetchedAt != nil {
		s := mr.FetchedAt.Format(time.RFC3339)
		fetched = &s
	}
	return merchantRateDTO{
		DirectionID: mr.DirectionID, DirectionSlug: mr.DirectionSlug,
		FromCode: mr.FromCode, ToCode: mr.ToCode,
		Rate: mr.Rate, Reserve: mr.Reserve, FetchedAt: fetched, Feed: feedStatus(mr),
	}
}

func (a *api) merchantRates(w http.ResponseWriter, r *http.Request) {
	m := merchantFromContext(r.Context())
	rates, err := a.deps.Store.MerchantRates(r.Context(), m.ExchangerID)
	if err != nil {
		a.serverError(w, "merchant rates", err)
		return
	}
	out := make([]merchantRateDTO, 0, len(rates))
	for _, mr := range rates {
		out = append(out, merchantRateDTOOf(mr))
	}
	writeJSON(w, http.StatusOK, out)
}

// merchantRefreshRate re-polls one direction's feed (re-stamps fetched_at) and
// republishes the direction so public readers see the fresh timestamp.
func (a *api) merchantRefreshRate(w http.ResponseWriter, r *http.Request) {
	m := merchantFromContext(r.Context())
	directionID := chi.URLParam(r, "directionId")

	dir, ok, err := a.deps.Store.DirectionByID(r.Context(), directionID)
	if err != nil {
		a.serverError(w, "direction lookup", err)
		return
	}
	if !ok {
		writeError(w, http.StatusNotFound, "not_found", "Направление не найдено")
		return
	}
	touched, err := a.deps.Store.TouchRate(r.Context(), m.ExchangerID, directionID)
	if err != nil {
		a.serverError(w, "touch rate", err)
		return
	}
	if !touched {
		writeError(w, http.StatusNotFound, "not_found", "Курс по этому направлению не настроен")
		return
	}
	a.deps.Svc.RepublishDirection(r.Context(), dir)

	// Return the refreshed row.
	rates, err := a.deps.Store.MerchantRates(r.Context(), m.ExchangerID)
	if err != nil {
		a.serverError(w, "merchant rates", err)
		return
	}
	for _, mr := range rates {
		if mr.DirectionID == directionID {
			writeJSON(w, http.StatusOK, merchantRateDTOOf(mr))
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- reviews + replies ---

type merchantReviewDTO struct {
	ID        string    `json:"id"`
	Author    string    `json:"author"`
	Rating    int       `json:"rating"`
	Title     *string   `json:"title,omitempty"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
	Reply     *string   `json:"reply,omitempty"`
	ReplyAt   *string   `json:"replyAt,omitempty"`
}

func (a *api) merchantReviews(w http.ResponseWriter, r *http.Request) {
	m := merchantFromContext(r.Context())
	reviews, err := a.deps.Store.MerchantReviews(r.Context(), m.ExchangerID)
	if err != nil {
		a.serverError(w, "merchant reviews", err)
		return
	}
	out := make([]merchantReviewDTO, 0, len(reviews))
	for _, rv := range reviews {
		dto := merchantReviewDTO{
			ID: rv.ID, Author: rv.Author, Rating: rv.Rating, Title: rv.Title,
			Body: rv.Body, CreatedAt: rv.CreatedAt, Reply: rv.Reply,
		}
		if rv.ReplyAt != nil {
			s := rv.ReplyAt.Format(time.RFC3339)
			dto.ReplyAt = &s
		}
		out = append(out, dto)
	}
	writeJSON(w, http.StatusOK, out)
}

type replyRequest struct {
	Body string `json:"body"`
}

// merchantReplyReview creates/updates the exchanger's reply to one of its reviews.
func (a *api) merchantReplyReview(w http.ResponseWriter, r *http.Request) {
	m := merchantFromContext(r.Context())
	reviewID := chi.URLParam(r, "id")
	var req replyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "invalid JSON body")
		return
	}
	body := strings.TrimSpace(req.Body)
	if len([]rune(body)) < 2 {
		writeError(w, http.StatusBadRequest, "invalid_body", "Ответ слишком короткий")
		return
	}
	if len([]rune(body)) > 600 {
		writeError(w, http.StatusBadRequest, "invalid_body", "Ответ длиннее 600 символов")
		return
	}
	owns, err := a.deps.Store.ReviewBelongsToExchanger(r.Context(), reviewID, m.ExchangerID)
	if err != nil {
		a.serverError(w, "review ownership", err)
		return
	}
	if !owns {
		writeError(w, http.StatusForbidden, "forbidden", "Отзыв относится к другому обменнику")
		return
	}
	if err := a.deps.Store.UpsertReviewReply(r.Context(), reviewID, m.UserID, body); err != nil {
		a.serverError(w, "upsert reply", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- traffic ---

type trafficDirectionDTO struct {
	DirectionSlug string `json:"directionSlug"`
	FromCode      string `json:"fromCode"`
	ToCode        string `json:"toCode"`
	Clicks        int    `json:"clicks"`
}

type trafficDTO struct {
	Days       int                   `json:"days"`
	Total      int                   `json:"total"`
	Series     []trafficPointDTO     `json:"series"`
	Directions []trafficDirectionDTO `json:"directions"`
}

// clampDays keeps the traffic window to the three tabs the UI offers.
func clampDays(raw string) int {
	switch raw {
	case "30":
		return 30
	case "7":
		return 7
	default:
		return 14
	}
}

func (a *api) merchantTraffic(w http.ResponseWriter, r *http.Request) {
	m := merchantFromContext(r.Context())
	days := clampDays(r.URL.Query().Get("days"))

	series, err := a.deps.Store.TrafficSeries(r.Context(), m.ExchangerID, days)
	if err != nil {
		a.serverError(w, "traffic series", err)
		return
	}
	byDir, err := a.deps.Store.TrafficByDirection(r.Context(), m.ExchangerID, days)
	if err != nil {
		a.serverError(w, "traffic by direction", err)
		return
	}
	total := 0
	for _, p := range series {
		total += p.Clicks
	}
	dirs := make([]trafficDirectionDTO, 0, len(byDir))
	for _, d := range byDir {
		dirs = append(dirs, trafficDirectionDTO{
			DirectionSlug: d.DirectionSlug, FromCode: d.FromCode, ToCode: d.ToCode, Clicks: d.Clicks,
		})
	}
	writeJSON(w, http.StatusOK, trafficDTO{
		Days: days, Total: total, Series: trafficPointDTOs(series), Directions: dirs,
	})
}

// --- profile ---

type merchantProfileDTO struct {
	Slug         string   `json:"slug"`
	Name         string   `json:"name"`
	Description  *string  `json:"description"`
	WebsiteURL   *string  `json:"websiteUrl"`
	LogoURL      *string  `json:"logoUrl"`
	Status       string   `json:"status"`
	IsVerified   bool     `json:"isVerified"`
	RatingAvg    *string  `json:"ratingAvg"`
	ReviewsCount int      `json:"reviewsCount"`
	OnSince      int      `json:"onSince"`
	Assets       []string `json:"assets"`
	ReserveTotal *string  `json:"reserveTotal"`
}

func (a *api) merchantProfile(w http.ResponseWriter, r *http.Request) {
	m := merchantFromContext(r.Context())
	e, ok, err := a.deps.Svc.ExchangerBySlug(r.Context(), m.ExchangerSlug)
	if err != nil {
		a.serverError(w, "merchant profile", err)
		return
	}
	if !ok {
		writeError(w, http.StatusNotFound, "not_found", "Обменник не найден")
		return
	}
	writeJSON(w, http.StatusOK, profileDTO(e))
}

func profileDTO(e domain.Exchanger) merchantProfileDTO {
	var rating *string
	if e.RatingAvg != nil {
		s := strconv.FormatFloat(*e.RatingAvg, 'f', 2, 64)
		rating = &s
	}
	assets := e.Assets
	if assets == nil {
		assets = []string{}
	}
	return merchantProfileDTO{
		Slug: e.Slug, Name: e.Name, Description: e.Description, WebsiteURL: e.WebsiteURL,
		LogoURL: e.LogoURL, Status: string(e.Status), IsVerified: e.IsVerified,
		RatingAvg: rating, ReviewsCount: e.ReviewsCount, OnSince: e.OnSinceYear,
		Assets: assets, ReserveTotal: e.ReserveTotal,
	}
}

type updateProfileRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	WebsiteURL  string `json:"websiteUrl"`
}

func (a *api) merchantUpdateProfile(w http.ResponseWriter, r *http.Request) {
	m := merchantFromContext(r.Context())
	if m.Role == "viewer" {
		writeError(w, http.StatusForbidden, "forbidden", "Недостаточно прав для изменения профиля")
		return
	}
	var req updateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "invalid JSON body")
		return
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		writeError(w, http.StatusBadRequest, "invalid_name", "Название не может быть пустым")
		return
	}
	if err := a.deps.Store.UpdateExchangerProfile(r.Context(), m.ExchangerID, name,
		strings.TrimSpace(req.Description), strings.TrimSpace(req.WebsiteURL)); err != nil {
		a.serverError(w, "update profile", err)
		return
	}
	e, _, err := a.deps.Svc.ExchangerBySlug(r.Context(), m.ExchangerSlug)
	if err != nil {
		a.serverError(w, "reload profile", err)
		return
	}
	writeJSON(w, http.StatusOK, profileDTO(e))
}

// --- complaints ---

type merchantComplaintDTO struct {
	ID         string    `json:"id"`
	ReviewID   string    `json:"reviewId"`
	Reason     string    `json:"reason"`
	Details    *string   `json:"details,omitempty"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"createdAt"`
	ReviewBody string    `json:"reviewBody"`
	Author     string    `json:"author"`
	Rating     int       `json:"rating"`
}

func (a *api) merchantComplaints(w http.ResponseWriter, r *http.Request) {
	m := merchantFromContext(r.Context())
	items, err := a.deps.Store.MerchantComplaints(r.Context(), m.ExchangerID)
	if err != nil {
		a.serverError(w, "merchant complaints", err)
		return
	}
	out := make([]merchantComplaintDTO, 0, len(items))
	for _, c := range items {
		out = append(out, merchantComplaintDTO{
			ID: c.ID, ReviewID: c.ReviewID, Reason: c.Reason, Details: c.Details,
			Status: c.Status, CreatedAt: c.CreatedAt, ReviewBody: c.ReviewBody,
			Author: c.Author, Rating: c.Rating,
		})
	}
	writeJSON(w, http.StatusOK, out)
}

// --- billing ---

type payoutDTO struct {
	ID          string  `json:"id"`
	PeriodStart string  `json:"periodStart"`
	PeriodEnd   string  `json:"periodEnd"`
	ClicksCount int     `json:"clicksCount"`
	Amount      string  `json:"amount"`
	Currency    string  `json:"currency"`
	Status      string  `json:"status"`
	PaidAt      *string `json:"paidAt,omitempty"`
}

type billingDTO struct {
	PerClick     int `json:"perClick"`
	CurrentMonth struct {
		Clicks    int `json:"clicks"`
		Estimated int `json:"estimated"`
	} `json:"currentMonth"`
	Payouts []payoutDTO `json:"payouts"`
}

func (a *api) merchantBilling(w http.ResponseWriter, r *http.Request) {
	m := merchantFromContext(r.Context())
	clicks, err := a.deps.Store.ClicksThisMonth(r.Context(), m.ExchangerID)
	if err != nil {
		a.serverError(w, "clicks this month", err)
		return
	}
	payouts, err := a.deps.Store.MerchantPayouts(r.Context(), m.ExchangerID)
	if err != nil {
		a.serverError(w, "merchant payouts", err)
		return
	}
	var out billingDTO
	out.PerClick = perClickRUB
	out.CurrentMonth.Clicks = clicks
	out.CurrentMonth.Estimated = clicks * perClickRUB
	out.Payouts = make([]payoutDTO, 0, len(payouts))
	for _, p := range payouts {
		dto := payoutDTO{
			ID: p.ID, PeriodStart: p.PeriodStart.Format("2006-01-02"),
			PeriodEnd: p.PeriodEnd.Format("2006-01-02"), ClicksCount: p.ClicksCount,
			Amount: p.Amount, Currency: p.Currency, Status: p.Status,
		}
		if p.PaidAt != nil {
			s := p.PaidAt.Format(time.RFC3339)
			dto.PaidAt = &s
		}
		out.Payouts = append(out.Payouts, dto)
	}
	writeJSON(w, http.StatusOK, out)
}
