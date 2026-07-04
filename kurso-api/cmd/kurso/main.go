// Command kurso is the entrypoint for the Kurso API monolith.
//
// It wires configuration, logging, and the driven adapters (Postgres, Redis)
// into the inbound HTTP adapter, then runs the server until an interrupt or
// termination signal triggers a graceful shutdown.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpapi "github.com/ANTON-IVANOVICH/kurso-v0/kurso-api/internal/adapter/http"
	"github.com/ANTON-IVANOVICH/kurso-v0/kurso-api/internal/adapter/store"
	"github.com/ANTON-IVANOVICH/kurso-v0/kurso-api/internal/config"
	"github.com/ANTON-IVANOVICH/kurso-v0/kurso-api/internal/platform/logger"
	"github.com/ANTON-IVANOVICH/kurso-v0/kurso-api/internal/platform/postgres"
	redisclient "github.com/ANTON-IVANOVICH/kurso-v0/kurso-api/internal/platform/redis"
	"github.com/ANTON-IVANOVICH/kurso-v0/kurso-api/internal/platform/server"
	"github.com/ANTON-IVANOVICH/kurso-v0/kurso-api/internal/service/auth"
	"github.com/ANTON-IVANOVICH/kurso-v0/kurso-api/internal/service/rates"
)

func main() {
	if err := run(); err != nil {
		slog.Error("kurso-api terminated", "err", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	log := logger.New(cfg.Log.Level, cfg.Log.Format)
	slog.SetDefault(log)
	log.Info("starting kurso-api", "env", cfg.Env, "http_port", cfg.HTTP.Port)

	// Root context cancelled on SIGINT/SIGTERM — drives graceful shutdown.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Driven adapter: Postgres. Connection is lazy; a failed startup ping is a
	// warning (not fatal) so the service boots for local development without
	// infra. Readiness (/readyz) reports the real state.
	pool, err := postgres.New(ctx, cfg.DB.URL)
	if err != nil {
		return fmt.Errorf("init postgres: %w", err)
	}
	defer pool.Close()
	if err := pingWithTimeout(ctx, 3*time.Second, pool.Ping); err != nil {
		log.Warn("postgres not reachable at startup, continuing", "err", err)
	} else {
		log.Info("postgres connected")
	}

	// Driven adapter: Redis.
	rdb, err := redisclient.New(cfg.Redis.URL)
	if err != nil {
		return fmt.Errorf("init redis: %w", err)
	}
	defer func() { _ = rdb.Close() }()
	if err := pingWithTimeout(ctx, 3*time.Second, func(c context.Context) error { return rdb.Ping(c).Err() }); err != nil {
		log.Warn("redis not reachable at startup, continuing", "err", err)
	} else {
		log.Info("redis connected")
	}

	// Application wiring: store → seed → rates service + SSE hub + runner.
	st := store.New(pool)
	if err := st.Seed(ctx); err != nil {
		log.Warn("catalogue seed failed, continuing", "err", err)
	} else {
		log.Info("catalogue seeded/verified")
	}
	if err := st.SeedGeo(ctx); err != nil {
		log.Warn("map geo seed failed, continuing", "err", err)
	}

	hub := rates.NewHub()
	svc := rates.NewService(st, rdb, hub)
	ticker := rates.NewTicker(st, svc, log, cfg.Rates.TickInterval)
	go ticker.Run(ctx)
	log.Info("rate runner started", "interval", cfg.Rates.TickInterval)

	// Auth: one service per identity store (admin + end user + merchant), distinct secrets.
	adminAuth := auth.NewService(adminAuthRepo{st}, cfg.Admin.JWTSecret, cfg.Admin.AccessTTL, cfg.Admin.RefreshTTL)
	userAuth := auth.NewService(userAuthRepo{st}, cfg.Admin.UserJWTSecret, cfg.Admin.AccessTTL, cfg.Admin.RefreshTTL)
	partnerAuth := auth.NewService(merchantAuthRepo{st}, cfg.Partner.JWTSecret, cfg.Admin.AccessTTL, cfg.Admin.RefreshTTL)

	// Ensure a seed administrator exists so a fresh DB can log in.
	if hash, err := auth.HashPassword(cfg.Admin.SeedPassword); err != nil {
		log.Warn("hash seed admin password failed", "err", err)
	} else if err := st.EnsureAdmin(ctx, cfg.Admin.SeedEmail, hash, "superadmin"); err != nil {
		log.Warn("seed admin failed, continuing", "err", err)
	} else {
		log.Info("admin account ensured", "email", cfg.Admin.SeedEmail)
	}

	// Ensure a seed merchant representative for the seed exchanger so the
	// partner cabinet has a working login on a fresh database.
	if hash, err := auth.HashPassword(cfg.Partner.SeedPassword); err != nil {
		log.Warn("hash seed partner password failed", "err", err)
	} else if err := st.EnsureExchangerUser(ctx, cfg.Partner.SeedExchanger, cfg.Partner.SeedEmail, hash, "owner"); err != nil {
		log.Warn("seed partner user failed, continuing", "err", err)
	} else {
		log.Info("partner account ensured", "email", cfg.Partner.SeedEmail, "exchanger", cfg.Partner.SeedExchanger)
	}

	// Inbound adapter: HTTP router.
	router := httpapi.NewRouter(httpapi.Deps{
		Log:            log,
		DB:             pool,
		Redis:          rdb,
		Svc:            svc,
		Auth:           adminAuth,
		UserAuth:       userAuth,
		PartnerAuth:    partnerAuth,
		Store:          st,
		AllowedOrigins: cfg.HTTP.AllowedOrigins,
		CookieSecure:   cfg.Admin.CookieSecure,
		CookieDomain:   cfg.Admin.CookieDomain,
	})

	srv := server.New(server.Config{
		Addr:            ":" + cfg.HTTP.Port,
		ReadTimeout:     cfg.HTTP.ReadTimeout,
		WriteTimeout:    cfg.HTTP.WriteTimeout,
		IdleTimeout:     cfg.HTTP.IdleTimeout,
		ShutdownTimeout: cfg.HTTP.ShutdownTimeout,
	}, router, log)

	return srv.Run(ctx)
}

// pingWithTimeout runs a ping function against a bounded child context.
func pingWithTimeout(ctx context.Context, d time.Duration, ping func(context.Context) error) error {
	ctx, cancel := context.WithTimeout(ctx, d)
	defer cancel()
	return ping(ctx)
}

// adminAuthRepo / userAuthRepo adapt the store's admin/user lookups to the
// auth.Repo port, so one auth.Service type serves both identity stores.
type adminAuthRepo struct{ st *store.Store }

func (r adminAuthRepo) AccountByEmail(ctx context.Context, email string) (auth.Account, error) {
	a, err := r.st.AdminByEmail(ctx, email)
	if err != nil {
		return auth.Account{}, err
	}
	return auth.Account{ID: a.ID, Email: a.Email, Role: string(a.Role), PasswordHash: a.PasswordHash, Status: a.Status}, nil
}
func (r adminAuthRepo) TouchLogin(ctx context.Context, id string) error {
	return r.st.TouchAdminLogin(ctx, id)
}

type userAuthRepo struct{ st *store.Store }

func (r userAuthRepo) AccountByEmail(ctx context.Context, email string) (auth.Account, error) {
	u, err := r.st.UserByEmail(ctx, email)
	if err != nil {
		return auth.Account{}, err
	}
	return auth.Account{ID: u.ID, Email: u.Email, Role: "user", PasswordHash: u.PasswordHash, Status: u.Status}, nil
}
func (r userAuthRepo) TouchLogin(ctx context.Context, id string) error {
	return r.st.TouchUserLogin(ctx, id)
}

type merchantAuthRepo struct{ st *store.Store }

func (r merchantAuthRepo) AccountByEmail(ctx context.Context, email string) (auth.Account, error) {
	u, err := r.st.ExchangerUserByEmail(ctx, email)
	if err != nil {
		return auth.Account{}, err
	}
	return auth.Account{ID: u.ID, Email: u.Email, Role: u.Role, PasswordHash: u.PasswordHash, Status: u.Status}, nil
}
func (r merchantAuthRepo) TouchLogin(ctx context.Context, id string) error {
	return r.st.TouchExchangerUserLogin(ctx, id)
}
