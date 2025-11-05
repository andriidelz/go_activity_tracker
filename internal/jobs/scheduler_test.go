package jobs

import (
	"context"
	"testing"
	"time"

	"github.com/andriidelzz/go-activity-tracker/internal/model"
	"github.com/stretchr/testify/require"
)

type MockRepository struct {
	Called bool
}

func (m *MockRepository) CreateEvent(_ context.Context, _ *model.Event) error { return nil }
func (m *MockRepository) GetEvents(_ context.Context, _ int) ([]model.Event, error) {
	return nil, nil
}
func (m *MockRepository) AggregateLastPeriod(_ context.Context) error {
	m.Called = true
	return nil
}
func (m *MockRepository) GetStats(_ context.Context) ([]model.Stat, error) { return nil, nil }

func TestAggregateEvents(t *testing.T) {
	mockRepo := &MockRepository{}
	ctx := context.Background()
	AggregateEvents(ctx, mockRepo)
	require.True(t, mockRepo.Called)
}

func TestScheduler_StartScheduler(t *testing.T) {
	mockRepo := &MockRepository{}
	ctx := context.Background()
	c := StartScheduler(ctx, mockRepo, true)
	time.Sleep(20 * time.Millisecond)
	c.Stop()
	require.NotNil(t, c)
}

func BenchmarkAggregateEvents(b *testing.B) {
	mockRepo := &MockRepository{}
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AggregateEvents(ctx, mockRepo)
	}
}
