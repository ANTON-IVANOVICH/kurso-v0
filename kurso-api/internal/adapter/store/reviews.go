package store

import (
	"context"
	"fmt"

	"github.com/ANTON-IVANOVICH/kurso-v0/kurso-api/internal/domain"
	"github.com/jackc/pgx/v5"
)

// ReviewsByExchanger returns reviews for an exchanger slug filtered by status,
// newest first.
func (s *Store) ReviewsByExchanger(ctx context.Context, slug string, status domain.ReviewStatus) ([]domain.Review, error) {
	rows, err := s.db.Query(ctx, `
		SELECT r.id::text, r.exchanger_id::text, COALESCE(r.author_name, 'Аноним'),
		       r.rating, r.title, r.body, r.status, r.created_at
		FROM reviews r
		JOIN exchangers e ON e.id = r.exchanger_id
		WHERE e.slug = $1 AND r.status = $2
		ORDER BY r.created_at DESC`, slug, string(status))
	if err != nil {
		return nil, fmt.Errorf("query reviews: %w", err)
	}
	defer rows.Close()

	var out []domain.Review
	for rows.Next() {
		var r domain.Review
		var st string
		if err := rows.Scan(&r.ID, &r.ExchangerID, &r.AuthorName, &r.Rating, &r.Title, &r.Body, &st, &r.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan review: %w", err)
		}
		r.Status = domain.ReviewStatus(st)
		r.ExchangerSlug = slug
		out = append(out, r)
	}
	return out, rows.Err()
}

// RatingSummary aggregates published reviews (average, count, star histogram).
func (s *Store) RatingSummary(ctx context.Context, slug string) (domain.RatingSummary, error) {
	rows, err := s.db.Query(ctx, `
		SELECT r.rating, count(*)
		FROM reviews r
		JOIN exchangers e ON e.id = r.exchanger_id
		WHERE e.slug = $1 AND r.status = 'published'
		GROUP BY r.rating`, slug)
	if err != nil {
		return domain.RatingSummary{}, fmt.Errorf("query rating summary: %w", err)
	}
	defer rows.Close()

	var sum domain.RatingSummary
	var weighted int
	for rows.Next() {
		var star, n int
		if err := rows.Scan(&star, &n); err != nil {
			return domain.RatingSummary{}, fmt.Errorf("scan rating: %w", err)
		}
		if star >= 1 && star <= 5 {
			sum.Histogram[star-1] = n
			sum.Count += n
			weighted += star * n
		}
	}
	if sum.Count > 0 {
		sum.Average = float64(weighted) / float64(sum.Count)
	}
	return sum, rows.Err()
}

// CreateReview inserts a review and returns it.
func (s *Store) CreateReview(ctx context.Context, in domain.NewReview) (domain.Review, error) {
	var r domain.Review
	var st string
	err := s.db.QueryRow(ctx, `
		INSERT INTO reviews (exchanger_id, author_name, author_email, rating, title, body, status, ip, published_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, CASE WHEN $7 = 'published' THEN now() END)
		RETURNING id::text, exchanger_id::text, COALESCE(author_name,'Аноним'), rating, title, body, status, created_at`,
		in.ExchangerID, in.AuthorName, in.AuthorEmail, in.Rating, in.Title, in.Body, string(in.Status), in.IP).
		Scan(&r.ID, &r.ExchangerID, &r.AuthorName, &r.Rating, &r.Title, &r.Body, &st, &r.CreatedAt)
	if err != nil {
		return domain.Review{}, fmt.Errorf("insert review: %w", err)
	}
	r.Status = domain.ReviewStatus(st)
	return r, nil
}

const adminReviewSelect = `
	SELECT r.id::text, r.exchanger_id::text, e.slug, e.name,
	       COALESCE(r.author_name,'Аноним'), r.rating, r.title, r.body, r.status, r.created_at
	FROM reviews r
	JOIN exchangers e ON e.id = r.exchanger_id`

