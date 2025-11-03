package repository

import (
	"context"
	"testing"
	"time"

	"github.com/andriidelzz/go-activity-tracker/internal/model"
	"github.com/stretchr/testify/require"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *Repository {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&model.Event{}, &model.Stat{})
	require.NoError(t, err)

	return NewRepository(db)
}

func TestCreateEvent(t *testing.T) {
	repo := setupTestDB(t)

	ctx := context.Background()

	event := &model.Event{
		UserID:    1,
		Type:      "login",
		Metadata:  map[string]any{"ip": "127.0.0.1"},
		CreatedAt: time.Now(),
	}

	err := repo.CreateEvent(ctx, event)
	require.NoError(t, err)
	require.NotZero(t, event.ID)
}

func TestGetEvents(t *testing.T) {
	repo := setupTestDB(t)

	ctx := context.Background()

	events := []model.Event{
		{UserID: 1, Type: "podcast"},
		{UserID: 1, Type: "circus"},
		{UserID: 2, Type: "concert of Baidak in Kyiv"},
	}

	for _, e := range events {
		err := repo.CreateEvent(ctx, &e)
		require.NoError(t, err)
	}

	result, err := repo.GetEvents(ctx, 1)
	require.NoError(t, err)
	require.Len(t, result, 2)
	require.Equal(t, "podcast", result[0].Type)
	require.Equal(t, 1, result[0].UserID)
}

func TestAggregateLastPeriod(t *testing.T) {
	repo := setupTestDB(t)

	ctx := context.Background()

	events := []model.Event{
		{UserID: 1, Type: "click", CreatedAt: time.Now()},
		{UserID: 1, Type: "scroll", CreatedAt: time.Now()},
		{UserID: 2, Type: "click", CreatedAt: time.Now()},
	}

	for _, e := range events {
		require.NoError(t, repo.CreateEvent(ctx, &e))
	}

	err := repo.AggregateLastPeriod(ctx)
	require.NoError(t, err)

	stats, err := repo.GetStats(ctx)
	require.NoError(t, err)

	require.Len(t, stats, 2)
	require.Equal(t, 2, stats[0].EventCount)
	require.Equal(t, 1, stats[1].EventCount)
}

func TestGetStats(t *testing.T) {
	repo := setupTestDB(t)

	ctx := context.Background()

	statsToInsert := []model.Stat{
		{UserID: 1, EventCount: 5},
		{UserID: 2, EventCount: 3},
	}

	for _, s := range statsToInsert {
		require.NoError(t, repo.db.WithContext(ctx).Create(&s).Error)
	}

	stats, err := repo.GetStats(ctx)
	require.NoError(t, err)
	require.Len(t, stats, 2)

	require.Equal(t, 1, stats[0].UserID)
	require.Equal(t, 5, stats[0].EventCount)
	require.Equal(t, 2, stats[1].UserID)
	require.Equal(t, 3, stats[1].EventCount)
}

func BenchmarkCreateEvent(b *testing.B) {
	repo := setupTestDB(nil)

	ctx := context.Background()

	event := &model.Event{
		UserID:   1,
		Type:     "benchmark",
		Metadata: map[string]any{"ip": "127.0.0.1"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = repo.CreateEvent(ctx, event)
	}
}

func BenchmarkAggregateLastPeriod(b *testing.B) {
	repo := setupTestDB(nil)

	ctx := context.Background()

	// Preload DB with events
	for i := range 10 {
		_ = repo.CreateEvent(ctx, &model.Event{
			UserID:   i % 10,
			Type:     "click",
			Metadata: map[string]any{},
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = repo.AggregateLastPeriod(ctx)
	}
}
