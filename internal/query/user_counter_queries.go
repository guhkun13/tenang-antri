package query

import (
	"context"
)

type UserCounterQueries struct{}

func NewUserCounterQueries() *UserCounterQueries {
	return &UserCounterQueries{}
}

func (q *UserCounterQueries) Create(ctx context.Context) string {
	return `INSERT INTO user_counters (user_id, counter_id) VALUES ($1, $2) RETURNING id, created_at, updated_at`
}

func (q *UserCounterQueries) GetByUserID(ctx context.Context) string {
	return `SELECT id, user_id, counter_id, created_at, updated_at FROM user_counters WHERE user_id = $1`
}

func (q *UserCounterQueries) GetByCounterID(ctx context.Context) string {
	return `SELECT id, user_id, counter_id, created_at, updated_at FROM user_counters WHERE counter_id = $1`
}

func (q *UserCounterQueries) GetCounterIDByUserID(ctx context.Context) string {
	return `SELECT counter_id FROM user_counters WHERE user_id = $1 LIMIT 1`
}

func (q *UserCounterQueries) DeleteByUserID(ctx context.Context) string {
	return `DELETE FROM user_counters WHERE user_id = $1`
}

func (q *UserCounterQueries) DeleteByCounterID(ctx context.Context) string {
	return `DELETE FROM user_counters WHERE counter_id = $1`
}

func (q *UserCounterQueries) ListAll(ctx context.Context) string {
	return `SELECT id, user_id, counter_id, created_at, updated_at FROM user_counters ORDER BY created_at DESC`
}
