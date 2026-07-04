// Package domain holds the core business types, free of any transport or
// storage concern. Adapters map these onto HTTP/SQL; services operate on them.
package domain

import "time"

// CurrencyKind classifies a currency.
type CurrencyKind string

const (
	CurrencyCrypto CurrencyKind = "crypto"
	CurrencyFiat   CurrencyKind = "fiat"
	CurrencyCash   CurrencyKind = "cash"
)

// Currency is an entry in the currency catalogue.
type Currency struct {
	ID      string
	Code    string
	Name    string
	Kind    CurrencyKind
	Network *string
	IconURL *string
}

// Direction is an ordered from→to currency pair users compare rates for.
type Direction struct {
	ID        string
	Slug      string
	FromID    string
	ToID      string
	FromCode  string
	FromName  string
	ToCode    string
	ToName    string
	IsPopular bool
}

// ExchangerStatus is the moderation/availability state of an exchanger.
type ExchangerStatus string

const (
	ExchangerActive ExchangerStatus = "active"
	ExchangerPaused ExchangerStatus = "paused"
	ExchangerBanned ExchangerStatus = "banned"
)

// Exchanger is a third-party exchange service in the catalogue.
type Exchanger struct {
	ID           string
	Slug         string
	Name         string
	Status       ExchangerStatus
	WebsiteURL   *string
	ReferralTmpl *string
	LogoURL      *string
	Description  *string
	RatingAvg    *float64
	ReviewsCount int
	IsVerified   bool
	// Partner is true when the exchanger has a referral arrangement (drives the
	// "Партнёр" badge). Derived from ReferralTmpl presence.
	Partner bool

	// Aggregates over the exchanger's active rates, computed in the catalogue
	// query (see Store.Exchangers) — used by the exchangers catalogue page.
	ReserveTotal    *string  // Σ reserve across active rates, decimal string (nil if none)
	DirectionsCount int      // distinct directions with an active rate
	Assets          []string // distinct source-currency codes traded, e.g. USDT, BTC
	OnSinceYear     int      // year the exchanger was added (created_at)
}

// MapPoint is a located exchanger (cash desk) shown on the map, with its current
// rate for the selected direction.
type MapPoint struct {
	Slug         string
	Name         string
	Latitude     float64
	Longitude    float64
	Address      *string
	City         *string
	Hours        *string
	RatingAvg    *float64
	ReviewsCount int
	Partner      bool
	Rate         *string // best rate for the queried direction (nil if none)
}

// RateHistoryPoint is one time-bucket's best (highest) rate for a direction,
// used to draw the sparkline behind the price-alert builder.
type RateHistoryPoint struct {
	At   time.Time
	Rate string // decimal string
}

// RateRow is one exchanger's current rate for a direction — the denormalised
// shape the public rates API and SSE stream return.
type RateRow struct {
	ExchangerID   string
	ExchangerSlug string
	ExchangerName string
	Partner       bool
	Rating        *float64
	ReviewsCount  int
	Rate          string // decimal string — never a float, to keep precision
	Reserve       *string
	MinAmount     *string
	MaxAmount     *string
	FetchedAt     time.Time
}
