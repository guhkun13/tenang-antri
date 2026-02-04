package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"tenangantri/internal/model"
	"tenangantri/internal/query"
)

type TicketRepository struct {
	pool      *pgxpool.Pool
	ticketQry *query.TicketQueries
}

func NewTicketRepository(pool *pgxpool.Pool) *TicketRepository {
	return &TicketRepository{
		pool:      pool,
		ticketQry: query.NewTicketQueries(),
	}
}

func (r *TicketRepository) GetByID(ctx context.Context, id int) (*model.Ticket, error) {
	sql := r.ticketQry.GetTicketByID(ctx)
	row := r.pool.QueryRow(ctx, sql, id)

	ticket := &model.Ticket{}
	var catID int
	err := row.Scan(
		&ticket.ID, &ticket.TicketNumber, &catID, &ticket.CounterID,
		&ticket.Status, &ticket.Priority, &ticket.CreatedAt, &ticket.CalledAt,
		&ticket.CompletedAt, &ticket.WaitTime, &ticket.ServiceTime, &ticket.Notes,
	)
	if err != nil {
		return nil, err
	}
	return ticket, nil
}

func (r *TicketRepository) GetWithDetails(ctx context.Context, id int) (*model.Ticket, error) {
	sql := r.ticketQry.GetTicketWithDetails(ctx)
	row := r.pool.QueryRow(ctx, sql, id)

	ticket := &model.Ticket{Category: &model.Category{}, Counter: &model.Counter{}}
	var catID int
	var catName, catPrefix, catColor string
	var coID *int
	var coNumber, coName string

	err := row.Scan(
		&ticket.ID, &ticket.TicketNumber, &catID, &ticket.CounterID,
		&ticket.Status, &ticket.Priority, &ticket.CreatedAt, &ticket.CalledAt,
		&ticket.CompletedAt, &ticket.WaitTime, &ticket.ServiceTime, &ticket.Notes,
		&catID, &catName, &catPrefix, &catColor,
		&coID, &coNumber, &coName,
	)
	if err != nil {
		return nil, err
	}

	ticket.Category.ID = catID
	ticket.Category.Name = catName
	ticket.Category.Prefix = catPrefix
	ticket.Category.ColorCode = catColor
	if coID != nil {
		ticket.Counter.ID = *coID
		ticket.Counter.Number = coNumber
		ticket.Counter.Name = coName
	}

	return ticket, nil
}

func (r *TicketRepository) GetByTicketNumber(ctx context.Context, ticketNumber string) (*model.Ticket, error) {
	sql := r.ticketQry.GetTicketByNumber(ctx)
	row := r.pool.QueryRow(ctx, sql, ticketNumber)

	ticket := &model.Ticket{Category: &model.Category{}, Counter: &model.Counter{}}
	var catID int
	var catName, catPrefix, catColor string
	var coID *int
	var coNumber, coName string

	err := row.Scan(
		&ticket.ID, &ticket.TicketNumber, &catID, &ticket.CounterID,
		&ticket.Status, &ticket.Priority, &ticket.CreatedAt, &ticket.CalledAt,
		&ticket.CompletedAt, &ticket.WaitTime, &ticket.ServiceTime, &ticket.Notes,
		&catID, &catName, &catPrefix, &catColor,
		&coID, &coNumber, &coName,
	)
	if err != nil {
		return nil, err
	}

	ticket.Category.ID = catID
	ticket.Category.Name = catName
	ticket.Category.Prefix = catPrefix
	ticket.Category.ColorCode = catColor
	if coID != nil {
		ticket.Counter.ID = *coID
		ticket.Counter.Number = coNumber
		ticket.Counter.Name = coName
	}

	return ticket, nil
}

func (r *TicketRepository) Create(ctx context.Context, ticket *model.Ticket) (*model.Ticket, error) {
	sql := r.ticketQry.CreateTicket(ctx)
	var id int
	var createdAt time.Time
	err := r.pool.QueryRow(ctx, sql, ticket.TicketNumber, ticket.Category.ID, ticket.Status, ticket.Priority, ticket.Notes).Scan(&id, &createdAt)
	if err != nil {
		return nil, err
	}

	ticket.ID = id
	ticket.CreatedAt = createdAt
	return ticket, nil
}

func (r *TicketRepository) UpdateStatus(ctx context.Context, id int, status string) error {
	sql := r.ticketQry.UpdateTicketStatus(ctx, status)
	_, err := r.pool.Exec(ctx, sql, status, id)
	return err
}

func (r *TicketRepository) AssignToCounter(ctx context.Context, ticketID, counterID int) error {
	sql := r.ticketQry.AssignTicketToCounter(ctx)
	_, err := r.pool.Exec(ctx, sql, counterID, ticketID)
	return err
}

