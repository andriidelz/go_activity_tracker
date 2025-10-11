package jobs

import (
	"log/slog"

	"github.com/andriidelzz/go-activity-tracker/internal/repository"
)

func AggregateEvents(repo *repository.Repository) {
	slog.Info("Running aggregation job...")
	if err := repo.AggregateLastPeriod(); err != nil {
		slog.Info("Aggregation failed:", "error", err)
	} else {
		slog.Info("Aggregation completed")
	}
}
