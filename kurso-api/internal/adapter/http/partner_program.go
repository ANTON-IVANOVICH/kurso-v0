package httpapi

import (
	"net/http"
	"strconv"

	"github.com/ANTON-IVANOVICH/kurso-v0/kurso-api/internal/service/auth"
)

// revsharePerClickRUB is the flat estimate used to project affiliate earnings
// from real clicks until the confirmed-exchange pipeline lands.
const revsharePerClickRUB = 10

type referralTagDTO struct {
	Tag    string `json:"tag"`
	Clicks int    `json:"clicks"`
}

type referralPointDTO struct {
	Day    string `json:"day"`
	Clicks int    `json:"clicks"`
}

type partnerOverviewDTO struct {
	Code          string             `json:"code"`
	RevsharePct   int                `json:"revsharePct"`
	CookieDays    int                `json:"cookieDays"`
	Clicks        int                `json:"clicks"`
	Registrations int                `json:"registrations"`
	EstimatedRUB  int                `json:"estimatedRub"`
	Series        []referralPointDTO `json:"series"`
	ByTag         []referralTagDTO   `json:"byTag"`
}

// userPartnerOverview returns the authenticated user's affiliate code and real
// performance (clicks from clickouts.ref_code, referred registrations).
func (a *api) userPartnerOverview(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value(userClaimsKey).(*auth.Claims)
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized", "Требуется вход")
		return
	}
	code, err := a.deps.Store.EnsureReferralCode(r.Context(), claims.Subject)
	if err != nil {
		a.serverError(w, "referral code", err)
		return
	}
	days := 30
	if v, perr := strconv.Atoi(r.URL.Query().Get("days")); perr == nil && v >= 1 && v <= 90 {
		days = v
	}
	stats, err := a.deps.Store.ReferralStats(r.Context(), claims.Subject, code, days)
	if err != nil {
		a.serverError(w, "referral stats", err)
		return
	}
	out := partnerOverviewDTO{
		Code:          code,
		RevsharePct:   30,
		CookieDays:    90,
		Clicks:        stats.Clicks,
		Registrations: stats.Registrations,
		EstimatedRUB:  stats.Clicks * revsharePerClickRUB,
		Series:        make([]referralPointDTO, 0, len(stats.Series)),
		ByTag:         make([]referralTagDTO, 0, len(stats.ByTag)),
	}
	for _, p := range stats.Series {
		out.Series = append(out.Series, referralPointDTO{Day: p.Day.Format("2006-01-02"), Clicks: p.Clicks})
	}
	for _, t := range stats.ByTag {
		out.ByTag = append(out.ByTag, referralTagDTO{Tag: t.Tag, Clicks: t.Clicks})
	}
	writeJSON(w, http.StatusOK, out)
}
