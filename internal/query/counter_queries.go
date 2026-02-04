package query

import (
	"context"
)

type CounterQueries struct{}

func NewCounterQueries() *CounterQueries {
	return &CounterQueries{}
}

func (q *CounterQueries) CreateCounter(ctx context.Context) string {
	return `INSERT INTO counters (number, name, location, status, is_active, category_id) 
	VALUES ($1, $2, $3, $4, $5, $6) 
	RETURNING id, created_at, updated_at`
}

func (q *CounterQueries) GetCounterByID(ctx context.Context) string {
	return `SELECT id, number, name, location, status, is_active, category_id, current_staff_id, created_at, updated_at FROM counters WHERE id = $1`
}

func (q *CounterQueries) UpdateCounter(ctx context.Context) string {
	return `UPDATE counters SET number = $1, name = $2, location = $3, status = $4, is_active = $5, category_id = $6, updated_at = NOW() WHERE id = $7`
}

func (q *CounterQueries) UpdateCounterStatus(ctx context.Context) string {
	return `UPDATE counters SET status = $1, updated_at = NOW() WHERE id = $2`
}

func (q *CounterQueries) UpdateCounterStaff(ctx context.Context) string {
	return `UPDATE counters SET current_staff_id = $1, updated_at = NOW() WHERE id = $2`
}

func (q *CounterQueries) DeleteCounter(ctx context.Context) string {
	return `DELETE FROM counters WHERE id = $1`
}

func (q *CounterQueries) ListCounters(ctx context.Context) string {
	return `SELECT id, number, name, location, status, is_active, current_staff_id, category_id, created_at, updated_at 
	FROM counters ORDER BY id asc`
}
