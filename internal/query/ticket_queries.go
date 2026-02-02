package query

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TicketQueries contains all ticket-related SQL queries
type TicketQueries struct {
	pool *pgxpool.Pool
}

func NewTicketQueries(pool *pgxpool.Pool) *TicketQueries {
	return &TicketQueries{pool: pool}
}

// CreateTicket inserts a new ticket
func (q *TicketQueries) CreateTicket(ctx context.Context, ticketNumber string, categoryID int, status string, priority int, notes string) (int, time.Time, error) {
	query := `
		INSERT INTO tickets (ticket_number, category_id, status, priority, notes)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`

	var id int
	var createdAt time.Time
	err := q.pool.QueryRow(ctx, query, ticketNumber, categoryID, status, priority, notes).
		Scan(&id, &createdAt)
	return id, createdAt, err
}

// GetTicketByID retrieves a ticket by ID
func (q *TicketQueries) GetTicketByID(ctx context.Context, id int) pgx.Row {
	query := `
		SELECT t.id, t.ticket_number, t.category_id, t.counter_id, t.status, t.priority,
		       t.created_at, t.called_at, t.completed_at, t.wait_time, t.service_time, t.notes
		FROM tickets t WHERE t.id = $1`
	return q.pool.QueryRow(ctx, query, id)
}

// GetTicketWithDetails retrieves a ticket with full details
func (q *TicketQueries) GetTicketWithDetails(ctx context.Context, id int) pgx.Row {
	query := `
		SELECT t.id, t.ticket_number, t.category_id, t.counter_id, t.status, t.priority,
		       t.created_at, t.called_at, t.completed_at, t.wait_time, t.service_time, t.notes,
		       c.id, c.name, c.prefix, c.color_code,
		       co.id, co.number, co.name
		FROM tickets t
		LEFT JOIN categories c ON t.category_id = c.id
		LEFT JOIN counters co ON t.counter_id = co.id
		WHERE t.id = $1`
	return q.pool.QueryRow(ctx, query, id)
}

// UpdateTicketStatus updates the status of a ticket
func (q *TicketQueries) UpdateTicketStatus(ctx context.Context, id int, status string) error {
	var query string

	switch status {
	case "serving":
		query = `UPDATE tickets SET status = $1, called_at = NOW() WHERE id = $2`
	case "completed", "no_show":
		query = `
			UPDATE tickets 
			SET status = $1, completed_at = NOW(),
			    wait_time = EXTRACT(EPOCH FROM (called_at - created_at))::INT,
			    service_time = EXTRACT(EPOCH FROM (NOW() - called_at))::INT
			WHERE id = $2`
	default:
		query = `UPDATE tickets SET status = $1 WHERE id = $2`
	}

	_, err := q.pool.Exec(ctx, query, status, id)
	return err
}

// AssignTicketToCounter assigns a ticket to a counter
func (q *TicketQueries) AssignTicketToCounter(ctx context.Context, ticketID, counterID int) error {
	query := `UPDATE tickets SET counter_id = $1, status = 'serving', called_at = NOW() WHERE id = $2`
	_, err := q.pool.Exec(ctx, query, counterID, ticketID)
	return err
}

// GetNextTicket retrieves the next ticket in queue for given categories
func (q *TicketQueries) GetNextTicket(ctx context.Context, categoryIDs []int) pgx.Row {
	if len(categoryIDs) == 0 {
		return nil
	}

	query := `
		SELECT t.id, t.ticket_number, t.category_id, t.status, t.priority, t.created_at, t.notes
		FROM tickets t
		WHERE t.category_id = ANY($1) AND t.status = 'waiting'
		ORDER BY 
			(SELECT priority FROM categories WHERE id = t.category_id) DESC,
			t.created_at ASC
		LIMIT 1`

	return q.pool.QueryRow(ctx, query, categoryIDs)
}

// GetCurrentTicketForCounter retrieves the current ticket being served at a counter
func (q *TicketQueries) GetCurrentTicketForCounter(ctx context.Context, counterID int) pgx.Row {
	query := `
		SELECT t.id, t.ticket_number, t.category_id, t.counter_id, t.status, t.priority,
		       t.created_at, t.called_at, t.completed_at, t.wait_time, t.service_time, t.notes
		FROM tickets t
		WHERE t.counter_id = $1 AND t.status = 'serving'
		ORDER BY t.called_at DESC
		LIMIT 1`
	return q.pool.QueryRow(ctx, query, counterID)
}

