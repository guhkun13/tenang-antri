package model

import (
	"database/sql"
	"time"
)

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
