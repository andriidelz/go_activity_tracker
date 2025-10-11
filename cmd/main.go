package main

import (
	"log/slog"
	"os"

	"github.com/andriidelzz/go-activity-tracker/internal/handler"
	"github.com/andriidelzz/go-activity-tracker/internal/jobs"
	"github.com/andriidelzz/go-activity-tracker/internal/metrics"
	"github.com/andriidelzz/go-activity-tracker/internal/repository"
	"github.com/andriidelzz/go-activity-tracker/internal/server"

	"github.com/robfig/cron/v3"
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

	c := cron.New()
	_, err = c.AddFunc("0 */4 * * *", func() { jobs.AggregateEvents(repo) }) // Every 4 hours.
	if err != nil {
		slog.Error("Critical error occurred during cron start", "err", err)
		os.Exit(1)
	} // todo move to jobs

	c.Start() // todo add stop logic if not exists.
	slog.Info("Cron job started for aggregation every 4 hours")

	r := server.RegisterRoutes(handler)
	slog.Info("Starting API on :8080")
	r.Run(":8080")
}