func (r *TicketRepository) GetNextTicket(ctx context.Context, categoryIDs []int) (*model.Ticket, error) {
	if len(categoryIDs) == 0 {
		return nil, fmt.Errorf("no categories provided")
	}

	sql := r.ticketQry.GetNextTicket(ctx, categoryIDs)
	row := r.pool.QueryRow(ctx, sql, categoryIDs)

	ticket := &model.Ticket{}
	var catID int
	err := row.Scan(
		&ticket.ID, &ticket.TicketNumber, &catID,
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

func (r *TicketRepository) GetCurrentForCounter(ctx context.Context, counterID int) (*model.Ticket, error) {
	sql := r.ticketQry.GetCurrentTicketForCounter(ctx)
	row := r.pool.QueryRow(ctx, sql, counterID)

	ticket := &model.Ticket{}
	var catID int
	err := row.Scan(
		&ticket.ID, &ticket.TicketNumber, &catID, &ticket.CounterID,
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

func (r *TicketRepository) List(ctx context.Context, filters map[string]interface{}) ([]model.Ticket, error) {
	sql := r.ticketQry.ListTickets(ctx, filters)
	rows, err := r.pool.Query(ctx, sql)
	if err != nil {
		log.Error().Err(err).Str("layer", "repository").Str("func", "List").Msg("Failed to query tickets")
		return nil, err
	}
	defer rows.Close()

	tickets, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (model.Ticket, error) {
		var t model.Ticket
		t.Category = &model.Category{}
		t.Counter = &model.Counter{}

		var catID int
		var catName, catPrefix, catColor string
		var coNumber, coName string

		err := row.Scan(
			&t.ID, &t.TicketNumber, &catID, &t.CounterID, &t.Status, &t.Priority,
			&t.CreatedAt, &t.CalledAt, &t.CompletedAt, &t.WaitTime, &t.ServiceTime, &t.Notes,
			&catName, &catPrefix, &catColor,
			&coNumber, &coName,
		)
		if err != nil {
			return model.Ticket{}, err
		}

		t.Category.ID = catID
		t.Category.Name = catName
		t.Category.Prefix = catPrefix
		t.Category.ColorCode = catColor
		if coNumber != "" {
			t.Counter.Number = coNumber
			t.Counter.Name = coName
		} else {
			t.Counter = nil
		}

		return t, nil
	})

	if err != nil {
		log.Error().Err(err).Str("layer", "repository").Str("func", "List").Msg("Failed to collect rows")
		return nil, err
	}

	return tickets, nil
}

func (r *TicketRepository) GetTodayCount(ctx context.Context) (int, error) {
	sql := r.ticketQry.GetTodayTicketCount(ctx)
	var count int
	err := r.pool.QueryRow(ctx, sql).Scan(&count)
	return count, err
}

func (r *TicketRepository) GetTodayCountByCategory(ctx context.Context, categoryID int) (int, error) {
	sql := r.ticketQry.GetTodayTicketCountByCategory(ctx)
	var count int
	err := r.pool.QueryRow(ctx, sql, categoryID).Scan(&count)
	return count, err
}

func (r *TicketRepository) GenerateNumber(ctx context.Context, prefix string) (string, error) {
	sql := r.ticketQry.GenerateTicketNumber(ctx)
	var number int
	err := r.pool.QueryRow(ctx, sql, prefix+"%").Scan(&number)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%03d", prefix, number), nil
}

func (r *TicketRepository) GetWaitingPreview(ctx context.Context, limit int) ([]model.Ticket, error) {
	sql := r.ticketQry.GetWaitingTicketsPreview(ctx)
	rows, err := r.pool.Query(ctx, sql, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tickets, err := pgx.CollectRows(rows, pgx.RowToStructByPos[model.Ticket])
	if err != nil {
		log.Error().Err(err).Str("layer", "repository").Str("func", "GetWaitingPreview").Msg("Failed to collect rows")
		return nil, err
	}

	return tickets, nil
}

func (r *TicketRepository) GetWaitingPreviewByCategories(ctx context.Context, categoryIDs []int, limit int) ([]model.Ticket, error) {
	sql := r.ticketQry.GetWaitingTicketsPreviewByCategories(ctx, categoryIDs)
	args := make([]any, 0, len(categoryIDs)+1)
	args = append(args, limit)
	for _, id := range categoryIDs {
		args = append(args, id)
	}
	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tickets, err := pgx.CollectRows(rows, pgx.RowToStructByPos[model.Ticket])
	if err != nil {
		log.Error().Err(err).Str("layer", "repository").Str("func", "GetWaitingPreviewByCategories").Msg("Failed to collect rows")
		return nil, err
	}

	return tickets, nil
}

func (r *TicketRepository) GetTodayCompletedByCategories(ctx context.Context, categoryIDs []int) ([]model.Ticket, error) {
	sql := r.ticketQry.GetTodayCompletedTicketsByCategories(ctx, categoryIDs)
	args := make([]any, len(categoryIDs))
	for i, id := range categoryIDs {
		args[i] = id
	}
	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tickets, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (model.Ticket, error) {
		var t model.Ticket
		t.Category = &model.Category{}
		t.Counter = &model.Counter{}

		var catID int
		var catName, catPrefix, catColor string
		var coID *int
		var coNumber, coName string

		err := row.Scan(
			&t.ID, &t.TicketNumber, &catID, &t.CounterID,
			&t.Status, &t.Priority, &t.CreatedAt, &t.CalledAt,
			&t.CompletedAt, &t.WaitTime, &t.ServiceTime, &t.Notes,
			&catID, &catName, &catPrefix, &catColor,
			&coID, &coNumber, &coName,
		)
		if err != nil {
			return model.Ticket{}, err
		}

		t.Category.ID = catID
		t.Category.Name = catName
		t.Category.Prefix = catPrefix
		t.Category.ColorCode = catColor
		if coID != nil {
			t.Counter.ID = *coID
			t.Counter.Number = coNumber
			t.Counter.Name = coName
		}

		return t, nil
	})

	if err != nil {
		log.Error().Err(err).Str("layer", "repository").Str("func", "GetTodayCompletedByCategories").Msg("Failed to collect rows")
		return nil, err
	}

	return tickets, nil
}

func (r *TicketRepository) GetLastCalledByCategoryID(ctx context.Context, categoryID int) (string, error) {
	sql := r.ticketQry.GetLastCalledTicketByCategory(ctx)
	var ticketNumber string
	err := r.pool.QueryRow(ctx, sql, categoryID).Scan(&ticketNumber)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return ticketNumber, nil
}
