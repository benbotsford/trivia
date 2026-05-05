// Package config loads runtime configuration from environment variables.
package config

import (
	"fmt"
	"os"
	"time"
)

// Config holds all runtime configuration for the API server.
type Config struct {
	// HTTP server listen address (e.g. ":8080")
	Addr string

	// Postgres connection string — format: postgres://user:pass@host/dbname
	DatabaseURL string

	// Redis address (host:port) and optional password
	RedisAddr     string
	RedisPassword string

	// Auth0 tenant domain and API audience identifier.
	// Optional when DevAuthToken is set — leave blank during local development
	// to skip JWT validation entirely and rely solely on the dev bypass.
	Auth0Domain   string
	Auth0Audience string

	// DevAuthToken is a shared secret used to bypass JWT validation in local
	// development. Any request bearing "Authorization: Bearer <DevAuthToken>"
	// is treated as a hardcoded dev user without contacting Auth0.
	// NEVER set this in production — leave it blank or unset.
	DevAuthToken string

	// Minimum log level: "debug" | "info" | "warn" | "error"
	LogLevel string

	// AutoMigrate runs embedded goose migrations on startup when true.
	// Leave off in environments where migrations are managed out-of-band
	// (e.g. a dedicated ops step before a rolling deploy).
	AutoMigrate bool

	// BootstrapSamples seeds sample question banks and questions on startup
	// when true. Intended for local development; leave off in production.
	BootstrapSamples bool

	// --- Connection pool ---
	// Neon's lower tiers cap concurrent connections in the low double digits.
	// Use the Neon pooler host (*.pooler.neon.tech) in DATABASE_URL for the
	// app, and the direct host only for migrations.

	// DBMaxConns caps the total number of open connections in the pool.
	// Default: 10 — safe for Neon Starter/Launch. Raise if using the pooler.
	DBMaxConns int32

	// DBMinConns is the minimum connections to keep open when idle.
	// 0 lets the pool drain fully — preferred for Neon serverless endpoints
	// which sleep after inactivity and don't need warm standby connections.
	DBMinConns int32

	// DBMaxConnLifetime is the maximum age of a connection before it is
	// recycled. Prevents accumulating stale connections across Neon restarts.
	DBMaxConnLifetime time.Duration

	// DBMaxConnIdleTime closes connections that have been idle longer than
	// this. Keep below Neon's idle-connection timeout (~5 min) so the pool
	// proactively releases connections before Neon drops them, avoiding
	// "connection reset" errors on the next request after a quiet period.
	DBMaxConnIdleTime time.Duration
}

// Load reads configuration from environment variables.
// Only DATABASE_URL is strictly required. Auth0 vars are optional when
// DEV_AUTH_TOKEN is set for local development.
func Load() (*Config, error) {
	c := &Config{
		Addr:          getEnv("ADDR", ":8080"),
		DatabaseURL:   mustGetEnv("DATABASE_URL"),
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		Auth0Domain:   getEnv("AUTH0_DOMAIN", ""),
		Auth0Audience: getEnv("AUTH0_AUDIENCE", ""),
		DevAuthToken:     getEnv("DEV_AUTH_TOKEN", ""),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
		AutoMigrate:      getEnvBool("AUTO_MIGRATE", false),
		BootstrapSamples: getEnvBool("BOOTSTRAP_SAMPLES", false),

		DBMaxConns:        getEnvInt32("DB_MAX_CONNS", 10),
		DBMinConns:        getEnvInt32("DB_MIN_CONNS", 0),
		DBMaxConnLifetime: getEnvDuration("DB_MAX_CONN_LIFETIME", 30*time.Minute),
		DBMaxConnIdleTime: getEnvDuration("DB_MAX_CONN_IDLE_TIME", 3*time.Minute),
	}
	return c, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	switch os.Getenv(key) {
	case "1", "true", "TRUE", "True", "yes", "YES":
		return true
	case "0", "false", "FALSE", "False", "no", "NO":
		return false
	default:
		return fallback
	}
}

func mustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("required env var %q is not set", key))
	}
	return v
}

func getEnvInt32(key string, fallback int32) int32 {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	var n int32
	if _, err := fmt.Sscanf(v, "%d", &n); err != nil || n < 0 {
		return fallback
	}
	return n
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}
	return d
}
