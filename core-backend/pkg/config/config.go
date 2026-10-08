package config

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

// Config represents root application configuration.
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Logger   LoggerConfig   `yaml:"logger"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
	Env  string `yaml:"env"`
}

type DatabaseConfig struct {
	URL      string `yaml:"url"`
	MaxConns int32  `yaml:"max_conns"`
}

type LoggerConfig struct {
	Level string `yaml:"level"`
}

// Default returns sensible defaults for development.
func Default() *Config {
	return &Config{
		Server: ServerConfig{
			Port: "8080",
			Env:  "development",
		},
		Database: DatabaseConfig{
			URL:      "postgres://postgres:postgres@localhost:5432/saas_dev?sslmode=disable",
			MaxConns: 10,
		},
		Logger: LoggerConfig{
			Level: "info",
		},
	}
}

// Load loads configuration from an optional file and overlays environment variables.
func Load(filePath string) (*Config, error) {
	cfg := Default()

	if filePath != "" {
		data, err := os.ReadFile(filePath)
		if err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf("reading config file %s: %w", filePath, err)
		}
		if err == nil {
			if err := yaml.Unmarshal(data, cfg); err != nil {
				return nil, fmt.Errorf("parsing config file: %w", err)
			}
		}
	}

	// Environment variable overrides
	if port := os.Getenv("PORT"); port != "" {
		cfg.Server.Port = port
	}
	if env := os.Getenv("ENV"); env != "" {
		cfg.Server.Env = env
	}
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		cfg.Database.URL = dbURL
	}
	if maxConns := os.Getenv("DB_MAX_CONNS"); maxConns != "" {
		if val, err := strconv.ParseInt(maxConns, 10, 32); err == nil {
			cfg.Database.MaxConns = int32(val)
		}
	}
	if lvl := os.Getenv("LOG_LEVEL"); lvl != "" {
		cfg.Logger.Level = lvl
	}

	return cfg, nil
}
