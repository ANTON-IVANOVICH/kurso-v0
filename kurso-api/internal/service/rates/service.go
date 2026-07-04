package rates

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ANTON-IVANOVICH/kurso-v0/kurso-api/internal/domain"
	"github.com/redis/go-redis/v9"
)

// Store is the persistence port the rates use cases depend on.
type Store interface {
	Currencies(ctx context.Context) ([]domain.Currency, error)
	Directions(ctx context.Context) ([]domain.Direction, error)
	DirectionBySlug(ctx context.Context, slug string) (domain.Direction, bool, error)
	Exchangers(ctx context.Context) ([]domain.Exchanger, error)
	ExchangerBySlug(ctx context.Context, slug string) (domain.Exchanger, bool, error)
	RatesByDirection(ctx context.Context, directionID string) ([]domain.RateRow, error)
	UpsertRate(ctx context.Context, exchangerID, directionID, rate, reserve string) error
}

// cacheTTL is short: the SSE stream carries live updates, so the cached read
// path only needs to absorb bursts of first-render traffic.
const cacheTTL = 8 * time.Second

// Service exposes the catalogue and rate reads used by the HTTP adapter.
type Service struct {
	store Store
	rdb   *redis.Client
	hub   *Hub
}

// NewService wires the store, Redis cache, and SSE hub together.
func NewService(store Store, rdb *redis.Client, hub *Hub) *Service {
	return &Service{store: store, rdb: rdb, hub: hub}
}

// Hub exposes the SSE fan-out for the streaming handler.
func (s *Service) Hub() *Hub { return s.hub }

func (s *Service) Currencies(ctx context.Context) ([]domain.Currency, error) {
	return s.store.Currencies(ctx)
}
func (s *Service) Directions(ctx context.Context) ([]domain.Direction, error) {
	return s.store.Directions(ctx)
}
func (s *Service) DirectionBySlug(ctx context.Context, slug string) (domain.Direction, bool, error) {
	return s.store.DirectionBySlug(ctx, slug)
}
func (s *Service) Exchangers(ctx context.Context) ([]domain.Exchanger, error) {
	return s.store.Exchangers(ctx)
}
func (s *Service) ExchangerBySlug(ctx context.Context, slug string) (domain.Exchanger, bool, error) {
	return s.store.ExchangerBySlug(ctx, slug)
}

// RatesByDirection returns the current rates for a direction, best-first, using
// a short Redis cache keyed by the direction slug.
func (s *Service) RatesByDirection(ctx context.Context, dir domain.Direction) ([]domain.RateRow, error) {
	key := cacheKey(dir.Slug)
	if s.rdb != nil {
		if b, err := s.rdb.Get(ctx, key).Bytes(); err == nil {
			var rows []domain.RateRow
			if json.Unmarshal(b, &rows) == nil {
				return rows, nil
			}
		}
	}
	rows, err := s.store.RatesByDirection(ctx, dir.ID)
	if err != nil {
		return nil, err
	}
	if s.rdb != nil {
		if b, err := json.Marshal(rows); err == nil {
			_ = s.rdb.Set(ctx, key, b, cacheTTL).Err()
		}
	}
	return rows, nil
}

// RepublishDirection re-reads a direction's rates, drops the cache, and pushes
// the fresh rows to the SSE hub. Called after an out-of-band write (e.g. a
// merchant manually refreshing a feed) so public readers see it immediately.
func (s *Service) RepublishDirection(ctx context.Context, dir domain.Direction) {
	fresh, err := s.store.RatesByDirection(ctx, dir.ID)
	if err != nil {
		return
	}
	s.invalidate(ctx, dir.Slug)
	s.hub.Publish(dir.Slug, fresh)
}

// invalidate drops the cached rates for a direction (called after a fresh tick).
func (s *Service) invalidate(ctx context.Context, slug string) {
	if s.rdb != nil {
		_ = s.rdb.Del(ctx, cacheKey(slug)).Err()
	}
}

func cacheKey(slug string) string { return "rates:" + slug }
