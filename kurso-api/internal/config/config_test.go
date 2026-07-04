package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	// Force the env to empty so the built-in defaults apply deterministically.
	for _, k := range []string{"KURSO_ENV", "HTTP_PORT", "DATABASE_URL", "REDIS_URL", "LOG_LEVEL", "LOG_FORMAT"} {
		t.Setenv(k, "")
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}
	if cfg.Env != "dev" {
		t.Errorf("Env = %q, want dev", cfg.Env)
	}
	if cfg.HTTP.Port != "8080" {
		t.Errorf("HTTP.Port = %q, want 8080", cfg.HTTP.Port)
	}
	// dev defaults to text logs.
	if cfg.Log.Format != "text" {
		t.Errorf("Log.Format = %q, want text", cfg.Log.Format)
	}
	if cfg.DB.URL == "" || cfg.Redis.URL == "" {
		t.Error("DB/Redis URLs should have non-empty defaults")
	}
}

func TestLoadProductionLogFormat(t *testing.T) {
	t.Setenv("KURSO_ENV", "production")
	t.Setenv("LOG_FORMAT", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}
	if cfg.Log.Format != "json" {
		t.Errorf("Log.Format = %q, want json in production", cfg.Log.Format)
	}
}

func TestLoadInvalidPort(t *testing.T) {
	t.Setenv("HTTP_PORT", "not-a-number")
	if _, err := Load(); err == nil {
		t.Error("Load() should fail on a non-numeric HTTP_PORT")
	}
}
