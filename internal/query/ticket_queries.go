package query

import (
	"context"
	"fmt"
	"strings"
)

type TicketQueries struct{}

func NewTicketQueries() *TicketQueries {
	return &TicketQueries{}
}

func (q *TicketQueries) CreateTicket(ctx context.Context) string {
	return `INSERT INTO tickets (ticket_number, category_id, status, priority, notes, daily_sequence, queue_date) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at`
}

func (q *TicketQueries) GetTicketByID(ctx context.Context) string {
	return `SELECT t.id, t.ticket_number, t.category_id, t.counter_id, t.status, t.priority, t.created_at, t.called_at, t.completed_at, t.wait_time, t.service_time, t.daily_sequence, t.queue_date, t.notes FROM tickets t WHERE t.id = $1`
}

func (q *TicketQueries) GetTicketWithDetails(ctx context.Context) string {
	return `SELECT t.id, t.ticket_number, t.category_id, t.counter_id, t.status, t.priority, t.created_at, t.called_at, t.completed_at, t.wait_time, t.service_time, t.daily_sequence, t.queue_date, t.notes, c.id, c.name, c.prefix, c.color_code, co.id, co.number, co.name FROM tickets t LEFT JOIN categories c ON t.category_id = c.id LEFT JOIN counters co ON t.counter_id = co.id WHERE t.id = $1`
}

func (q *TicketQueries) GetTicketByNumber(ctx context.Context) string {
	return `SELECT t.id, t.ticket_number, t.category_id, t.counter_id, t.status, t.priority, t.created_at, t.called_at, t.completed_at, t.wait_time, t.service_time, t.daily_sequence, t.queue_date, t.notes, c.id, c.name, c.prefix, c.color_code, co.id, co.number, co.name FROM tickets t LEFT JOIN categories c ON t.category_id = c.id LEFT JOIN counters co ON t.counter_id = co.id WHERE t.ticket_number = $1`
}

func (q *TicketQueries) UpdateTicketStatus(ctx context.Context, status string) string {
	switch status {
	case "serving":
		return `UPDATE tickets SET status = $1, called_at = NOW() WHERE id = $2`
	case "completed", "no_show":
		return `UPDATE tickets SET status = $1, completed_at = NOW(), wait_time = EXTRACT(EPOCH FROM (called_at - created_at))::INT, service_time = EXTRACT(EPOCH FROM (NOW() - called_at))::INT WHERE id = $2`
	default:
		return `UPDATE tickets SET status = $1 WHERE id = $2`
	}
}

func (q *TicketQueries) AssignTicketToCounter(ctx context.Context) string {
	return `UPDATE tickets SET counter_id = $1, status = 'serving', called_at = NOW() WHERE id = $2`
}

func (q *TicketQueries) GetNextTicket(ctx context.Context, categoryIDs []int) string {
	return `SELECT t.id, t.ticket_number, t.category_id, t.status, t.priority, t.created_at, t.daily_sequence, t.queue_date, t.notes FROM tickets t WHERE t.category_id = ANY($1) AND t.status = 'waiting' ORDER BY (SELECT priority FROM categories WHERE id = t.category_id) DESC, t.created_at ASC LIMIT 1`
}

func (q *TicketQueries) GetCurrentTicketForCounter(ctx context.Context) string {
	return `SELECT t.id, t.ticket_number, t.category_id, t.counter_id, t.status, t.priority, t.created_at, t.called_at, t.completed_at, t.wait_time, t.service_time, t.daily_sequence, t.queue_date, t.notes FROM tickets t WHERE t.counter_id = $1 AND t.status = 'serving' ORDER BY t.called_at DESC LIMIT 1`
}

type ListTicketsResult struct {
	Query string
	Args  []any
}

