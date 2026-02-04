package query

import (
	"context"
)

type UserQueries struct{}

func NewUserQueries() *UserQueries {
	return &UserQueries{}
}

func (q *UserQueries) CreateUser(ctx context.Context) string {
	return `INSERT INTO users (username, password, full_name, email, phone, role, counter_id, is_active) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, created_at, updated_at`
}

func (q *UserQueries) GetUserByUsername(ctx context.Context) string {
	return `SELECT id, username, password, full_name, email, phone, role, is_active, counter_id, created_at, updated_at, last_login FROM users WHERE username = $1`
}

func (q *UserQueries) GetUserByID(ctx context.Context) string {
	return `SELECT id, username, full_name, email, phone, role, is_active, counter_id, created_at, updated_at, last_login FROM users WHERE id = $1`
}

func (q *UserQueries) UpdateUser(ctx context.Context) string {
	return `UPDATE users SET full_name = $1, email = $2, phone = $3, role = $4, counter_id = $5, is_active = $6, updated_at = NOW() WHERE id = $7`
}

func (q *UserQueries) UpdateUserPassword(ctx context.Context) string {
	return `UPDATE users SET password = $1, updated_at = NOW() WHERE id = $2`
}

func (q *UserQueries) UpdateLastLogin(ctx context.Context) string {
	return `UPDATE users SET last_login = NOW() WHERE id = $1`
}

func (q *UserQueries) DeleteUser(ctx context.Context) string {
	return `DELETE FROM users WHERE id = $1`
}

func (q *UserQueries) ListUsers(ctx context.Context, role string) string {
	if role != "" {
		return `SELECT id, username, full_name, email, phone, role, is_active, counter_id, created_at, updated_at, last_login FROM users WHERE role = $1 ORDER BY created_at DESC`
	}
	return `SELECT id, username, full_name, email, phone, role, is_active, counter_id, created_at, updated_at, last_login FROM users ORDER BY created_at DESC`
}
