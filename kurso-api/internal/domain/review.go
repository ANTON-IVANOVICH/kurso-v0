package domain

import "time"

// ReviewStatus is the moderation state of a review.
type ReviewStatus string

const (
	ReviewPending   ReviewStatus = "pending"
	ReviewPublished ReviewStatus = "published"
	ReviewRejected  ReviewStatus = "rejected"
	ReviewNeedsInfo ReviewStatus = "needs_info"
)

// Review is a user's rating + text about an exchanger.
type Review struct {
	ID            string
	ExchangerID   string
	ExchangerSlug string // joined, for admin lists
	ExchangerName string // joined, for admin lists
	AuthorName    string
	Rating        int
	Title         *string
	Body          string
	Status        ReviewStatus
	CreatedAt     time.Time
}

// NewReview is the input for creating a review.
type NewReview struct {
	ExchangerID string
	AuthorName  string
	AuthorEmail *string
	Rating      int
	Title       *string
	Body        string
	Status      ReviewStatus
	IP          *string
}

// RatingSummary aggregates published reviews for an exchanger.
type RatingSummary struct {
	Average   float64
	Count     int
	Histogram [5]int // Histogram[0] = 1★ … Histogram[4] = 5★
}

// ReviewReport is a user complaint about a published review.
type ReviewReport struct {
	ID            string
	ReviewID      string
	Reason        string
	Details       *string
	Status        string // open | reviewed | dismissed
	CreatedAt     time.Time
	ReviewBody    string // joined
	ExchangerName string // joined
}
