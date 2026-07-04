package domain

import "time"

// ExchangerUser is a representative of an exchanger who signs into the merchant
// cabinet (partner.kurso.io).
type ExchangerUser struct {
	ID           string
	ExchangerID  string
	Email        string
	PasswordHash string
	Role         string // owner | manager | viewer
	Status       string // invited | active | disabled
}

// MerchantIdentity resolves an exchanger_user to the exchanger it belongs to;
// stashed in the request context by the merchant auth middleware.
type MerchantIdentity struct {
	UserID        string
	Email         string
	Role          string
	ExchangerID   string
	ExchangerSlug string
	ExchangerName string
}

// MerchantMetrics are the headline numbers on the cabinet dashboard.
type MerchantMetrics struct {
	RatesActive       int     // active rates with a fresh feed
	RatesTotal        int     // active rates total
	RatesStale        int     // active rates whose feed is delayed
	ClicksToday       int     // clickouts since local midnight (UTC)
	ClicksYesterday   int     // clickouts the previous day
	RatingAvg         float64 // published-review average
	ReviewsCount      int     // published reviews
	ReviewsUnanswered int     // published reviews without a reply
}

// MerchantRate is one direction's current rate row in the cabinet courses table.
type MerchantRate struct {
	DirectionID   string
	DirectionSlug string
	FromCode      string
	ToCode        string
	Rate          *string
	Reserve       *string
	FetchedAt     *time.Time
	IsActive      bool
}

// MerchantReview is a review shown in the cabinet with its optional reply.
type MerchantReview struct {
	ID        string
	Author    string
	Rating    int
	Title     *string
	Body      string
	CreatedAt time.Time
	Reply     *string
	ReplyAt   *time.Time
}

// TrafficPoint is one day's clickout total for the traffic sparkline.
type TrafficPoint struct {
	Day    time.Time
	Clicks int
}

// TrafficDirection is the per-direction clickout breakdown.
type TrafficDirection struct {
	DirectionSlug string
	FromCode      string
	ToCode        string
	Clicks        int
}

// MerchantComplaint is a user complaint (a report filed against one of the
// exchanger's reviews) surfaced in the cabinet.
type MerchantComplaint struct {
	ID         string
	ReviewID   string
	Reason     string
	Details    *string
	Status     string
	CreatedAt  time.Time
	ReviewBody string
	Author     string
	Rating     int
}

// Payout is a billing settlement row (partner_payouts).
type Payout struct {
	ID          string
	PeriodStart time.Time
	PeriodEnd   time.Time
	ClicksCount int
	Amount      string
	Currency    string
	Status      string
	CreatedAt   time.Time
	PaidAt      *time.Time
}
