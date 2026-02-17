package model

import (
	"database/sql"
	"time"
)

// CounterStatus constants
const (
	CounterStatusOffline = "offline"
	CounterStatusIdle    = "idle"
	CounterStatusServing = "serving"
	CounterStatusPaused  = "paused"
)

// Counter represents a service counter
type Counter struct {
	ID        int            `json:"id" db:"id"`
	Number    string         `json:"number" db:"number"`
	Name      sql.NullString `json:"name" db:"name"`
	Location  sql.NullString `json:"location" db:"location"`
	Status    string         `json:"status" db:"status"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at"`
}
