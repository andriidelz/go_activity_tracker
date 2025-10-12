package jobs

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/andriidelzz/go-activity-tracker/internal/repository"
	"github.com/robfig/cron/v3"
)

func StartScheduler(repo *repository.Repository) {
	c := cron.New()

	_, err := c.AddFunc("0 */4 * * *", func() {
		AggregateEvents(repo)
	})
	if err != nil {
		slog.Error("Critical error occurred during cron start", "err", err)
		os.Exit(1)
	}

	c.Start()
	slog.Info("Cron job started for aggregation every 4 hours")

	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop

		slog.Info("Shutdown signal received. Stopping cron scheduler...")

		c.Stop()

		slog.Info("Cron scheduler stopped gracefully.")
	}()
}
