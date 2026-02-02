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
	Role      string         `json:"role" db:"role"` // admin, staff
	IsActive  bool           `json:"is_active" db:"is_active"`
	CounterID sql.NullInt64  `json:"counter_id,omitempty" db:"counter_id"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at"`
	LastLogin *time.Time     `json:"last_login,omitempty" db:"last_login"`
}

// Category represents a service category
type Category struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Prefix      string    `json:"prefix" db:"prefix"`
	Priority    int       `json:"priority" db:"priority"`
	ColorCode   string    `json:"color_code" db:"color_code"`
	Description string    `json:"description" db:"description"`
	Icon        string    `json:"icon" db:"icon"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Counter represents a service counter
type Counter struct {
	ID             int        `json:"id" db:"id"`
	Number         string     `json:"number" db:"number"`
	Name           string     `json:"name" db:"name"`
	Location       string     `json:"location" db:"location"`
	Status         string     `json:"status" db:"status"` // active, paused, inactive
	IsActive       bool       `json:"is_active" db:"is_active"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
	Categories     []Category `json:"categories,omitempty" db:"-"`
	CurrentStaffID *int       `json:"-" db:"current_staff_id"`
	CurrentStaff   *User      `json:"-" db:"-"`
}

// CounterCategory links counters to categories
type CounterCategory struct {
	CounterID  int `json:"counter_id" db:"counter_id"`
	CategoryID int `json:"category_id" db:"category_id"`
}

// Ticket represents a queue ticket
type Ticket struct {
	ID           int        `json:"id" db:"id"`
	TicketNumber string     `json:"ticket_number" db:"ticket_number"`
	CategoryID   int        `json:"category_id" db:"category_id"`
	Category     *Category  `json:"category,omitempty" db:"-"`
	CounterID    *int       `json:"counter_id,omitempty" db:"counter_id"`
	Counter      *Counter   `json:"counter,omitempty" db:"-"`
	Status       string     `json:"status" db:"status"` // waiting, serving, completed, no_show, cancelled
	Priority     int        `json:"priority" db:"priority"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	CalledAt     *time.Time `json:"called_at,omitempty" db:"called_at"`
	CompletedAt  *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	WaitTime     *int       `json:"wait_time,omitempty" db:"wait_time"`       // in seconds
	ServiceTime  *int       `json:"service_time,omitempty" db:"service_time"` // in seconds
	Notes        string     `json:"notes" db:"notes"`
}
