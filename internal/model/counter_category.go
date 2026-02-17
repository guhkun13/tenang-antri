package model

import (
	"time"
)

// CounterCategory maps counters to the categories they serve (one counter can serve multiple categories)
type CounterCategory struct {
	ID         int       `json:"id" db:"id"`
	CounterID  int       `json:"counter_id" db:"counter_id"`
	CategoryID int       `json:"category_id" db:"category_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}
