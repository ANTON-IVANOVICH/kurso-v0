package domain

// User is an end user of kurso-web (alerts, favorites, reviews).
type User struct {
	ID           string
	Email        string
	PasswordHash string // "" when the account has no password (Telegram/Google only)
	DisplayName  *string
	Status       string // active | blocked | deleted
}

// ReferralTag is per-UTM-tag clickout attribution for the affiliate cabinet.
type ReferralTag struct {
	Tag    string // suffix after the base code, or "" for the untagged link
	Clicks int
}

// ReferralStats aggregates a partner's affiliate performance from real clickouts
// and referred registrations.
type ReferralStats struct {
	Clicks        int
	Registrations int
	Series        []TrafficPoint
	ByTag         []ReferralTag
}
