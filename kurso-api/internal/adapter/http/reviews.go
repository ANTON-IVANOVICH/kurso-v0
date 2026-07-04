package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/ANTON-IVANOVICH/kurso-v0/kurso-api/internal/domain"
	"github.com/go-chi/chi/v5"
)

// --- DTOs ---

type reviewDTO struct {
	ID        string    `json:"id"`
	Author    string    `json:"author"`
	Rating    int       `json:"rating"`
	Title     *string   `json:"title,omitempty"`
	Body      string    `json:"body"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

type adminReviewDTO struct {
	reviewDTO
	ExchangerSlug string `json:"exchangerSlug"`
	ExchangerName string `json:"exchangerName"`
}

type ratingSummaryDTO struct {
	Average   float64 `json:"average"`
	Count     int     `json:"count"`
	Histogram [5]int  `json:"histogram"`
}

type reviewsResponseDTO struct {
	Summary ratingSummaryDTO `json:"summary"`
	Reviews []reviewDTO      `json:"reviews"`
}

func toReviewDTO(r domain.Review) reviewDTO {
	return reviewDTO{
		ID: r.ID, Author: r.AuthorName, Rating: r.Rating, Title: r.Title,
		Body: r.Body, Status: string(r.Status), CreatedAt: r.CreatedAt,
	}
}

// --- auto-moderation ---

// stopWords route a review to manual moderation instead of auto-publishing.
// Accusatory/spam terms are held so a human confirms them (they are often the
// start of a dispute), while ordinary reviews publish immediately.
var stopWords = []string{"скам", "scam", "мошенн", "развод", "кидал", "лохотрон", "http://", "https://", "www."}

func autoModerate(body string) domain.ReviewStatus {
	trimmed := strings.TrimSpace(body)
	if len([]rune(trimmed)) < 10 {
		return domain.ReviewPending
	}
	lower := strings.ToLower(trimmed)
	for _, w := range stopWords {
		if strings.Contains(lower, w) {
			return domain.ReviewPending
		}
	}
	return domain.ReviewPublished
}

// --- public handlers ---

// listReviews returns published reviews for an exchanger plus the rating summary.
func (a *api) listReviews(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	reviews, err := a.deps.Store.ReviewsByExchanger(r.Context(), slug, domain.ReviewPublished)
	if err != nil {
		a.serverError(w, "reviews", err)
		return
	}
	summary, err := a.deps.Store.RatingSummary(r.Context(), slug)
	if err != nil {
		a.serverError(w, "rating summary", err)
		return
	}
	out := reviewsResponseDTO{
		Summary: ratingSummaryDTO{Average: summary.Average, Count: summary.Count, Histogram: summary.Histogram},
		Reviews: make([]reviewDTO, 0, len(reviews)),
	}
	for _, rv := range reviews {
		out.Reviews = append(out.Reviews, toReviewDTO(rv))
	}
	writeJSON(w, http.StatusOK, out)
}

type createReviewRequest struct {
	Author string `json:"author"`
	Email  string `json:"email"`
	Rating int    `json:"rating"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

// createReview accepts a new review, auto-moderates it, and returns it. A clean
// review is published immediately; flagged ones await moderation.
func (a *api) createReview(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	var req createReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "invalid JSON body")
		return
	}
	if req.Rating < 1 || req.Rating > 5 {
		writeError(w, http.StatusBadRequest, "invalid_rating", "Оценка должна быть от 1 до 5")
		return
	}
	if len([]rune(strings.TrimSpace(req.Body))) < 10 {
		writeError(w, http.StatusBadRequest, "invalid_body", "Отзыв слишком короткий")
		return
	}
	ex, ok, err := a.deps.Svc.ExchangerBySlug(r.Context(), slug)
	if err != nil {
		a.serverError(w, "exchanger lookup", err)
		return
	}
	if !ok {
		writeError(w, http.StatusNotFound, "not_found", "Обменник не найден")
		return
	}

	status := autoModerate(req.Body)
	author := strings.TrimSpace(req.Author)
	if author == "" {
		author = "Аноним"
	}
	in := domain.NewReview{
		ExchangerID: ex.ID,
		AuthorName:  author,
		Rating:      req.Rating,
		Body:        strings.TrimSpace(req.Body),
		Status:      status,
	}
	if t := strings.TrimSpace(req.Title); t != "" {
		in.Title = &t
	}
	if e := strings.TrimSpace(req.Email); e != "" {
		in.AuthorEmail = &e
	}

	created, err := a.deps.Store.CreateReview(r.Context(), in)
	if err != nil {
		a.serverError(w, "create review", err)
		return
	}
	if created.Status == domain.ReviewPublished {
		_ = a.deps.Store.RecomputeExchangerRating(r.Context(), ex.ID)
	}
	writeJSON(w, http.StatusCreated, toReviewDTO(created))
}

