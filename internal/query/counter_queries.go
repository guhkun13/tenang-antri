package query

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CounterQueries contains all counter-related SQL queries
type CounterQueries struct {
	pool *pgxpool.Pool
}

func NewCounterQueries(pool *pgxpool.Pool) *CounterQueries {
	return &CounterQueries{pool: pool}
}

// CreateCounter inserts a new counter
func (q *CounterQueries) CreateCounter(ctx context.Context, number, name, location, status string, isActive bool) (int, time.Time, time.Time, error) {
	query := `
		INSERT INTO counters (number, name, location, status, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`

	var id int
	var createdAt, updatedAt time.Time
	err := q.pool.QueryRow(ctx, query, number, name, location, status, isActive).
		Scan(&id, &createdAt, &updatedAt)
	return id, createdAt, updatedAt, err
}

// GetCounterByID retrieves a counter by ID
func (q *CounterQueries) GetCounterByID(ctx context.Context, id int) pgx.Row {
	query := `SELECT id, number, name, location, status, is_active, current_staff_id, created_at, updated_at
			  FROM counters WHERE id = $1`
	return q.pool.QueryRow(ctx, query, id)
}

// UpdateCounter updates counter information
func (q *CounterQueries) UpdateCounter(ctx context.Context, number, name, location, status string, isActive bool, id int) error {
	query := `
		UPDATE counters SET number = $1, name = $2, location = $3, 
		status = $4, is_active = $5, updated_at = NOW()
		WHERE id = $6`

	_, err := q.pool.Exec(ctx, query, number, name, location, status, isActive, id)
	return err
}

// UpdateCounterStatus updates only the status of a counter
func (q *CounterQueries) UpdateCounterStatus(ctx context.Context, status string, id int) error {
	query := `UPDATE counters SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := q.pool.Exec(ctx, query, status, id)
	return err
}

// UpdateCounterStaff updates the current staff assigned to a counter
func (q *CounterQueries) UpdateCounterStaff(ctx context.Context, counterID int, staffID *int) error {
	query := `UPDATE counters SET current_staff_id = $1, updated_at = NOW() WHERE id = $2`
	_, err := q.pool.Exec(ctx, query, staffID, counterID)
	return err
}

// DeleteCounter deletes a counter by ID
func (q *CounterQueries) DeleteCounter(ctx context.Context, id int) error {
	query := `DELETE FROM counters WHERE id = $1`
	_, err := q.pool.Exec(ctx, query, id)
	return err
}

// ListCounters retrieves all counters
func (q *CounterQueries) ListCounters(ctx context.Context) (pgx.Rows, error) {
	query := `SELECT id, number, name, location, status, is_active, current_staff_id, created_at, updated_at
			  FROM counters ORDER BY number`
	return q.pool.Query(ctx, query)
}

// GetCounterCategories retrieves categories for a specific counter
func (q *CounterQueries) GetCounterCategories(ctx context.Context, counterID int) (pgx.Rows, error) {
	query := `
		SELECT c.id, c.name, c.prefix, c.priority, c.color_code, c.description, c.icon, c.is_active, c.created_at, c.updated_at
		FROM categories c
		JOIN counter_categories cc ON c.id = cc.category_id
		WHERE cc.counter_id = $1 AND c.is_active = true
		ORDER BY c.priority DESC, c.name`
	return q.pool.Query(ctx, query, counterID)
}

// AssignCategoryToCounter assigns a category to a counter
func (q *CounterQueries) AssignCategoryToCounter(ctx context.Context, counterID, categoryID int) error {
	query := `INSERT INTO counter_categories (counter_id, category_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := q.pool.Exec(ctx, query, counterID, categoryID)
	return err
}

// RemoveCategoryFromCounter removes a category from a counter
func (q *CounterQueries) RemoveCategoryFromCounter(ctx context.Context, counterID, categoryID int) error {
	query := `DELETE FROM counter_categories WHERE counter_id = $1 AND category_id = $2`
	_, err := q.pool.Exec(ctx, query, counterID, categoryID)
	return err
}

// ClearCounterCategories removes all categories from a counter
func (q *CounterQueries) ClearCounterCategories(ctx context.Context, counterID int) error {
	query := `DELETE FROM counter_categories WHERE counter_id = $1`
	_, err := q.pool.Exec(ctx, query, counterID)
	return err
}