func (q *TicketQueries) ListTickets(ctx context.Context, filters map[string]interface{}) ListTicketsResult {
	query := `SELECT t.id, t.ticket_number, t.category_id, t.counter_id, t.status, t.priority, t.created_at, t.called_at, t.completed_at, t.wait_time, t.service_time, t.daily_sequence, t.queue_date, t.notes, c.name as category_name, c.prefix as category_prefix, c.color_code as category_color, co.number as counter_number, co.name as counter_name FROM tickets t LEFT JOIN categories c ON t.category_id = c.id LEFT JOIN counters co ON t.counter_id = co.id WHERE 1=1`
	args := make([]any, 0)
	argCount := 1

	if search, ok := filters["search"]; ok && search != "" {
		query += fmt.Sprintf(" AND (t.ticket_number ILIKE $%d OR c.name ILIKE $%d)", argCount, argCount)
		args = append(args, search)
		argCount++
	}
	if status, ok := filters["status"]; ok && status != "" {
		query += fmt.Sprintf(" AND t.status = $%d", argCount)
		args = append(args, status)
		argCount++
	}
	if categoryID, ok := filters["category_id"]; ok && categoryID != 0 {
		query += fmt.Sprintf(" AND t.category_id = $%d", argCount)
		args = append(args, categoryID)
		argCount++
	}
	if counterID, ok := filters["counter_id"]; ok && counterID != 0 {
		query += fmt.Sprintf(" AND t.counter_id = $%d", argCount)
		args = append(args, counterID)
		argCount++
	}
	if dateFrom, ok := filters["date_from"]; ok && dateFrom != "" {
		query += fmt.Sprintf(" AND t.created_at >= $%d", argCount)
		args = append(args, dateFrom)
		argCount++
	}
	if dateTo, ok := filters["date_to"]; ok && dateTo != "" {
		query += fmt.Sprintf(" AND t.created_at < ($%d::date + interval '1 day')", argCount)
		args = append(args, dateTo)
		argCount++
	}

	query += ` ORDER BY t.created_at DESC`

	if limit, ok := filters["limit"].(int); ok && limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
		if offset, ok := filters["offset"].(int); ok && offset > 0 {
			query += fmt.Sprintf(" OFFSET %d", offset)
		}
	} else {
		query += " LIMIT 50"
	}

	return ListTicketsResult{Query: query, Args: args}
}

func (q *TicketQueries) GetTodayTicketCount(ctx context.Context) string {
	return `SELECT COUNT(*) FROM tickets WHERE queue_date = CURRENT_DATE`
}

func (q *TicketQueries) GetTodayTicketCountByCategory(ctx context.Context) string {
	return `SELECT COUNT(*) FROM tickets WHERE category_id = $1 AND queue_date = CURRENT_DATE`
}

func (q *TicketQueries) GenerateTicketNumber(ctx context.Context) string {
	return `SELECT COALESCE(MAX(daily_sequence), 0) + 1 FROM tickets WHERE category_id = $1 AND queue_date = CURRENT_DATE`
}

func (q *TicketQueries) GetWaitingTicketsPreview(ctx context.Context) string {
	return `SELECT t.id, t.ticket_number, t.category_id, t.counter_id, t.status, t.priority, t.created_at, t.called_at, t.completed_at, t.wait_time, t.service_time, t.daily_sequence, t.queue_date, t.notes FROM tickets t WHERE t.status = 'waiting' ORDER BY (SELECT priority FROM categories WHERE id = t.category_id) DESC, t.created_at ASC LIMIT $1`
}

func (q *TicketQueries) GetWaitingTicketsPreviewByCategories(ctx context.Context, categoryIDs []int) string {
	placeholders := make([]string, len(categoryIDs))
	for i := range categoryIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+2)
	}
	return fmt.Sprintf(`SELECT t.id, t.ticket_number, t.category_id, t.counter_id, t.status, 
	t.priority, t.created_at, t.called_at, t.completed_at, t.wait_time, t.service_time, t.daily_sequence, t.queue_date, t.notes 
	FROM tickets t 
	WHERE t.status = 'waiting' AND t.category_id IN (%s) 
	ORDER BY (
			SELECT priority 
			FROM categories 
			WHERE id = t.category_id
		) 
	DESC, t.created_at ASC LIMIT $1`, strings.Join(placeholders, ","))
}

