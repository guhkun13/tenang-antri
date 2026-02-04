package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"tenangantri/internal/model"
	"tenangantri/internal/query"
)

// TicketRepository handles ticket data operations
type TicketRepository struct {
	ticketQueries *query.TicketQueries
}

func NewTicketRepository(pool *pgxpool.Pool) *TicketRepository {
	return &TicketRepository{
		ticketQueries: query.NewTicketQueries(pool),
	}
}

// GetByID retrieves a ticket by ID
func (r *TicketRepository) GetByID(ctx context.Context, id int) (*model.Ticket, error) {
	row := r.ticketQueries.GetTicketByID(ctx, id)

	ticket := &model.Ticket{}
	err := row.Scan(
		&ticket.ID, &ticket.TicketNumber, &ticket.CategoryID, &ticket.CounterID,
		&ticket.Status, &ticket.Priority, &ticket.CreatedAt, &ticket.CalledAt,
		&ticket.CompletedAt, &ticket.WaitTime, &ticket.ServiceTime, &ticket.Notes,
	)
	if err != nil {
		return nil, err
	}
	return ticket, nil
}

// GetWithDetails retrieves a ticket with full details
func (r *TicketRepository) GetWithDetails(ctx context.Context, id int) (*model.Ticket, error) {
	row := r.ticketQueries.GetTicketWithDetails(ctx, id)

	ticket := &model.Ticket{Category: &model.Category{}, Counter: &model.Counter{}}

	var catID *int
	var catName, catPrefix, catColor *string
	var coID *int
	var coNumber, coName *string

	err := row.Scan(
		&ticket.ID, &ticket.TicketNumber, &ticket.CategoryID, &ticket.CounterID,
		&ticket.Status, &ticket.Priority, &ticket.CreatedAt, &ticket.CalledAt,
		&ticket.CompletedAt, &ticket.WaitTime, &ticket.ServiceTime, &ticket.Notes,
		&catID, &catName, &catPrefix, &catColor,
		&coID, &coNumber, &coName,
	)
	if err != nil {
		return nil, err
	}

	if catID != nil {
		ticket.Category.ID = *catID
		ticket.Category.Name = *catName
		ticket.Category.Prefix = *catPrefix
		ticket.Category.ColorCode = *catColor
	}
	if coID != nil {
		ticket.Counter.ID = *coID
		ticket.Counter.Number = *coNumber
		ticket.Counter.Name = *coName
	}

	return ticket, nil
}

// GetByTicketNumber retrieves a ticket with full details by ticket number
func (r *TicketRepository) GetByTicketNumber(ctx context.Context, ticketNumber string) (*model.Ticket, error) {
	row := r.ticketQueries.GetTicketByNumber(ctx, ticketNumber)

	ticket := &model.Ticket{Category: &model.Category{}, Counter: &model.Counter{}}

	var catID *int
	var catName, catPrefix, catColor *string
	var coID *int
	var coNumber, coName *string

	err := row.Scan(
		&ticket.ID, &ticket.TicketNumber, &ticket.CategoryID, &ticket.CounterID,
		&ticket.Status, &ticket.Priority, &ticket.CreatedAt, &ticket.CalledAt,
		&ticket.CompletedAt, &ticket.WaitTime, &ticket.ServiceTime, &ticket.Notes,
		&catID, &catName, &catPrefix, &catColor,
		&coID, &coNumber, &coName,
	)
	if err != nil {
		return nil, err
	}

	if catID != nil {
		ticket.Category.ID = *catID
		ticket.Category.Name = *catName
		ticket.Category.Prefix = *catPrefix
		ticket.Category.ColorCode = *catColor
	}
	if coID != nil {
		ticket.Counter.ID = *coID
		ticket.Counter.Number = *coNumber
		ticket.Counter.Name = *coName
	}

	return ticket, nil
}

