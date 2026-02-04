package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Logging.Level != "info" {
		t.Errorf("expected logging level 'info', got '%s'", cfg.Logging.Level)
	}
	if cfg.Logging.Output != "stdout" {
		t.Errorf("expected logging output 'stdout', got '%s'", cfg.Logging.Output)
	}
	if cfg.Logging.Format != "compact" {
		t.Errorf("expected logging format 'compact', got '%s'", cfg.Logging.Format)
	}
	if cfg.Server.Port != 3000 {
		t.Errorf("expected server port 3000, got %d", cfg.Server.Port)
	}
	if cfg.Server.BindIP != "0.0.0.0" {
		t.Errorf("expected bind IP '0.0.0.0', got '%s'", cfg.Server.BindIP)
	}
}

func TestLoadConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.toml")

	content := `[logging]
level = "debug"
output = "stderr"
format = "json"

[server]
port = 8080
bind_ip = "127.0.0.1"
`
	err := os.WriteFile(configFile, []byte(content), 0o644)
	if err != nil {
		t.Fatalf("failed to create test config file: %v", err)
	}

	cfg, err := LoadConfig(configFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Logging.Level != "debug" {
		t.Errorf("expected logging level 'debug', got '%s'", cfg.Logging.Level)
	}
	if cfg.Logging.Output != "stderr" {
		t.Errorf("expected logging output 'stderr', got '%s'", cfg.Logging.Output)
	}
	if cfg.Logging.Format != "json" {
		t.Errorf("expected logging format 'json', got '%s'", cfg.Logging.Format)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("expected server port 8080, got %d", cfg.Server.Port)
	}
	if cfg.Server.BindIP != "127.0.0.1" {
		t.Errorf("expected bind IP '127.0.0.1', got '%s'", cfg.Server.BindIP)
	}
}

func TestLoadConfig_InvalidFile(t *testing.T) {
	_, err := LoadConfig("/nonexistent/config.toml")
	if err == nil {
		t.Error("expected error for nonexistent file, got nil")
	}
}

func TestLoadConfig_InvalidTOML(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.toml")

	content := `this is not valid toml`
	err := os.WriteFile(configFile, []byte(content), 0o644)
	if err != nil {
		t.Fatalf("failed to create test config file: %v", err)
	}

	_, err = LoadConfig(configFile)
	if err == nil {
		t.Error("expected error for invalid TOML, got nil")
	}
}

func TestLoadConfigWithPriority_ExplicitPath(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.toml")

	content := `[server]
port = 9999
`
	err := os.WriteFile(configFile, []byte(content), 0o644)
	if err != nil {
		t.Fatalf("failed to create test config file: %v", err)
	}

	cfg, err := LoadConfigWithPriority(configFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Server.Port != 9999 {
		t.Errorf("expected port 9999, got %d", cfg.Server.Port)
	}
}

func TestLoadConfigWithPriority_EnvVar(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.toml")

	content := `[server]
port = 7777
`
	err := os.WriteFile(configFile, []byte(content), 0o644)
	if err != nil {
		t.Fatalf("failed to create test config file: %v", err)
	}

	os.Setenv("CONFIG_FILE", configFile)
	defer os.Unsetenv("CONFIG_FILE")

	cfg, err := LoadConfigWithPriority("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Server.Port != 7777 {
		t.Errorf("expected port 7777, got %d", cfg.Server.Port)
	}
}

func TestLoadConfigWithPriority_Default(t *testing.T) {
	os.Unsetenv("CONFIG_FILE")

	cfg, err := LoadConfigWithPriority("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Server.Port != 3000 {
		t.Errorf("expected default port 3000, got %d", cfg.Server.Port)
	}
}
