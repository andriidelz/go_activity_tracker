package repository

import (
	"context"

	"github.com/andriidelzz/go-activity-tracker/internal/model"

	"gorm.io/gorm"
)

type RepositoryInterface interface {
	CreateEvent(ctx context.Context, event *model.Event) error
	GetEvents(ctx context.Context, userID int) ([]model.Event, error)
	AggregateLastPeriod(ctx context.Context) error
	GetStats(ctx context.Context) ([]model.Stat, error)
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

type Repository struct {
	db *gorm.DB
}

func (r *Repository) CreateEvent(ctx context.Context, event *model.Event) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *Repository) GetEvents(ctx context.Context, userID int) ([]model.Event, error) {
	var events []model.Event
	return events, r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&events).Error
}

func (r *Repository) AggregateLastPeriod(ctx context.Context) error {
	var results []struct {
		UserID     int
		EventCount int
	}
	if err := r.db.WithContext(ctx).Model(&model.Event{}).
		Select("user_id, COUNT(*) as event_count").
		Group("user_id").
		Scan(&results).Error; err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Exec("DELETE FROM stats").Error; err != nil {
		return err
	}

	for _, res := range results {
		stat := model.Stat{
			UserID:     res.UserID,
			EventCount: res.EventCount,
		}
		if err := r.db.WithContext(ctx).Create(&stat).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) GetStats(ctx context.Context) ([]model.Stat, error) {
	var stats []model.Stat
	err := r.db.WithContext(ctx).Order("user_id ASC").Find(&stats).Error
	return stats, err
}
