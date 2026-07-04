package store

import (
	"context"
	"fmt"
	"strconv"
)

type seedCurrency struct {
	code, name, kind, network string
	sort                      int
}

type seedExchanger struct {
	slug, name string
	rating     float64
	reviews    int
	verified   bool
	partner    bool
}

type seedDirection struct {
	from, to, slug string
	popular        bool
	base           float64 // reference rate (to-per-from) the ticker jitters around
}

var (
	seedCurrencies = []seedCurrency{
		{"USDT", "Tether", "crypto", "TRC20", 1},
		{"BTC", "Bitcoin", "crypto", "", 2},
		{"ETH", "Ethereum", "crypto", "", 3},
		{"TON", "Toncoin", "crypto", "", 4},
		{"TINKOFF", "Тинькофф", "fiat", "", 10},
		{"SBER", "Сбербанк", "fiat", "", 11},
		{"RUBCASH", "Наличные RUB", "cash", "", 12},
	}

	// order matters — the "best" exchanger is applied via per-exchanger factor below.
	seedExchangers = []seedExchanger{
		{"cryptobridge", "CryptoBridge", 4.9, 1203, true, true},
		{"netex24", "NetEx24", 4.9, 2104, true, true},
		{"24paybank", "24Paybank", 4.7, 560, true, false},
		{"coino", "Coino", 4.8, 420, false, false},
		{"baksman", "BaksMan", 4.4, 310, false, false},
		{"bitx", "BitX", 4.5, 180, false, false},
	}

	seedDirections = []seedDirection{
		{"USDT", "TINKOFF", "usdt-tinkoff", true, 81.20},
		{"USDT", "SBER", "usdt-sber", true, 81.05},
		{"USDT", "RUBCASH", "usdt-cash", true, 80.60},
		{"BTC", "TINKOFF", "btc-tinkoff", true, 8110000},
		{"ETH", "TINKOFF", "eth-tinkoff", false, 312400},
		{"TON", "TINKOFF", "ton-tinkoff", false, 521.40},
	}

	// Per-exchanger multiplier applied to a direction's base rate. Index-aligned
	// with seedExchangers; CryptoBridge/Coino lead, BaksMan trails.
	exchangerFactor = []float64{1.0000, 0.9969, 0.9943, 1.0005, 0.9741, 0.9915}
)

// Seed idempotently loads the Stage-1 catalogue and an initial set of rates.
// It is a no-op once exchangers exist, so restarts keep live (ticker) rates.
func (s *Store) Seed(ctx context.Context) error {
	var existing int
	if err := s.db.QueryRow(ctx, `SELECT count(*) FROM exchangers`).Scan(&existing); err != nil {
		return fmt.Errorf("seed: count exchangers: %w", err)
	}
	if existing > 0 {
		return nil
	}

	// currencies
	for _, c := range seedCurrencies {
		if _, err := s.db.Exec(ctx, `
			INSERT INTO currencies (code, name, kind, network, sort_order)
			VALUES ($1, $2, $3, NULLIF($4,''), $5)
			ON CONFLICT DO NOTHING`, c.code, c.name, c.kind, c.network, c.sort); err != nil {
			return fmt.Errorf("seed currency %s: %w", c.code, err)
		}
	}
	currencyID, err := s.idMap(ctx, `SELECT code, id::text FROM currencies`)
	if err != nil {
		return err
	}

	// exchangers
	for _, e := range seedExchangers {
		var ref *string
		if e.partner {
			tmpl := fmt.Sprintf("https://%s.example/exchange?ref=kurso&d={direction}", e.slug)
			ref = &tmpl
		}
		site := fmt.Sprintf("https://%s.example", e.slug)
		if _, err := s.db.Exec(ctx, `
			INSERT INTO exchangers (slug, name, website_url, referral_url_template, rating_avg, reviews_count, is_verified)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT (slug) DO NOTHING`,
			e.slug, e.name, site, ref, e.rating, e.reviews, e.verified); err != nil {
			return fmt.Errorf("seed exchanger %s: %w", e.slug, err)
		}
	}
	exchangerID, err := s.idMap(ctx, `SELECT slug, id::text FROM exchangers`)
	if err != nil {
		return err
	}

	// directions
	for i, d := range seedDirections {
		from, to := currencyID[d.from], currencyID[d.to]
		if from == "" || to == "" {
			continue
		}
		if _, err := s.db.Exec(ctx, `
			INSERT INTO directions (from_currency_id, to_currency_id, slug, is_popular, sort_order)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (slug) DO NOTHING`, from, to, d.slug, d.popular, i); err != nil {
			return fmt.Errorf("seed direction %s: %w", d.slug, err)
		}
	}
	directionID, err := s.idMap(ctx, `SELECT slug, id::text FROM directions`)
	if err != nil {
		return err
	}

	// initial rates
	for _, d := range seedDirections {
		dirID := directionID[d.slug]
		if dirID == "" {
			continue
		}
		for i, e := range seedExchangers {
			exID := exchangerID[e.slug]
			if exID == "" {
				continue
			}
			rate := d.base * exchangerFactor[i]
			reserve := 5_000_000 + float64(i)*3_200_000
			if err := s.UpsertRate(ctx, exID, dirID, ftoa(rate), ftoa(reserve)); err != nil {
				return fmt.Errorf("seed rate %s/%s: %w", e.slug, d.slug, err)
			}
		}
	}
	return nil
}

// idMap builds a key→id lookup from a two-column query.
func (s *Store) idMap(ctx context.Context, query string) (map[string]string, error) {
	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("seed idMap: %w", err)
	}
	defer rows.Close()
	m := map[string]string{}
	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err != nil {
			return nil, fmt.Errorf("seed idMap scan: %w", err)
		}
		m[k] = v
	}
	return m, rows.Err()
}

func ftoa(v float64) string { return strconv.FormatFloat(v, 'f', 8, 64) }
