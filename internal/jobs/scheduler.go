package jobs

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/andriidelzz/go-activity-tracker/internal/repository"
	"github.com/robfig/cron/v3"
)

func StartScheduler(ctx context.Context, repo repository.RepositoryInterface, testMode bool) *cron.Cron {
	c := cron.New()

	schedule := "0 */4 * * *"
	if testMode {
		schedule = "@every 1s"
	}

	_, err := c.AddFunc(schedule, func() {
		// Create a child context with timeout for each aggregation run
		aggCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
		defer cancel()
		AggregateEvents(aggCtx, repo)
	})
	if err != nil {
		slog.Error("Critical error occurred during cron start", "err", err)
		os.Exit(1)
	}

	c.Start()
	slog.Info("Cron job started", "schedule", schedule)

	if !testMode {
		go func() {
			// Listen for OS signals or context cancellation
			stop := make(chan os.Signal, 1)
			signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
			select {
			case <-stop:
				slog.Info("Shutdown signal received.")
			case <-ctx.Done():
				slog.Info("Context cancellation received.")
			}
			slog.Info("Stopping cron scheduler...")
			c.Stop()
			slog.Info("Cron scheduler stopped gracefully.")
		}()
	}

	return c
}
