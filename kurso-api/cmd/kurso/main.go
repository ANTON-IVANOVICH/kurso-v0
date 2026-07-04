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

	hub := rates.NewHub()
	svc := rates.NewService(st, rdb, hub)
	ticker := rates.NewTicker(st, svc, log, cfg.Rates.TickInterval)
	go ticker.Run(ctx)
	log.Info("rate runner started", "interval", cfg.Rates.TickInterval)

	// Admin auth: ensure a seed administrator exists so a fresh DB can log in.
	authSvc := auth.NewService(st, cfg.Admin.JWTSecret, cfg.Admin.AccessTTL, cfg.Admin.RefreshTTL)
	if hash, err := auth.HashPassword(cfg.Admin.SeedPassword); err != nil {
		log.Warn("hash seed admin password failed", "err", err)
	} else if err := st.EnsureAdmin(ctx, cfg.Admin.SeedEmail, hash, "superadmin"); err != nil {
		log.Warn("seed admin failed, continuing", "err", err)
	} else {
		log.Info("admin account ensured", "email", cfg.Admin.SeedEmail)
	}

	// Inbound adapter: HTTP router.
	router := httpapi.NewRouter(httpapi.Deps{
		Log:            log,
		DB:             pool,
		Redis:          rdb,
		Svc:            svc,
		Auth:           authSvc,
		AllowedOrigins: cfg.HTTP.AllowedOrigins,
		CookieSecure:   cfg.Admin.CookieSecure,
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