// Create creates a new ticket
func (r *TicketRepository) Create(ctx context.Context, ticket *model.Ticket) (*model.Ticket, error) {
	id, createdAt, err := r.ticketQueries.CreateTicket(
		ctx,
		ticket.TicketNumber,
		ticket.CategoryID,
		ticket.Status,
		ticket.Priority,
		ticket.Notes,
	)
	if err != nil {
		return nil, err
	}

	ticket.ID = id
	ticket.CreatedAt = createdAt
	return ticket, nil
}

// UpdateStatus updates the status of a ticket
func (r *TicketRepository) UpdateStatus(ctx context.Context, id int, status string) error {
	return r.ticketQueries.UpdateTicketStatus(ctx, id, status)
}

// AssignToCounter assigns a ticket to a counter
func (r *TicketRepository) AssignToCounter(ctx context.Context, ticketID, counterID int) error {
	return r.ticketQueries.AssignTicketToCounter(ctx, ticketID, counterID)
}

// GetNextTicket retrieves the next ticket in queue for given categories
func (r *TicketRepository) GetNextTicket(ctx context.Context, categoryIDs []int) (*model.Ticket, error) {
	if len(categoryIDs) == 0 {
		return nil, fmt.Errorf("no categories provided")
	}

	row := r.ticketQueries.GetNextTicket(ctx, categoryIDs)

	ticket := &model.Ticket{}
	err := row.Scan(
		&ticket.ID, &ticket.TicketNumber, &ticket.CategoryID,
		&ticket.Status, &ticket.Priority, &ticket.CreatedAt, &ticket.Notes,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return ticket, nil
}

// GetCurrentForCounter retrieves the current ticket being served at a counter
func (r *TicketRepository) GetCurrentForCounter(ctx context.Context, counterID int) (*model.Ticket, error) {
	row := r.ticketQueries.GetCurrentTicketForCounter(ctx, counterID)

	ticket := &model.Ticket{}
	err := row.Scan(
		&ticket.ID, &ticket.TicketNumber, &ticket.CategoryID, &ticket.CounterID,
		&ticket.Status, &ticket.Priority, &ticket.CreatedAt, &ticket.CalledAt,
		&ticket.CompletedAt, &ticket.WaitTime, &ticket.ServiceTime, &ticket.Notes,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		log.Error().Err(err).Msg("Error getting current ticket for counter")
		return nil, err
	}
	return ticket, nil
}

// List retrieves tickets with optional filters
func (r *TicketRepository) List(ctx context.Context, filters map[string]interface{}) ([]model.Ticket, error) {
	rows, err := r.ticketQueries.ListTickets(ctx, filters)
	if err != nil {
		log.Error().Err(err).Str("layer", "repository").Str("func", "List").Msg("Failed to query tickets")
		return nil, err
	}
	defer rows.Close()

	var tickets []model.Ticket
	for rows.Next() {
		var t model.Ticket
		t.Category = &model.Category{}
		t.Counter = &model.Counter{}

		var catName, catPrefix, catColor *string
		var coNumber, coName *string

		err := rows.Scan(
			&t.ID, &t.TicketNumber, &t.CategoryID, &t.CounterID, &t.Status, &t.Priority,
			&t.CreatedAt, &t.CalledAt, &t.CompletedAt, &t.WaitTime, &t.ServiceTime, &t.Notes,
			&catName, &catPrefix, &catColor,
			&coNumber, &coName,
		)
		if err != nil {
			log.Error().Err(err).Str("layer", "repository").Str("func", "List").Msg("Failed to scan ticket")
			return nil, err
		}

		if catName != nil {
			t.Category.Name = *catName
			t.Category.Prefix = *catPrefix
			t.Category.ColorCode = *catColor
		}
		if coNumber != nil {
			t.Counter.Number = *coNumber
			t.Counter.Name = *coName
		} else {
			t.Counter = nil
		}

		tickets = append(tickets, t)
	}

	return tickets, nil
}

// GetTodayCount retrieves today's ticket count
func (r *TicketRepository) GetTodayCount(ctx context.Context) (int, error) {
	var count int
	err := r.ticketQueries.GetTodayTicketCount(ctx).Scan(&count)
	return count, err
}

// GetTodayCountByCategory retrieves today's ticket count for a specific category
func (r *TicketRepository) GetTodayCountByCategory(ctx context.Context, categoryID int) (int, error) {
	var count int
	err := r.ticketQueries.GetTodayTicketCountByCategory(ctx, categoryID).Scan(&count)
	return count, err
}

// GenerateNumber generates a unique ticket number for a given prefix
func (r *TicketRepository) GenerateNumber(ctx context.Context, prefix string) (string, error) {
	var number int
	err := r.ticketQueries.GenerateTicketNumber(ctx, prefix).Scan(&number)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%03d", prefix, number), nil
}

// GetWaitingPreview retrieves a preview of waiting tickets
func (r *TicketRepository) GetWaitingPreview(ctx context.Context, limit int) ([]model.Ticket, error) {
	rows, err := r.ticketQueries.GetWaitingTicketsPreview(ctx, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Ticket])
	if err != nil {
		log.Error().Err(err).Str("layer", "repository").Str("func", "GetWaitingPreview").Msg("Failed to collect rows")
		return nil, err
	}

	return res, nil
}

