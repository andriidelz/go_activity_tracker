package repository

import (
	"github.com/andriidelzz/go-activity-tracker/internal/model"

	"gorm.io/gorm"
)

type RepositoryInterface interface {
	CreateEvent(event *model.Event) error
	GetEvents(userID int) ([]model.Event, error)
	AggregateLastPeriod() error
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

type Repository struct {
	db *gorm.DB
}

func (r *Repository) CreateEvent(event *model.Event) error {
	return r.db.Create(event).Error
}

func (r *Repository) GetEvents(userID int) ([]model.Event, error) {
	var events []model.Event
	return events, r.db.Where("user_id = ?", userID).Find(&events).Error
}

func (r *Repository) AggregateLastPeriod() error {
	var results []struct {
		UserID     int
		EventCount int // total count of created events by user.
	}
	err := r.db.Model(&model.Event{}).
		Select("user_id, COUNT(*) as event_count").
		Group("user_id").
		Scan(&results).Error
	if err != nil {
		return err
	}

	for _, res := range results {
		stat := model.Stat{
			UserID:     res.UserID,
			EventCount: res.EventCount,
		}
		r.db.Create(&stat)
	}
	return nil
}

func (r *Repository) GetStats() ([]model.Stat, error) {
	var stats []model.Stat
	err := r.db.Find(&stats).Error
	return stats, err
}
