package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"

	"github.com/btbots1994/trivia/internal/auth"
	"github.com/btbots1994/trivia/internal/billing"
	"github.com/btbots1994/trivia/internal/config"
	"github.com/btbots1994/trivia/internal/game"
	"github.com/btbots1994/trivia/internal/realtime"
	"github.com/btbots1994/trivia/internal/store"
	"github.com/btbots1994/trivia/internal/user"
)

func main() {
	// -------------------------------------------------------------------------
	// Logger
	// -------------------------------------------------------------------------
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo, // overridden below after config is loaded
	}))
	slog.SetDefault(logger)

	// -------------------------------------------------------------------------
	// Config
	// -------------------------------------------------------------------------
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "err", err)
		os.Exit(1)
	}

	// -------------------------------------------------------------------------
	// Postgres
	// -------------------------------------------------------------------------
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to postgres", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		slog.Error("postgres ping failed", "err", err)
		os.Exit(1)
	}
	slog.Info("postgres connected")

	// -------------------------------------------------------------------------
	// Redis
	// -------------------------------------------------------------------------
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
	})
	defer rdb.Close()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		slog.Error("redis ping failed", "err", err)
		os.Exit(1)
	}
	slog.Info("redis connected")

	// -------------------------------------------------------------------------
	// Services
	// -------------------------------------------------------------------------
	queries := store.New(pool)
	_ = user.New(queries)
	entitlements := billing.NoopChecker{}
	gameSvc := game.New(queries, entitlements)
	hub := realtime.New()

	// -------------------------------------------------------------------------
	// Auth middleware
	// -------------------------------------------------------------------------
	authMiddleware := auth.New(cfg.Auth0Domain, cfg.Auth0Audience)

	// -------------------------------------------------------------------------
	// Router
	// -------------------------------------------------------------------------
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	// Health / readiness probes
	r.Get("/healthz", healthz(pool, rdb))
	r.Get("/readyz", healthz(pool, rdb))

	// Prometheus metrics
	r.Handle("/metrics", promhttp.Handler())

	// WebSocket (unauthenticated — players join with game code + display name)
	hub.RegisterRoutes(r)

	// Authenticated API routes
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.Handler)
		gameSvc.RegisterRoutes(r)
	})

	// -------------------------------------------------------------------------
	// HTTP server with graceful shutdown
	// -------------------------------------------------------------------------
	srv := &http.Server{
		Addr:         cfg.Addr,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("server starting", "addr", cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	<-quit
	slog.Info("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("graceful shutdown failed", "err", err)
	}
	slog.Info("server stopped")
}

// healthz checks both Postgres and Redis connectivity.
func healthz(pool *pgxpool.Pool, rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		if err := pool.Ping(ctx); err != nil {
			http.Error(w, "postgres unhealthy: "+err.Error(), http.StatusServiceUnavailable)
			return
		}
		if _, err := rdb.Ping(ctx).Result(); err != nil {
			http.Error(w, "redis unhealthy: "+err.Error(), http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}
}