// ReviewsByStatus returns reviews across all exchangers in a moderation state.
func (s *Store) ReviewsByStatus(ctx context.Context, status domain.ReviewStatus) ([]domain.Review, error) {
	rows, err := s.db.Query(ctx, adminReviewSelect+`
		WHERE r.status = $1
		ORDER BY r.created_at ASC`, string(status))
	if err != nil {
		return nil, fmt.Errorf("query reviews by status: %w", err)
	}
	defer rows.Close()

	var out []domain.Review
	for rows.Next() {
		r, err := scanAdminReview(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

func scanAdminReview(row pgx.Row) (domain.Review, error) {
	var r domain.Review
	var st string
	err := row.Scan(&r.ID, &r.ExchangerID, &r.ExchangerSlug, &r.ExchangerName,
		&r.AuthorName, &r.Rating, &r.Title, &r.Body, &st, &r.CreatedAt)
	r.Status = domain.ReviewStatus(st)
	return r, err
}

// SetReviewStatus updates a review's moderation state and returns the affected
// exchanger id (so the caller can recompute its rating).
func (s *Store) SetReviewStatus(ctx context.Context, id string, status domain.ReviewStatus, reason string) (string, error) {
	var exchangerID string
	err := s.db.QueryRow(ctx, `
		UPDATE reviews
		SET status = $2,
		    moderation_reason = NULLIF($3, ''),
		    published_at = CASE WHEN $2 = 'published' THEN now() ELSE published_at END
		WHERE id = $1
		RETURNING exchanger_id::text`, id, string(status), reason).Scan(&exchangerID)
	if err != nil {
		return "", fmt.Errorf("set review status: %w", err)
	}
	return exchangerID, nil
}

// RecomputeExchangerRating refreshes rating_avg + reviews_count from published
// reviews.
func (s *Store) RecomputeExchangerRating(ctx context.Context, exchangerID string) error {
	_, err := s.db.Exec(ctx, `
		UPDATE exchangers e
		SET rating_avg = sub.avg, reviews_count = sub.cnt
		FROM (
			SELECT round(avg(rating)::numeric, 2) AS avg, count(*) AS cnt
			FROM reviews WHERE exchanger_id = $1 AND status = 'published'
		) sub
		WHERE e.id = $1`, exchangerID)
	if err != nil {
		return fmt.Errorf("recompute rating: %w", err)
	}
	return nil
}

// CreateReport files a complaint against a review.
func (s *Store) CreateReport(ctx context.Context, reviewID, reason string, details *string) error {
	_, err := s.db.Exec(ctx, `
		INSERT INTO review_reports (review_id, reason, details)
		VALUES ($1, $2, $3)`, reviewID, reason, details)
	if err != nil {
		return fmt.Errorf("insert report: %w", err)
	}
	return nil
}

// OpenReports lists open review reports for the admin.
func (s *Store) OpenReports(ctx context.Context) ([]domain.ReviewReport, error) {
	rows, err := s.db.Query(ctx, `
		SELECT rp.id::text, rp.review_id::text, rp.reason, rp.details, rp.status, rp.created_at,
		       r.body, e.name
		FROM review_reports rp
		JOIN reviews r ON r.id = rp.review_id
		JOIN exchangers e ON e.id = r.exchanger_id
		WHERE rp.status = 'open'
		ORDER BY rp.created_at ASC`)
	if err != nil {
		return nil, fmt.Errorf("query reports: %w", err)
	}
	defer rows.Close()

	var out []domain.ReviewReport
	for rows.Next() {
		var rp domain.ReviewReport
		if err := rows.Scan(&rp.ID, &rp.ReviewID, &rp.Reason, &rp.Details, &rp.Status,
			&rp.CreatedAt, &rp.ReviewBody, &rp.ExchangerName); err != nil {
			return nil, fmt.Errorf("scan report: %w", err)
		}
		out = append(out, rp)
	}
	return out, rows.Err()
}

// SetReportStatus resolves a report (reviewed | dismissed).
func (s *Store) SetReportStatus(ctx context.Context, id, status string) error {
	_, err := s.db.Exec(ctx, `UPDATE review_reports SET status = $2 WHERE id = $1`, id, status)
	if err != nil {
		return fmt.Errorf("set report status: %w", err)
	}
	return nil
}
