package repository

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/andriidelzz/go-activity-tracker/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		slog.Warn("DB_DSN not set, using default local connection")
		dsn = "host=localhost user=postgres password=postgres dbname=activity port=5432 sslmode=disable" // Of coz, in prod, you should not use this line of fallback code.
	}

	var database *gorm.DB
	var err error

	for attempt := 1; attempt <= 10; attempt++ {
		database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			slog.Info("Connected to PostgreSQL successfully")
			return database, nil
		}
		slog.Warn(fmt.Sprintf("DB connection attempt %d failed: %v", attempt, err))
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("could not connect to database after retries: %w", err)
}

func Migrate(database *gorm.DB) error {
	slog.Info("Running migrations...")

	err := database.AutoMigrate(
		&model.Event{},
		&model.Stat{},
	)
	if err != nil {
		return err
	}

	slog.Info("Migrations completed successfully")
	return nil
}
