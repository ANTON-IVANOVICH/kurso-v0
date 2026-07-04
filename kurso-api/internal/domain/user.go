package domain

// User is an end user of kurso-web (alerts, favorites, reviews).
type User struct {
	ID           string
	Email        string
	PasswordHash string // "" when the account has no password (Telegram/Google only)
	DisplayName  *string
	Status       string // active | blocked | deleted
}
