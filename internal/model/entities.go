package model

import (
	"database/sql"
	"time"
)

// User represents a system user (admin or staff)
type User struct {
	ID        int            `json:"id" db:"id"`
	Username  string         `json:"username" db:"username"`
	Password  string         `json:"-" db:"-"`
	FullName  sql.NullString `json:"full_name" db:"full_name"`
	Email     sql.NullString `json:"email" db:"email"`
	Phone     sql.NullString `json:"phone" db:"phone"`
	Role      string         `json:"role" db:"role"`
	IsActive  bool           `json:"is_active" db:"is_active"`
	CounterID sql.NullInt64  `json:"counter_id,omitempty" db:"counter_id"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at"`
	LastLogin sql.NullTime   `json:"last_login,omitempty" db:"last_login"`
}

// Category represents a service category
type Category struct {
	ID          int            `json:"id" db:"id"`
	Name        string         `json:"name" db:"name"`
	Prefix      string         `json:"prefix" db:"prefix"`
	Priority    int            `json:"priority" db:"priority"`
	ColorCode   string         `json:"color_code" db:"color_code"`
	Description sql.NullString `json:"description" db:"description"`
	Icon        sql.NullString `json:"icon" db:"icon"`
	IsActive    bool           `json:"is_active" db:"is_active"`
	CreatedAt   time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" db:"updated_at"`
}

// CounterStatus constants
const (
	CounterStatusOffline = "offline"
	CounterStatusIdle    = "idle"
	CounterStatusServing = "serving"
	CounterStatusPaused  = "paused"
)

// Counter represents a service counter
type Counter struct {
	ID             int            `json:"id" db:"id"`
	Number         string         `json:"number" db:"number"`
	Name           sql.NullString `json:"name" db:"name"`
	Location       sql.NullString `json:"location" db:"location"`
	Status         string         `json:"status" db:"status"`
	CategoryID     sql.NullInt64  `json:"category_id,omitempty" db:"category_id"`
	Category       *Category      `json:"category,omitempty" db:"-"`
	CreatedAt      time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at" db:"updated_at"`
	CurrentStaffID sql.NullInt64  `json:"-" db:"current_staff_id"`
}

// Ticket represents a queue ticket
type Ticket struct {
	ID            int            `json:"id" db:"id"`
	TicketNumber  string         `json:"ticket_number" db:"ticket_number"`
	CategoryID    sql.NullInt64  `json:"category_id,omitempty" db:"category_id"`
	Category      *Category      `json:"category,omitempty" db:"-"`
	CounterID     sql.NullInt64  `json:"counter_id,omitempty" db:"counter_id"`
	Counter       *Counter       `json:"counter,omitempty" db:"-"`
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
