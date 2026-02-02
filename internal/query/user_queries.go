package query

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

// UserQueries contains all user-related SQL queries
type UserQueries struct {
	pool *pgxpool.Pool
}

func NewUserQueries(pool *pgxpool.Pool) *UserQueries {
	return &UserQueries{pool: pool}
}

// CreateUser inserts a new user into the database
func (q *UserQueries) CreateUser(ctx context.Context, username, password, fullName, email, phone, role string, counterID *int, isActive bool) (int, time.Time, time.Time, error) {
	query := `
		INSERT INTO users (username, password, full_name, email, phone, role, counter_id, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`

	var id int
	var createdAt, updatedAt time.Time
	err := q.pool.QueryRow(ctx, query, username, password, fullName, email, phone, role, counterID, isActive).
		Scan(&id, &createdAt, &updatedAt)
	return id, createdAt, updatedAt, err
}

// GetUserByUsername retrieves a user by username
func (q *UserQueries) GetUserByUsername(ctx context.Context, username string) pgx.Row {
	query := `
		SELECT id, username, password, full_name, email, phone, role, is_active, counter_id, created_at, updated_at, last_login
		FROM users WHERE username = $1`

	return q.pool.QueryRow(ctx, query, username)
}

// GetUserByID retrieves a user by ID
func (q *UserQueries) GetUserByID(ctx context.Context, id int) pgx.Row {
	query := `
		SELECT id, username, full_name, email, phone, role, is_active, counter_id, created_at, updated_at, last_login
		FROM users WHERE id = $1`

	return q.pool.QueryRow(ctx, query, id)
}

// UpdateUser updates user information
func (q *UserQueries) UpdateUser(ctx context.Context, fullName, email, phone, role string, counterID *int, isActive bool, id int) error {
	query := `
		UPDATE users SET full_name = $1, email = $2, phone = $3, role = $4, 
		counter_id = $5, is_active = $6, updated_at = NOW()
		WHERE id = $7`

	_, err := q.pool.Exec(ctx, query, fullName, email, phone, role, counterID, isActive, id)
	return err
}

// UpdateUserPassword updates user password
func (q *UserQueries) UpdateUserPassword(ctx context.Context, id int, password string) error {
	query := `UPDATE users SET password = $1, updated_at = NOW() WHERE id = $2`
	_, err := q.pool.Exec(ctx, query, password, id)
	return err
}

// UpdateLastLogin updates the last login timestamp for a user
func (q *UserQueries) UpdateLastLogin(ctx context.Context, id int) error {
	query := `UPDATE users SET last_login = NOW() WHERE id = $1`
	_, err := q.pool.Exec(ctx, query, id)
	return err
}

// DeleteUser deletes a user by ID
func (q *UserQueries) DeleteUser(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := q.pool.Exec(ctx, query, id)
	return err
}

// ListUsers retrieves users with optional role filter
func (q *UserQueries) ListUsers(ctx context.Context, role string) (pgx.Rows, error) {
	var query string
	var args []interface{}

	if role != "" {
		query = `SELECT id, username, full_name, email, phone, role, is_active, counter_id, created_at, updated_at, last_login
				 FROM users WHERE role = $1 ORDER BY created_at DESC`
		args = append(args, role)
	} else {
		query = `SELECT id, username, full_name, email, phone, role, is_active, counter_id, created_at, updated_at, last_login
				 FROM users ORDER BY created_at DESC`
	}

	res, err := q.pool.Query(ctx, query, args...)
	if err != nil {
		log.Error().Err(err).Str("layer", "query").Msg("Failed to list users")
		return nil, err
	}

	return res, nil
}
