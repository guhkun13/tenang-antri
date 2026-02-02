package query

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// StatsQueries contains all statistics-related SQL queries
type StatsQueries struct {
	pool *pgxpool.Pool
}

func NewStatsQueries(pool *pgxpool.Pool) *StatsQueries {
	return &StatsQueries{pool: pool}
}

// GetDashboardStats retrieves comprehensive dashboard statistics
func (q *StatsQueries) GetDashboardStats(ctx context.Context) (pgx.Rows, error) {
	// This will be implemented in the service layer to combine multiple queries
	return nil, fmt.Errorf("complex query - implement in service layer")
}

// GetQueueLengthByCategory retrieves queue length for each category
func (q *StatsQueries) GetQueueLengthByCategory(ctx context.Context) (pgx.Rows, error) {
	query := `
		SELECT c.id, c.name, c.prefix, c.color_code, COUNT(t.id) as waiting_count
		FROM categories c
		LEFT JOIN tickets t ON c.id = t.category_id AND t.status = 'waiting'
		WHERE c.is_active = true
		GROUP BY c.id, c.name, c.prefix, c.color_code
		ORDER BY waiting_count DESC, c.priority DESC`
	return q.pool.Query(ctx, query)
}

// GetQueueLengthByCategories retrieves queue length for specific categories
func (q *StatsQueries) GetQueueLengthByCategories(ctx context.Context, categoryIDs []int) (pgx.Rows, error) {
	if len(categoryIDs) == 0 {
		return nil, fmt.Errorf("no categories provided")
	}

	placeholders := make([]string, len(categoryIDs))
	args := make([]interface{}, len(categoryIDs))
	for i, id := range categoryIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT c.id, c.name, c.prefix, c.color_code, COUNT(t.id) as waiting_count
		FROM categories c
		LEFT JOIN tickets t ON c.id = t.category_id AND t.status = 'waiting'
		WHERE c.is_active = true AND c.id IN (%s)
		GROUP BY c.id, c.name, c.prefix, c.color_code
		ORDER BY waiting_count DESC, c.priority DESC`, strings.Join(placeholders, ","))

	return q.pool.Query(ctx, query, args...)
}

// GetHourlyDistribution retrieves hourly ticket distribution
func (q *StatsQueries) GetHourlyDistribution(ctx context.Context) (pgx.Rows, error) {
	query := `
		SELECT EXTRACT(HOUR FROM created_at)::INT as hour, COUNT(*) as count
		FROM tickets
		WHERE DATE(created_at) = CURRENT_DATE
		GROUP BY EXTRACT(HOUR FROM created_at)
		ORDER BY hour`
	return q.pool.Query(ctx, query)
}

// GetCurrentlyServingTickets retrieves currently serving tickets for display
func (q *StatsQueries) GetCurrentlyServingTickets(ctx context.Context) (pgx.Rows, error) {
	query := `
		SELECT t.ticket_number, c.number, cat.prefix, cat.color_code, t.status
		FROM tickets t
		JOIN counters c ON t.counter_id = c.id
		JOIN categories cat ON t.category_id = cat.id
		WHERE t.status = 'serving'
		ORDER BY t.called_at DESC
		LIMIT 10`
	return q.pool.Query(ctx, query)
}

// Individual stat queries for dashboard stats

// GetTotalTicketsToday retrieves total tickets created today
func (q *StatsQueries) GetTotalTicketsToday(ctx context.Context) pgx.Row {
	query := `SELECT COUNT(*) FROM tickets WHERE DATE(created_at) = CURRENT_DATE`
	return q.pool.QueryRow(ctx, query)
}

// GetCurrentlyServingCount retrieves count of currently serving tickets
func (q *StatsQueries) GetCurrentlyServingCount(ctx context.Context) pgx.Row {
	query := `SELECT COUNT(*) FROM tickets WHERE status = 'serving'`
	return q.pool.QueryRow(ctx, query)
}

// GetWaitingTicketsCount retrieves count of waiting tickets
func (q *StatsQueries) GetWaitingTicketsCount(ctx context.Context) pgx.Row {
	query := `SELECT COUNT(*) FROM tickets WHERE status = 'waiting'`
	return q.pool.QueryRow(ctx, query)
}

// GetActiveCountersCount retrieves count of active counters
func (q *StatsQueries) GetActiveCountersCount(ctx context.Context) pgx.Row {
	query := `SELECT COUNT(*) FROM counters WHERE status = 'active' AND is_active = true`
	return q.pool.QueryRow(ctx, query)
}

// GetPausedCountersCount retrieves count of paused counters
func (q *StatsQueries) GetPausedCountersCount(ctx context.Context) pgx.Row {
	query := `SELECT COUNT(*) FROM counters WHERE status = 'paused' AND is_active = true`
	return q.pool.QueryRow(ctx, query)
}

// GetAvgWaitTimeToday retrieves average wait time for today
func (q *StatsQueries) GetAvgWaitTimeToday(ctx context.Context) pgx.Row {
	query := `SELECT COALESCE(AVG(wait_time)::INT, 0) FROM tickets 
			 WHERE DATE(created_at) = CURRENT_DATE AND wait_time IS NOT NULL`
	return q.pool.QueryRow(ctx, query)
}

// GetAvgServiceTimeToday retrieves average service time for today
func (q *StatsQueries) GetAvgServiceTimeToday(ctx context.Context) pgx.Row {
	query := `SELECT COALESCE(AVG(service_time)::INT, 0) FROM tickets 
			 WHERE DATE(created_at) = CURRENT_DATE AND service_time IS NOT NULL`
	return q.pool.QueryRow(ctx, query)
}

// GetTicketsByStatusToday retrieves tickets grouped by status for today
func (q *StatsQueries) GetTicketsByStatusToday(ctx context.Context) (pgx.Rows, error) {
	query := `SELECT status, COUNT(*) FROM tickets WHERE DATE(created_at) = CURRENT_DATE GROUP BY status`
	return q.pool.Query(ctx, query)
}
