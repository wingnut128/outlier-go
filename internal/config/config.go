package config

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

// Config represents the application configuration
type Config struct {
	Logging LoggingConfig `toml:"logging"`
	Server  ServerConfig  `toml:"server"`
}

// LoggingConfig represents logging configuration
type LoggingConfig struct {
	Level  string `toml:"level"`  // trace, debug, info, warn, error
	Output string `toml:"output"` // stdout, stderr, file
	Format string `toml:"format"` // compact, pretty, json
}

// ServerConfig represents server configuration
type ServerConfig struct {
	BindIP string `toml:"bind_ip"`
	Port   int    `toml:"port"`
}

// DefaultConfig returns a configuration with default values
func DefaultConfig() *Config {
	return &Config{
		Logging: LoggingConfig{
			Level:  "info",
			Output: "stdout",
			Format: "compact",
		},
		Server: ServerConfig{
			Port:   3000,
			BindIP: "0.0.0.0",
		},
	}
}

// LoadConfig loads configuration from a file path
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config := DefaultConfig()
	if err := toml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return config, nil
}

// LoadConfigWithPriority loads configuration with the following priority:
// 1. Provided configPath (if not empty)
// 2. CONFIG_FILE environment variable
// 3. Default configuration
func LoadConfigWithPriority(configPath string) (*Config, error) {
	// Priority 1: Explicit path provided
	if configPath != "" {
		return LoadConfig(configPath)
	}

	// Priority 2: Environment variable
	envPath := os.Getenv("CONFIG_FILE")
	if envPath != "" {
		return LoadConfig(envPath)
	}

	// Priority 3: Default configuration
	return DefaultConfig(), nil
}
