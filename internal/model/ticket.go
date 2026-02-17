package model

import (
	"database/sql"
	"time"
)

// Ticket represents a queue ticket
type Ticket struct {
	ID            int            `json:"id" db:"id"`
	TicketNumber  string         `json:"ticket_number" db:"ticket_number"`
	CategoryID    sql.NullInt64  `json:"category_id,omitempty" db:"category_id"`
	CounterID     sql.NullInt64  `json:"counter_id,omitempty" db:"counter_id"`
	Status        string         `json:"status" db:"status"`
	Priority      int            `json:"priority" db:"priority"`
	CreatedAt     time.Time      `json:"created_at" db:"created_at"`
	CalledAt      sql.NullTime   `json:"called_at,omitempty" db:"called_at"`
	CompletedAt   sql.NullTime   `json:"completed_at,omitempty" db:"completed_at"`
	WaitTime      sql.NullInt64  `json:"wait_time,omitempty" db:"wait_time"`
	ServiceTime   sql.NullInt64  `json:"service_time,omitempty" db:"service_time"`
	DailySequence int            `json:"daily_sequence" db:"daily_sequence"`
	QueueDate     time.Time      `json:"queue_date" db:"queue_date"`
	Notes         sql.NullString `json:"notes" db:"notes"`
}
