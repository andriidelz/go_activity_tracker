package repository

import (
	"os"
	"testing"
	"time"

	"github.com/andriidelzz/go-activity-tracker/internal/model"
	"github.com/stretchr/testify/require"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupIntegrationDB(t *testing.T) *Repository {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		t.Skip("Skipping integration test: TEST_DATABASE_DSN not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	err = db.Migrator().DropTable(&model.Event{}, &model.Stat{})
	require.NoError(t, err)

	err = db.AutoMigrate(&model.Event{}, &model.Stat{})
	require.NoError(t, err)

	return NewRepository(db)
}

func TestIntegration_CreateAndGetEvent(t *testing.T) {
	repo := setupIntegrationDB(t)

	event := &model.Event{
		UserID:    42,
		Type:      "integration_login",
		Metadata:  map[string]any{"ip": "192.168.1.10"},
		CreatedAt: time.Now(),
	}

	err := repo.CreateEvent(event)
	require.NoError(t, err)
	require.NotZero(t, event.ID)

	events, err := repo.GetEvents(42)
	require.NoError(t, err)
	require.Len(t, events, 1)
	require.Equal(t, "integration_login", events[0].Type)
}

func TestIntegration_AggregationFlow(t *testing.T) {
	repo := setupIntegrationDB(t)

	events := []model.Event{
		{UserID: 1, Type: "click", CreatedAt: time.Now()},
		{UserID: 1, Type: "scroll", CreatedAt: time.Now()},
		{UserID: 2, Type: "click", CreatedAt: time.Now()},
	}
	for _, e := range events {
		require.NoError(t, repo.CreateEvent(&e))
	}

	err := repo.AggregateLastPeriod()
	require.NoError(t, err)

	stats, err := repo.GetStats()
	require.NoError(t, err)
	require.Len(t, stats, 2)

	require.Equal(t, 2, stats[0].EventCount)
	require.Equal(t, 1, stats[1].EventCount)
}

func TestIntegration_NoEvents(t *testing.T) {
	repo := setupIntegrationDB(t)

	// Aggregate when there are no events
	err := repo.AggregateLastPeriod()
	require.NoError(t, err)

	stats, err := repo.GetStats()
	require.NoError(t, err)
	require.Len(t, stats, 0)
}

func TestIntegration_MultipleAggregations(t *testing.T) {
	repo := setupIntegrationDB(t)

	events1 := []model.Event{
		{UserID: 1, Type: "click", CreatedAt: time.Now()},
		{UserID: 1, Type: "scroll", CreatedAt: time.Now()},
	}
	for _, e := range events1 {
		require.NoError(t, repo.CreateEvent(&e))
	}

	err := repo.AggregateLastPeriod()
	require.NoError(t, err)

	stats, err := repo.GetStats()
	require.NoError(t, err)
	require.Len(t, stats, 1)
	require.Equal(t, 2, stats[0].EventCount)

	events2 := []model.Event{
		{UserID: 1, Type: "click", CreatedAt: time.Now()},
		{UserID: 2, Type: "scroll", CreatedAt: time.Now()},
	}
	for _, e := range events2 {
		require.NoError(t, repo.CreateEvent(&e))
	}

	err = repo.AggregateLastPeriod()
	require.NoError(t, err)

	stats, err = repo.GetStats()
	require.NoError(t, err)
	require.Len(t, stats, 2)

	m := make(map[int]int)
	for _, s := range stats {
		m[s.UserID] = s.EventCount
	}

	require.Equal(t, 3, m[1])
	require.Equal(t, 1, m[2])
}
