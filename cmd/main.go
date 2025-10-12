package main

import (
	"log/slog"
	"os"

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

	repo := repository.NewRepository(db)
	handler := handler.NewHandler(repo)

	jobs.StartScheduler(repo)

	r := server.RegisterRoutes(handler)
	slog.Info("Starting API on :8080")
	r.Run(":8080")
}
