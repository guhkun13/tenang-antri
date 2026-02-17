package model

import (
	"time"
)

// UserCounter maps users to their assigned counters
type UserCounter struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	CounterID int       `json:"counter_id" db:"counter_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
