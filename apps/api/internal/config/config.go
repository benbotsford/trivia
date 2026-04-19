package config

import (
	"fmt"
	"os"
)

// Config holds all runtime configuration loaded from environment variables.
type Config struct {
	// HTTP
	Addr string // e.g. ":8080"

	// Postgres (Neon)
	DatabaseURL string

	// Redis
	RedisAddr     string
	RedisPassword string

	// Auth0
	Auth0Domain   string // e.g. "your-tenant.us.auth0.com"
	Auth0Audience string // API identifier in Auth0

	// Observability
	LogLevel string // "debug" | "info" | "warn" | "error"
}

// Load reads configuration from environment variables.
// Missing required vars cause an error.
func Load() (*Config, error) {
	c := &Config{
		Addr:          getEnv("ADDR", ":8080"),
		DatabaseURL:   mustGetEnv("DATABASE_URL"),
		RedisAddr:     getEnv("REDIS_ADDR", "redis:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		Auth0Domain:   mustGetEnv("AUTH0_DOMAIN"),
		Auth0Audience: mustGetEnv("AUTH0_AUDIENCE"),
		LogLevel:      getEnv("LOG_LEVEL", "info"),
	}
	return c, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("required env var %q is not set", key))
	}
	return v
}
