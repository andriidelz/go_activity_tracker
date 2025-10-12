package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type JSONB map[string]any

func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONB) Scan(value any) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("Scan: expected []byte, got %T", value)
	}
	return json.Unmarshal(bytes, j)
}

type Event struct {
	ID        int       `gorm:"primaryKey"`
	UserID    int       `json:"user_id"`
	Type      string    `json:"type"`
	Metadata  JSONB     `json:"metadata" gorm:"type:jsonb"`
	CreatedAt time.Time `json:"created_at"`
}

type Stat struct {
	ID         int `gorm:"primaryKey"`
	UserID     int `json:"user_id"`
	EventCount int `json:"event_count"`
}
