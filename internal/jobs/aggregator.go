package jobs

import (
	"context"
	"log/slog"
	"time"

	"github.com/andriidelzz/go-activity-tracker/internal/repository"
)

func AggregateEvents(ctx context.Context, repo repository.RepositoryInterface) {
	aggCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	slog.Info("Running aggregation job...")
	if err := repo.AggregateLastPeriod(aggCtx); err != nil {
		slog.Info("Aggregation failed:", "error", err)
	} else {
		slog.Info("Aggregation completed")
	}
}
