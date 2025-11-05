package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/andriidelzz/go-activity-tracker/internal/handler"
	"github.com/andriidelzz/go-activity-tracker/internal/jobs"
	"github.com/andriidelzz/go-activity-tracker/internal/metrics"
	"github.com/andriidelzz/go-activity-tracker/internal/repository"
	"github.com/andriidelzz/go-activity-tracker/internal/server"
)

func main() {
	db, err := repository.Connect()
	if err != nil {
		slog.Error("Database connection failed", "error", err)
		os.Exit(1)
	}

	if err := repository.Migrate(db); err != nil {
		slog.Error("Migration failed", "error", err)
		os.Exit(1)
	}

	metrics.Register()
	metrics.CollectSystemMetrics(db)

	repo := repository.NewRepository(db)
	handler := handler.NewHandler(repo)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start scheduler with context
	jobs.StartScheduler(ctx, repo, false)

	// Shutdown handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
		slog.Info("Shutting down...")
		os.Exit(0)
	}()

	r := server.RegisterRoutes(handler)
	slog.Info("Starting API on :8080")
	if err := r.Run(":8080"); err != nil {
		slog.Error("API server failed", "error", err)
		os.Exit(1)
	}
}