type reportReviewRequest struct {
	Reason  string `json:"reason"`
	Details string `json:"details"`
}

// reportReview files a complaint against a review.
func (a *api) reportReview(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req reportReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "invalid JSON body")
		return
	}
	reason := strings.TrimSpace(req.Reason)
	if reason == "" {
		reason = "other"
	}
	var details *string
	if d := strings.TrimSpace(req.Details); d != "" {
		details = &d
	}
	if err := a.deps.Store.CreateReport(r.Context(), id, reason, details); err != nil {
		a.serverError(w, "create report", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- admin handlers ---

// adminListReviews lists reviews in a moderation state (default: pending).
func (a *api) adminListReviews(w http.ResponseWriter, r *http.Request) {
	status := domain.ReviewStatus(r.URL.Query().Get("status"))
	if status == "" {
		status = domain.ReviewPending
	}
	reviews, err := a.deps.Store.ReviewsByStatus(r.Context(), status)
	if err != nil {
		a.serverError(w, "admin reviews", err)
		return
	}
	out := make([]adminReviewDTO, 0, len(reviews))
	for _, rv := range reviews {
		out = append(out, adminReviewDTO{
			reviewDTO:     toReviewDTO(rv),
			ExchangerSlug: rv.ExchangerSlug,
			ExchangerName: rv.ExchangerName,
		})
	}
	writeJSON(w, http.StatusOK, out)
}

type moderateReviewRequest struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
}

// adminModerateReview sets a review's status and recomputes the exchanger rating.
func (a *api) adminModerateReview(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req moderateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "invalid JSON body")
		return
	}
	switch domain.ReviewStatus(req.Status) {
	case domain.ReviewPublished, domain.ReviewRejected, domain.ReviewNeedsInfo, domain.ReviewPending:
	default:
		writeError(w, http.StatusBadRequest, "invalid_status", "Недопустимый статус")
		return
	}
	exchangerID, err := a.deps.Store.SetReviewStatus(r.Context(), id, domain.ReviewStatus(req.Status), req.Reason)
	if err != nil {
		a.serverError(w, "moderate review", err)
		return
	}
	// Publishing or unpublishing changes the average — refresh it either way.
	_ = a.deps.Store.RecomputeExchangerRating(r.Context(), exchangerID)
	w.WriteHeader(http.StatusNoContent)
}

type reportDTO struct {
	ID            string    `json:"id"`
	ReviewID      string    `json:"reviewId"`
	Reason        string    `json:"reason"`
	Details       *string   `json:"details,omitempty"`
	ReviewBody    string    `json:"reviewBody"`
	ExchangerName string    `json:"exchangerName"`
	CreatedAt     time.Time `json:"createdAt"`
}

// adminListReports lists open complaints about reviews.
func (a *api) adminListReports(w http.ResponseWriter, r *http.Request) {
	reports, err := a.deps.Store.OpenReports(r.Context())
	if err != nil {
		a.serverError(w, "admin reports", err)
		return
	}
	out := make([]reportDTO, 0, len(reports))
	for _, rp := range reports {
		out = append(out, reportDTO{
			ID: rp.ID, ReviewID: rp.ReviewID, Reason: rp.Reason, Details: rp.Details,
			ReviewBody: rp.ReviewBody, ExchangerName: rp.ExchangerName, CreatedAt: rp.CreatedAt,
		})
	}
	writeJSON(w, http.StatusOK, out)
}

type resolveReportRequest struct {
	Status string `json:"status"`
}

// adminResolveReport marks a report reviewed or dismissed.
func (a *api) adminResolveReport(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req resolveReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "invalid JSON body")
		return
	}
	if req.Status != "reviewed" && req.Status != "dismissed" {
		writeError(w, http.StatusBadRequest, "invalid_status", "Недопустимый статус")
		return
	}
	if err := a.deps.Store.SetReportStatus(r.Context(), id, req.Status); err != nil {
		a.serverError(w, "resolve report", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
