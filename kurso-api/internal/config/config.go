// Package config loads runtime configuration from environment variables.
//
// Every value has a development-friendly default so the binary runs with an
// empty environment; production overrides via env vars (or a mounted .env).
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config is the fully-resolved application configuration.
type Config struct {
	Env   string
	HTTP  HTTPConfig
	DB    DBConfig
	Redis RedisConfig
	Log   LogConfig
	Rates RatesConfig
}

// HTTPConfig holds HTTP server tuning.
type HTTPConfig struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
	AllowedOrigins  []string
}

// RatesConfig tunes the rate runner (stand-in for the parser pipeline).
type RatesConfig struct {
	TickInterval time.Duration
}

// DBConfig holds the Postgres connection string.
type DBConfig struct {
	URL string
}

// RedisConfig holds the Redis connection string.
type RedisConfig struct {
	URL string
}

// LogConfig controls the slog handler.
type LogConfig struct {
	Level  string // debug|info|warn|error
	Format string // json|text
}

// Load reads configuration from the environment and validates it.
func Load() (Config, error) {
	env := getenv("KURSO_ENV", "dev")

	// In local/dev, default to human-readable text logs; JSON otherwise.
	logFormat := "json"
	switch env {
	case "dev", "development", "local":
		logFormat = "text"
	}

	cfg := Config{
		Env: env,
		HTTP: HTTPConfig{
			Port:            getenv("HTTP_PORT", "8080"),
			ReadTimeout:     getDuration("HTTP_READ_TIMEOUT", 10*time.Second),
			WriteTimeout:    getDuration("HTTP_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:     getDuration("HTTP_IDLE_TIMEOUT", 60*time.Second),
			ShutdownTimeout: getDuration("HTTP_SHUTDOWN_TIMEOUT", 10*time.Second),
			AllowedOrigins:  splitCSV(getenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")),
		},
		Rates: RatesConfig{
			TickInterval: getDuration("RATES_TICK_INTERVAL", 5*time.Second),
		},
		DB: DBConfig{
			URL: getenv("DATABASE_URL", "postgres://kurso:kurso@localhost:5432/kurso?sslmode=disable"),
		},
		Redis: RedisConfig{
			URL: getenv("REDIS_URL", "redis://localhost:6379/0"),
		},
		Log: LogConfig{
			Level:  getenv("LOG_LEVEL", "info"),
			Format: getenv("LOG_FORMAT", logFormat),
		},
	}

	if _, err := strconv.Atoi(cfg.HTTP.Port); err != nil {
		return Config{}, fmt.Errorf("invalid HTTP_PORT %q: %w", cfg.HTTP.Port, err)
	}
	if cfg.DB.URL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL must not be empty")
	}
	if cfg.Redis.URL == "" {
		return Config{}, fmt.Errorf("REDIS_URL must not be empty")
	}

	return cfg, nil
}

func getenv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

// splitCSV parses a comma-separated env value into a trimmed, non-empty slice.
func splitCSV(v string) []string {
	var out []string
	for _, part := range strings.Split(v, ",") {
		if s := strings.TrimSpace(part); s != "" {
			out = append(out, s)
		}
	}
	return out
}

func getDuration(key string, fallback time.Duration) time.Duration {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}
	return d
}
