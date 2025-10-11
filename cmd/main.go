package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/andriidelzz/go-activity-tracker/internal/api"
	"github.com/andriidelzz/go-activity-tracker/internal/db"
	"github.com/andriidelzz/go-activity-tracker/internal/jobs"

	"github.com/robfig/cron/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		slog.Info("DB_DSN not set, using default DSN")
		dsn = "host=localhost user=postgres password=postgres dbname=activity port=5432 sslmode=disable"
	} else {
		slog.Info("Using DB_DSN from environment", "dsn", dsn)
	}

	var gormDB *gorm.DB
	var err error
	for i := range 10 {
		gormDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			slog.Info("Successfully connected to database")
			break
		}
		slog.Info(fmt.Sprintf("Failed to connect to DB (attempt %d): %v", i+1, err))
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		slog.Error("Failed to connect to DB after retries:", "error", err)
		os.Exit(1)
	}

	slog.Info("Running database migrations")
	db.Migrate(gormDB)
	slog.Info("Database migrations completed")

	repo := db.NewRepository(gormDB)

	c := cron.New()
	_, err = c.AddFunc("0 */4 * * *", func() { jobs.AggregateEvents(repo) }) // Every 4 hours.
	if err != nil {
		slog.Error("Critical error occurred", "err", err)
		os.Exit(1)
	}
	c.Start()
	slog.Info("Cron job started for aggregation every 4 hours")

	router := api.SetupRouter(repo)
	slog.Info("Starting API on :8080")
	router.Run(":8080")
}
