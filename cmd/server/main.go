package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/falasefemi2/workaround-backend/internal/config"
	"github.com/falasefemi2/workaround-backend/internal/server"
	"github.com/jackc/pgx/v5/pgxpool"
)

const shutdownTimeout = 30 * time.Second

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	// Setup structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Connect to database
	poolCfg, err := pgxpool.ParseConfig(cfg.Database.URL)
	if err != nil {
		slog.Error("failed to parse database config", "error", err)
		os.Exit(1)
	}

	// Apply DB pool configuration

	poolCfg.MaxConns = int32(cfg.Database.MaxOpenConns)
	pool, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		slog.Error("failed to create pool", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Verify DB connection
	if err := pool.Ping(context.Background()); err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	slog.Info("database connected")

	// Initialize router
	handler := server.New(pool, cfg)

	// Create HTTP server using config timeouts
	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:           handler,
		ReadTimeout:       cfg.Server.ReadTimeout,
		WriteTimeout:      cfg.Server.WriteTimeout,
		IdleTimeout:       cfg.Server.IdleTimeout,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Listen for shutdown signals
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start server
	go func() {
		slog.Info("starting server",
			"port", cfg.Server.Port,
			"env", cfg.Primary.Env,
		)

		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt
	<-ctx.Done()
	slog.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
		os.Exit(1)
	}

	slog.Info("server exited cleanly")
}
