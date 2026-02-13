package query

import (
	"context"
)

type CategoryQueries struct{}

func NewCategoryQueries() *CategoryQueries {
	return &CategoryQueries{}
}

func (q *CategoryQueries) CreateCategory(ctx context.Context) string {
	return `INSERT INTO categories (name, prefix, priority, color_code, description, icon, is_active) 
	VALUES ($1, $2, $3, $4, $5, $6, $7) 
	RETURNING id, created_at, updated_at`
}

func (q *CategoryQueries) GetCategoryByID(ctx context.Context) string {
	return `SELECT id, name, prefix, priority, color_code, description, icon, is_active, created_at, updated_at 
	FROM categories WHERE id = $1`
}

func (q *CategoryQueries) UpdateCategory(ctx context.Context) string {
	return `UPDATE categories 
	SET name = $1, prefix = $2, priority = $3, color_code = $4, description = $5, icon = $6, is_active = $7, updated_at = NOW() 
	WHERE id = $8`
}

func (q *CategoryQueries) DeleteCategory(ctx context.Context) string {
	return `DELETE FROM categories WHERE id = $1`
}

func (q *CategoryQueries) ListCategories(ctx context.Context, activeOnly bool, withCountersOnly bool) string {
	query := `SELECT DISTINCT categories.id, categories.name, categories.prefix, categories.priority, categories.color_code, categories.description, categories.icon, categories.is_active, categories.created_at, categories.updated_at FROM categories`

	if withCountersOnly {
		query += ` INNER JOIN counters ON counters.category_id = categories.id AND counters.current_staff_id IS NOT NULL`
	}

	if activeOnly {
		query += ` WHERE categories.is_active = true`
	}

	query += ` ORDER BY categories.priority DESC, categories.name`
	return query
}