// GetWaitingPreviewByCategories retrieves waiting tickets preview for specific categories
func (r *TicketRepository) GetWaitingPreviewByCategories(ctx context.Context, categoryIDs []int, limit int) ([]model.Ticket, error) {
	rows, err := r.ticketQueries.GetWaitingTicketsPreviewByCategories(ctx, categoryIDs, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Ticket])
	if err != nil {
		log.Error().Err(err).Str("layer", "repository").Str("func", "GetWaitingPreviewByCategories").Msg("Failed to collect rows")
		return nil, err
	}

	return res, nil
}

// GetTodayCompletedByCategories retrieves completed tickets today for specific categories
func (r *TicketRepository) GetTodayCompletedByCategories(ctx context.Context, categoryIDs []int) ([]model.Ticket, error) {
	rows, err := r.ticketQueries.GetTodayCompletedTicketsByCategories(ctx, categoryIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []model.Ticket
	for rows.Next() {
		ticket := model.Ticket{Category: &model.Category{}, Counter: &model.Counter{}}
		var catID *int
		var catName, catPrefix, catColor *string
		var coID *int
		var coNumber, coName *string

		err := rows.Scan(
			&ticket.ID, &ticket.TicketNumber, &ticket.CategoryID, &ticket.CounterID,
			&ticket.Status, &ticket.Priority, &ticket.CreatedAt, &ticket.CalledAt,
			&ticket.CompletedAt, &ticket.WaitTime, &ticket.ServiceTime, &ticket.Notes,
			&catID, &catName, &catPrefix, &catColor,
			&coID, &coNumber, &coName,
		)
		if err != nil {
			log.Error().Err(err).Str("layer", "repository").Str("func", "GetTodayCompletedByCategories").Msg("Failed to scan row")
			return nil, err
		}

		if catID != nil {
			ticket.Category.ID = *catID
			ticket.Category.Name = *catName
			ticket.Category.Prefix = *catPrefix
			ticket.Category.ColorCode = *catColor
		}
		if coID != nil {
			ticket.Counter.ID = *coID
			ticket.Counter.Number = *coNumber
			ticket.Counter.Name = *coName
		}

		tickets = append(tickets, ticket)
	}

	return tickets, nil
}

// GetLastCalledByCategoryID retrieves the last called ticket number for a specific category
func (r *TicketRepository) GetLastCalledByCategoryID(ctx context.Context, categoryID int) (string, error) {
	var ticketNumber string
	err := r.ticketQueries.GetLastCalledTicketByCategory(ctx, categoryID).Scan(&ticketNumber)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return ticketNumber, nil
}
