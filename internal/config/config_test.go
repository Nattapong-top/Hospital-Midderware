package config_test

import (
	"os"
	"testing"

	"Hospital-Middleware/internal/config"
)

func TestConfig_Load_Defaults(t *testing.T) {
	// Clear env vars
	_ = os.Unsetenv("SERVER_PORT")
	_ = os.Unsetenv("JWT_SECRET")

	cfg := config.Load()
	if cfg.ServerPort != ":8080" {
		t.Errorf("expected default server port :8080, got %s", cfg.ServerPort)
	}
	if cfg.JWTSecret != "super-secret-key-5678" {
		t.Errorf("expected default jwt secret, got %s", cfg.JWTSecret)
	}
	if cfg.DSN() == "" {
		t.Error("expected non-empty DSN string")
	}
}

func TestConfig_Load_CustomEnv(t *testing.T) {
	_ = os.Setenv("SERVER_PORT", ":9090")
	_ = os.Setenv("DB_NAME", "custom_db")
	defer func() {
		_ = os.Unsetenv("SERVER_PORT")
		_ = os.Unsetenv("DB_NAME")
	}()

	cfg := config.Load()
	if cfg.ServerPort != ":9090" {
		t.Errorf("expected custom server port :9090, got %s", cfg.ServerPort)
	}
	if cfg.DBName != "custom_db" {
		t.Errorf("expected custom db name custom_db, got %s", cfg.DBName)
	}
}
