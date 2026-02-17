package query

import (
	"context"
)

type CounterCategoryQueries struct{}

func NewCounterCategoryQueries() *CounterCategoryQueries {
	return &CounterCategoryQueries{}
}

func (q *CounterCategoryQueries) Create(ctx context.Context) string {
	return `INSERT INTO counter_category (counter_id, category_id) VALUES ($1, $2) RETURNING id, created_at, updated_at`
}

func (q *CounterCategoryQueries) GetByID(ctx context.Context) string {
	return `SELECT id, counter_id, category_id, created_at, updated_at FROM counter_category WHERE id = $1`
}

func (q *CounterCategoryQueries) GetByCounterID(ctx context.Context) string {
	return `SELECT id, counter_id, category_id, created_at, updated_at FROM counter_category WHERE counter_id = $1`
}

func (q *CounterCategoryQueries) GetCategoryIDsByCounterID(ctx context.Context) string {
	return `SELECT category_id FROM counter_category WHERE counter_id = $1`
}

func (q *CounterCategoryQueries) GetByCategoryID(ctx context.Context) string {
	return `SELECT id, counter_id, category_id, created_at, updated_at FROM counter_category WHERE category_id = $1`
}

func (q *CounterCategoryQueries) GetCounterIDsByCategoryID(ctx context.Context) string {
	return `SELECT counter_id FROM counter_category WHERE category_id = $1`
}

func (q *CounterCategoryQueries) DeleteByID(ctx context.Context) string {
	return `DELETE FROM counter_category WHERE id = $1`
}

func (q *CounterCategoryQueries) DeleteByCounterID(ctx context.Context) string {
	return `DELETE FROM counter_category WHERE counter_id = $1`
}

func (q *CounterCategoryQueries) DeleteByCategoryID(ctx context.Context) string {
	return `DELETE FROM counter_category WHERE category_id = $1`
}

func (q *CounterCategoryQueries) DeleteByCounterAndCategory(ctx context.Context) string {
	return `DELETE FROM counter_category WHERE counter_id = $1 AND category_id = $2`
}

func (q *CounterCategoryQueries) ListAll(ctx context.Context) string {
	return `SELECT id, counter_id, category_id, created_at, updated_at FROM counter_category ORDER BY created_at DESC`
}
