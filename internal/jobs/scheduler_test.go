package jobs

import (
	"testing"
	"time"

	"github.com/andriidelzz/go-activity-tracker/internal/model"
	"github.com/stretchr/testify/require"
)

type MockRepository struct {
	Called bool
}

func (m *MockRepository) CreateEvent(event *model.Event) error        { return nil }
func (m *MockRepository) GetEvents(userID int) ([]model.Event, error) { return nil, nil }
func (m *MockRepository) AggregateLastPeriod() error {
	m.Called = true
	return nil
}
func (m *MockRepository) GetStats() ([]model.Stat, error) { return nil, nil }

func TestAggregateEvents(t *testing.T) {
	mockRepo := &MockRepository{}
	AggregateEvents(mockRepo)
	require.True(t, mockRepo.Called)
}

func TestScheduler_StartScheduler(t *testing.T) {
	mockRepo := &MockRepository{}
	c := StartScheduler(mockRepo, true)
	time.Sleep(20 * time.Millisecond)
	c.Stop()
	require.NotNil(t, c)
}

func BenchmarkAggregateEvents(b *testing.B) {
	mockRepo := &MockRepository{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AggregateEvents(mockRepo)
	}
}