// ListTickets retrieves tickets with optional filters
func (q *TicketQueries) ListTickets(ctx context.Context, filters map[string]interface{}) (pgx.Rows, error) {
	query := `
		SELECT t.id, t.ticket_number, t.category_id, t.counter_id, t.status, t.priority,
		       t.created_at, t.called_at, t.completed_at, t.wait_time, t.service_time, t.notes
		FROM tickets t WHERE 1=1`
	var args []interface{}
	argCount := 1

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
		query += fmt.Sprintf(" AND t.created_at <= $%d", argCount)
		args = append(args, dateTo)
		argCount++
	}

	query += ` ORDER BY t.created_at DESC`

	if limit, ok := filters["limit"]; ok && limit != 0 {
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, limit)
		argCount++
	}
	if offset, ok := filters["offset"]; ok && offset != 0 {
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, offset)
	}

	return q.pool.Query(ctx, query, args...)
}

// GetTodayTicketCount retrieves today's ticket count
func (q *TicketQueries) GetTodayTicketCount(ctx context.Context) pgx.Row {
	query := `SELECT COUNT(*) FROM tickets WHERE DATE(created_at) = CURRENT_DATE`
	return q.pool.QueryRow(ctx, query)
}

// GetTodayTicketCountByCategory retrieves today's ticket count for a specific category
func (q *TicketQueries) GetTodayTicketCountByCategory(ctx context.Context, categoryID int) pgx.Row {
	query := `SELECT COUNT(*) FROM tickets WHERE category_id = $1 AND DATE(created_at) = CURRENT_DATE`
	return q.pool.QueryRow(ctx, query, categoryID)
}

// GenerateTicketNumber generates a unique ticket number for a given prefix
func (q *TicketQueries) GenerateTicketNumber(ctx context.Context, prefix string) pgx.Row {
	query := `
		SELECT COALESCE(MAX(NULLIF(regexp_replace(ticket_number, '[^0-9]', '', 'g'), ''))::INT, 0) + 1
		FROM tickets 
		WHERE ticket_number LIKE $1 AND DATE(created_at) = CURRENT_DATE`
	return q.pool.QueryRow(ctx, query, prefix+"%")
}

// GetWaitingTicketsPreview retrieves a preview of waiting tickets
func (q *TicketQueries) GetWaitingTicketsPreview(ctx context.Context, limit int) (pgx.Rows, error) {
	query := `
		SELECT t.id, t.ticket_number, t.category_id, t.counter_id, t.status, t.priority, 
		t.created_at, t.called_at, t.completed_at, t.wait_time, t.service_time, t.notes 
		FROM tickets t
		WHERE t.status = 'waiting'
		ORDER BY 
			(SELECT priority FROM categories WHERE id = t.category_id) DESC,
			t.created_at ASC
		LIMIT $1`
	return q.pool.Query(ctx, query, limit)
}

// GetWaitingTicketsPreviewByCategories retrieves waiting tickets preview for specific categories
func (q *TicketQueries) GetWaitingTicketsPreviewByCategories(ctx context.Context, categoryIDs []int, limit int) (pgx.Rows, error) {
	if len(categoryIDs) == 0 {
		return nil, fmt.Errorf("no categories provided")
	}

	placeholders := make([]string, len(categoryIDs))
	args := make([]interface{}, len(categoryIDs)+1)
	for i, id := range categoryIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+2)
		args[i+1] = id
	}

	query := fmt.Sprintf(`
		SELECT t.id, t.ticket_number, t.category_id, t.counter_id, t.status, t.priority, 
		t.created_at, t.called_at, t.completed_at, t.wait_time, t.service_time, t.notes 
		FROM tickets t
		WHERE t.status = 'waiting' AND t.category_id IN (%s)
		ORDER BY 
			(SELECT priority FROM categories WHERE id = t.category_id) DESC,
			t.created_at ASC
		LIMIT $1`, strings.Join(placeholders, ","))

	args[0] = limit

	return q.pool.Query(ctx, query, args...)
}

// GetTodayCompletedTicketsByCategories retrieves completed tickets today for specific categories
func (q *TicketQueries) GetTodayCompletedTicketsByCategories(ctx context.Context, categoryIDs []int) (pgx.Rows, error) {
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
		SELECT t.id, t.ticket_number, t.category_id, t.counter_id, t.status, t.priority,
		       t.created_at, t.called_at, t.completed_at, t.wait_time, t.service_time, t.notes,
		       c.id, c.name, c.prefix, c.color_code,
		       co.id, co.number, co.name
		FROM tickets t
		LEFT JOIN categories c ON t.category_id = c.id
		LEFT JOIN counters co ON t.counter_id = co.id
		WHERE t.status = 'completed' 
		  AND t.completed_at >= CURRENT_DATE
		  AND t.category_id IN (%s)
		ORDER BY t.completed_at DESC
		LIMIT 50`, strings.Join(placeholders, ","))

	return q.pool.Query(ctx, query, args...)
}
