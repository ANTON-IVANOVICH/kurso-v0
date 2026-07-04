package rates

import (
	"context"
	"log/slog"
	"math/rand/v2"
	"strconv"
	"time"
)

// Ticker keeps rates fresh by applying a small random walk each interval,
// persisting the result, invalidating the cache, and publishing to the SSE hub.
//
// It stands in for the real parser runner (BestChange-XML + REST adapters): the
// public read/stream contract is identical, so swapping in real parsers later is
// a matter of replacing the rate source — not the delivery path.
type Ticker struct {
	store    Store
	svc      *Service
	log      *slog.Logger
	interval time.Duration
}

// NewTicker builds a runner over the given store/service.
func NewTicker(store Store, svc *Service, log *slog.Logger, interval time.Duration) *Ticker {
	return &Ticker{store: store, svc: svc, log: log, interval: interval}
}

// Run ticks until the context is cancelled. Safe to call in its own goroutine.
func (t *Ticker) Run(ctx context.Context) {
	tk := time.NewTicker(t.interval)
	defer tk.Stop()
	t.tick(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case <-tk.C:
			t.tick(ctx)
		}
	}
}

func (t *Ticker) tick(ctx context.Context) {
	dirs, err := t.store.Directions(ctx)
	if err != nil {
		t.log.Warn("rate ticker: list directions failed", "err", err)
		return
	}
	for _, d := range dirs {
		rows, err := t.store.RatesByDirection(ctx, d.ID)
		if err != nil {
			t.log.Warn("rate ticker: read rates failed", "direction", d.Slug, "err", err)
			continue
		}
		for _, r := range rows {
			old, perr := strconv.ParseFloat(r.Rate, 64)
			if perr != nil {
				continue
			}
			next := old * (1 + (rand.Float64()-0.5)*0.004) // ±0.2% random walk
			reserve := ""
			if r.Reserve != nil {
				reserve = *r.Reserve
			}
			if err := t.store.UpsertRate(ctx, r.ExchangerID, d.ID, formatRate(next), reserve); err != nil {
				t.log.Warn("rate ticker: upsert failed", "exchanger", r.ExchangerSlug, "err", err)
			}
		}
		fresh, err := t.store.RatesByDirection(ctx, d.ID)
		if err != nil {
			continue
		}
		t.svc.invalidate(ctx, d.Slug)
		t.svc.Hub().Publish(d.Slug, fresh)
	}
}

func formatRate(v float64) string { return strconv.FormatFloat(v, 'f', 8, 64) }