func (q *TicketQueries) GetTodayCompletedTicketsByCategories(ctx context.Context, categoryIDs []int) string {
	placeholders := make([]string, len(categoryIDs))
	for i := range categoryIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}
	return fmt.Sprintf(`SELECT t.id, t.ticket_number, t.category_id, t.counter_id, t.status, t.priority, t.created_at, t.called_at, t.completed_at, t.wait_time, t.service_time, t.daily_sequence, t.queue_date, t.notes, c.id, c.name, c.prefix, c.color_code, co.id, co.number, co.name FROM tickets t LEFT JOIN categories c ON t.category_id = c.id LEFT JOIN counters co ON t.counter_id = co.id WHERE t.status = 'completed' AND t.queue_date = CURRENT_DATE AND t.category_id IN (%s) ORDER BY t.completed_at DESC LIMIT 50`, strings.Join(placeholders, ","))
}

func (q *TicketQueries) GetLastCalledTicketByCategory(ctx context.Context) string {
	return `SELECT ticket_number FROM tickets WHERE category_id = $1 AND status IN ('serving', 'completed') ORDER BY called_at DESC LIMIT 1`
}

func (q *TicketQueries) GetTodayTicketsByCategories(ctx context.Context, categoryIDs []int) string {
	placeholders := make([]string, len(categoryIDs))
	for i := range categoryIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}
	return fmt.Sprintf(`SELECT t.id, t.ticket_number, t.category_id, t.counter_id, t.status, t.priority, t.created_at, t.called_at, t.completed_at, t.wait_time, t.service_time, t.daily_sequence, t.queue_date, t.notes, c.id, c.name, c.prefix, c.color_code, co.id, co.number, co.name FROM tickets t LEFT JOIN categories c ON t.category_id = c.id LEFT JOIN counters co ON t.counter_id = co.id WHERE t.queue_date = CURRENT_DATE AND t.category_id IN (%s) ORDER BY t.created_at DESC`, strings.Join(placeholders, ","))
}

func (q *TicketQueries) GetAllTodayTickets(ctx context.Context) string {
	return `SELECT t.id, t.ticket_number, t.category_id, t.counter_id, t.status, t.priority, t.created_at, t.called_at, t.completed_at, t.wait_time, t.service_time, t.daily_sequence, t.queue_date, t.notes, c.id, c.name, c.prefix, c.color_code, co.id, co.number, co.name FROM tickets t LEFT JOIN categories c ON t.category_id = c.id LEFT JOIN counters co ON t.counter_id = co.id WHERE t.queue_date = CURRENT_DATE ORDER BY t.created_at DESC`
}

func (q *TicketQueries) GetAllTicketsByCategories(ctx context.Context, categoryIDs []int) string {
	placeholders := make([]string, len(categoryIDs))
	for i := range categoryIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}
	return fmt.Sprintf(`SELECT t.id, t.ticket_number, t.category_id, t.counter_id, t.status, t.priority, t.created_at, t.called_at, t.completed_at, t.wait_time, t.service_time, t.daily_sequence, t.queue_date, t.notes, c.id, c.name, c.prefix, c.color_code, co.id, co.number, co.name FROM tickets t LEFT JOIN categories c ON t.category_id = c.id LEFT JOIN counters co ON t.counter_id = co.id WHERE t.category_id IN (%s) ORDER BY t.created_at DESC`, strings.Join(placeholders, ","))
}

func (q *TicketQueries) CancelYesterdayWaiting() string {
	return `UPDATE tickets SET status = 'cancelled' WHERE queue_date = CURRENT_DATE - INTERVAL '1 day' AND status = 'waiting'`
}
