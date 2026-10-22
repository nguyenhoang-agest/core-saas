package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

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

// loadDotEnv reads environment files (e.g. .env) and sets any variables not already set in the process.
func loadDotEnv(paths ...string) {
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			if len(val) >= 2 && ((val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'')) {
				val = val[1 : len(val)-1]
			}
			if os.Getenv(key) == "" {
				os.Setenv(key, val)
			}
		}
	}
}

// Load loads configuration from an optional file and overlays environment variables.
func Load(filePath string) (*Config, error) {
	loadDotEnv(".env", "../.env")

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
	} else if appEnv := os.Getenv("APP_ENV"); appEnv != "" {
		cfg.Server.Env = appEnv
	}

	// Database overrides: accept full URL or granular connection parameters
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		cfg.Database.URL = dbURL
	} else if dbURL := os.Getenv("DB_URL"); dbURL != "" {
		cfg.Database.URL = dbURL
	} else if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		dbPort := os.Getenv("DB_PORT")
		if dbPort == "" {
			dbPort = "5432"
		}
		dbUser := os.Getenv("DB_USER")
		if dbUser == "" {
			dbUser = "postgres"
		}
		dbPass := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")
		if dbName == "" {
			dbName = "saas_dev"
		}
		sslMode := os.Getenv("DB_SSLMODE")
		if sslMode == "" {
			sslMode = "disable"
		}
		if dbPass != "" {
			cfg.Database.URL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", dbUser, dbPass, dbHost, dbPort, dbName, sslMode)
		} else {
			cfg.Database.URL = fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=%s", dbUser, dbHost, dbPort, dbName, sslMode)
		}
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
