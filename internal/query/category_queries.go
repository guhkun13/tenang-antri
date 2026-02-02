package query

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CategoryQueries contains all category-related SQL queries
type CategoryQueries struct {
	pool *pgxpool.Pool
}

func NewCategoryQueries(pool *pgxpool.Pool) *CategoryQueries {
	return &CategoryQueries{pool: pool}
}

// CreateCategory inserts a new category
func (q *CategoryQueries) CreateCategory(ctx context.Context, name, prefix string, priority int, colorCode, description, icon string, isActive bool) (int, time.Time, time.Time, error) {
	query := `
		INSERT INTO categories (name, prefix, priority, color_code, description, icon, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`

	var id int
	var createdAt, updatedAt time.Time
	err := q.pool.QueryRow(ctx, query, name, prefix, priority, colorCode, description, icon, isActive).
		Scan(&id, &createdAt, &updatedAt)
	return id, createdAt, updatedAt, err
}

// GetCategoryByID retrieves a category by ID
func (q *CategoryQueries) GetCategoryByID(ctx context.Context, id int) pgx.Row {
	query := `SELECT id, name, prefix, priority, color_code, description, icon, is_active, created_at, updated_at
			  FROM categories WHERE id = $1`
	return q.pool.QueryRow(ctx, query, id)
}

// UpdateCategory updates category information
func (q *CategoryQueries) UpdateCategory(ctx context.Context, name, prefix string, priority int, colorCode, description, icon string, isActive bool, id int) error {
	query := `
		UPDATE categories SET name = $1, prefix = $2, priority = $3, color_code = $4,
		description = $5, icon = $6, is_active = $7, updated_at = NOW()
		WHERE id = $8`

	_, err := q.pool.Exec(ctx, query, name, prefix, priority, colorCode, description, icon, isActive, id)
	return err
}

// DeleteCategory deletes a category by ID
func (q *CategoryQueries) DeleteCategory(ctx context.Context, id int) error {
	query := `DELETE FROM categories WHERE id = $1`
	_, err := q.pool.Exec(ctx, query, id)
	return err
}

// ListCategories retrieves categories with optional active filter
func (q *CategoryQueries) ListCategories(ctx context.Context, activeOnly bool) (pgx.Rows, error) {
	query := `SELECT id, name, prefix, priority, color_code, description, icon, is_active, created_at, updated_at
			  FROM categories`
	if activeOnly {
		query += ` WHERE is_active = true`
	}
	query += ` ORDER BY priority DESC, name`

	return q.pool.Query(ctx, query)
}
