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
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at"`
	LastLogin sql.NullTime   `json:"last_login,omitempty" db:"last_login"`
}
