package query

import (
	"context"
	"fmt"
	"strings"
)

type StatsQueries struct{}

func NewStatsQueries() *StatsQueries {
	return &StatsQueries{}
}

func (q *StatsQueries) GetDashboardStats(ctx context.Context) string {
	return ""
}

func (q *StatsQueries) GetQueueLengthByCategory(ctx context.Context) string {
	return `SELECT c.id, c.name, c.prefix, c.color_code, COUNT(t.id) as waiting_count FROM categories c LEFT JOIN tickets t ON c.id = t.category_id AND t.status = 'waiting' WHERE c.is_active = true GROUP BY c.id, c.name, c.prefix, c.color_code ORDER BY waiting_count DESC, c.priority DESC`
}

func (q *StatsQueries) GetQueueLengthByCategories(ctx context.Context, categoryIDs []int) string {
	placeholders := make([]string, len(categoryIDs))
	for i := range categoryIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}
	return fmt.Sprintf(`SELECT c.id, c.name, c.prefix, c.color_code, COUNT(t.id) as waiting_count FROM categories c LEFT JOIN tickets t ON c.id = t.category_id AND t.status = 'waiting' WHERE c.is_active = true AND c.id IN (%s) GROUP BY c.id, c.name, c.prefix, c.color_code ORDER BY waiting_count DESC, c.priority DESC`, strings.Join(placeholders, ","))
}

func (q *StatsQueries) GetHourlyDistribution(ctx context.Context) string {
	return `SELECT EXTRACT(HOUR FROM created_at)::INT as hour, COUNT(*) as count FROM tickets WHERE DATE(created_at) = CURRENT_DATE GROUP BY EXTRACT(HOUR FROM created_at) ORDER BY hour`
}

func (q *StatsQueries) GetCurrentlyServingTickets(ctx context.Context) string {
	return `SELECT t.ticket_number, c.number, cat.prefix, cat.color_code, t.status FROM tickets t JOIN counters c ON t.counter_id = c.id JOIN categories cat ON t.category_id = cat.id WHERE t.status = 'serving' ORDER BY t.called_at DESC LIMIT 10`
}

func (q *StatsQueries) GetTotalTicketsToday(ctx context.Context) string {
	return `SELECT COUNT(*) FROM tickets WHERE DATE(created_at) = CURRENT_DATE`
}

func (q *StatsQueries) GetCurrentlyServingCount(ctx context.Context) string {
	return `SELECT COUNT(*) FROM tickets WHERE status = 'serving'`
}

func (q *StatsQueries) GetWaitingTicketsCount(ctx context.Context) string {
	return `SELECT COUNT(*) FROM tickets WHERE status = 'waiting'`
}

func (q *StatsQueries) GetActiveCountersCount(ctx context.Context) string {
	return `SELECT COUNT(*) FROM counters WHERE is_active = true`
}

func (q *StatsQueries) GetPausedCountersCount(ctx context.Context) string {
	return `SELECT COUNT(*) FROM counters WHERE status = 'paused' AND is_active = true`
}

func (q *StatsQueries) GetAvgWaitTimeToday(ctx context.Context) string {
	return `SELECT COALESCE(AVG(wait_time)::INT, 0) FROM tickets WHERE DATE(created_at) = CURRENT_DATE AND wait_time IS NOT NULL`
}

func (q *StatsQueries) GetAvgServiceTimeToday(ctx context.Context) string {
	return `SELECT COALESCE(AVG(service_time)::INT, 0) FROM tickets WHERE DATE(created_at) = CURRENT_DATE AND service_time IS NOT NULL`
}

func (q *StatsQueries) GetTicketsByStatusToday(ctx context.Context) string {
	return `SELECT status, COUNT(*) FROM tickets WHERE DATE(created_at) = CURRENT_DATE GROUP BY status`
}
