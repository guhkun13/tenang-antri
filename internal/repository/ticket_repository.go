package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"

	"tenangantri/internal/model"
	"tenangantri/internal/query"
)

type TicketRepository interface {
	GetByID(ctx context.Context, id int) (*model.Ticket, error)
	GetWithDetails(ctx context.Context, id int) (*model.Ticket, error)
	GetByTicketNumber(ctx context.Context, ticketNumber string) (*model.Ticket, error)
	Create(ctx context.Context, ticket *model.Ticket) (*model.Ticket, error)
	UpdateStatus(ctx context.Context, id int, status string) error
	AssignToCounter(ctx context.Context, ticketID, counterID int) error
	GetNextTicket(ctx context.Context, categoryIDs []int) (*model.Ticket, error)
	GetCurrentForCounter(ctx context.Context, counterID int) (*model.Ticket, error)
	List(ctx context.Context, filters map[string]interface{}) ([]model.Ticket, error)
	GetTodayCount(ctx context.Context) (int, error)
	GetTodayCountByCategory(ctx context.Context, categoryID int) (int, error)
	GenerateNumber(ctx context.Context, categoryID int, prefix string) (string, int, error)
	GetWaitingPreview(ctx context.Context, limit int) ([]model.Ticket, error)
	GetWaitingPreviewByCategories(ctx context.Context, categoryIDs []int, limit int) ([]model.Ticket, error)
	GetTodayCompletedByCategories(ctx context.Context, categoryIDs []int) ([]model.Ticket, error)
	GetLastCalledByCategoryID(ctx context.Context, categoryID int) (string, error)
}

type ticketRepository struct {
	pool      DB
	ticketQry *query.TicketQueries
}

func NewTicketRepository(pool DB) TicketRepository {
	return &ticketRepository{
		pool:      pool,
		ticketQry: query.NewTicketQueries(),
	}
}

func (r *ticketRepository) GetByID(ctx context.Context, id int) (*model.Ticket, error) {
	queryStr := r.ticketQry.GetTicketByID(ctx)
	row := r.pool.QueryRow(ctx, queryStr, id)

	ticket := &model.Ticket{}
	var catID int
	var coID sql.NullInt64
	var waitTime, serviceTime sql.NullInt64
	var calledAt, completedAt sql.NullTime

	err := row.Scan(
		&ticket.ID, &ticket.TicketNumber, &catID, &coID,
		&ticket.Status, &ticket.Priority, &ticket.CreatedAt, &calledAt,
		&completedAt, &waitTime, &serviceTime, &ticket.DailySequence, &ticket.QueueDate, &ticket.Notes,
	)
	if err != nil {
		return nil, err
	}

	ticket.CategoryID = &catID
	if coID.Valid {
		val := int(coID.Int64)
		ticket.CounterID = &val
	}
	if calledAt.Valid {
		ticket.CalledAt = &calledAt.Time
	}
	if completedAt.Valid {
		ticket.CompletedAt = &completedAt.Time
	}
	if waitTime.Valid {
		val := int(waitTime.Int64)
		ticket.WaitTime = &val
	}
	if serviceTime.Valid {
		val := int(serviceTime.Int64)
		ticket.ServiceTime = &val
	}

	return ticket, nil
}

func (r *ticketRepository) GetWithDetails(ctx context.Context, id int) (*model.Ticket, error) {
	queryStr := r.ticketQry.GetTicketWithDetails(ctx)
	row := r.pool.QueryRow(ctx, queryStr, id)

	ticket := &model.Ticket{Category: &model.Category{}, Counter: &model.Counter{}}
	var catID, catIDFromJoin int
	var catName, catPrefix, catColor string
	var coID *int
	var coNumber, coName *string

	err := row.Scan(
		&ticket.ID, &ticket.TicketNumber, &catID, &ticket.CounterID,
		&ticket.Status, &ticket.Priority, &ticket.CreatedAt, &ticket.CalledAt,
		&ticket.CompletedAt, &ticket.WaitTime, &ticket.ServiceTime, &ticket.DailySequence, &ticket.QueueDate, &ticket.Notes,
		&catIDFromJoin, &catName, &catPrefix, &catColor,
		&coID, &coNumber, &coName,
	)
	if err != nil {
		return nil, err
	}

	ticket.Category.ID = catIDFromJoin
	ticket.Category.Name = catName
	ticket.Category.Prefix = catPrefix
	ticket.Category.ColorCode = catColor
	if coID != nil {
		ticket.Counter.ID = *coID
		if coNumber != nil {
			ticket.Counter.Number = *coNumber
		}
		if coName != nil {
			ticket.Counter.Name = *coName
		}
	}

	return ticket, nil
}

func (r *ticketRepository) GetByTicketNumber(ctx context.Context, ticketNumber string) (*model.Ticket, error) {
	queryStr := r.ticketQry.GetTicketByNumber(ctx)
	row := r.pool.QueryRow(ctx, queryStr, ticketNumber)

	ticket := &model.Ticket{Category: &model.Category{}, Counter: &model.Counter{}}
	var catID, catIDFromJoin int
	var catName, catPrefix, catColor string
	var coID *int
	var coNumber, coName *string

	err := row.Scan(
		&ticket.ID, &ticket.TicketNumber, &catID, &ticket.CounterID,
		&ticket.Status, &ticket.Priority, &ticket.CreatedAt, &ticket.CalledAt,
		&ticket.CompletedAt, &ticket.WaitTime, &ticket.ServiceTime, &ticket.DailySequence, &ticket.QueueDate, &ticket.Notes,
		&catIDFromJoin, &catName, &catPrefix, &catColor,
		&coID, &coNumber, &coName,
	)
	if err != nil {
		return nil, err
	}

	ticket.Category.ID = catIDFromJoin
	ticket.Category.Name = catName
	ticket.Category.Prefix = catPrefix
	ticket.Category.ColorCode = catColor
	if coID != nil {
		ticket.Counter.ID = *coID
		if coNumber != nil {
			ticket.Counter.Number = *coNumber
		}
		if coName != nil {
			ticket.Counter.Name = *coName
		}
	}

	return ticket, nil
}

func (r *ticketRepository) Create(ctx context.Context, ticket *model.Ticket) (*model.Ticket, error) {
	queryStr := r.ticketQry.CreateTicket(ctx)
	var id int
	var createdAt time.Time
	err := r.pool.QueryRow(ctx, queryStr, ticket.TicketNumber, ticket.Category.ID, ticket.Status, ticket.Priority, ticket.Notes, ticket.DailySequence, ticket.QueueDate).Scan(&id, &createdAt)
	if err != nil {
		return nil, err
	}

	ticket.ID = id
	ticket.CreatedAt = createdAt
	return ticket, nil
}

func (r *ticketRepository) UpdateStatus(ctx context.Context, id int, status string) error {
	queryStr := r.ticketQry.UpdateTicketStatus(ctx, status)
	_, err := r.pool.Exec(ctx, queryStr, status, id)
	return err
}

func (r *ticketRepository) AssignToCounter(ctx context.Context, ticketID, counterID int) error {
	queryStr := r.ticketQry.AssignTicketToCounter(ctx)
	_, err := r.pool.Exec(ctx, queryStr, counterID, ticketID)
	return err
}

func (r *ticketRepository) GetNextTicket(ctx context.Context, categoryIDs []int) (*model.Ticket, error) {
	if len(categoryIDs) == 0 {
		return nil, fmt.Errorf("no categories provided")
	}

	sql := r.ticketQry.GetNextTicket(ctx, categoryIDs)
	row := r.pool.QueryRow(ctx, sql, categoryIDs)

	ticket := &model.Ticket{}
	var catID int
	err := row.Scan(
		&ticket.ID, &ticket.TicketNumber, &catID,
		&ticket.Status, &ticket.Priority, &ticket.CreatedAt, &ticket.DailySequence, &ticket.QueueDate, &ticket.Notes,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return ticket, nil
}

func (r *ticketRepository) GetCurrentForCounter(ctx context.Context, counterID int) (*model.Ticket, error) {
	queryStr := r.ticketQry.GetCurrentTicketForCounter(ctx)
	row := r.pool.QueryRow(ctx, queryStr, counterID)

	ticket := &model.Ticket{}
	err := row.Scan(
		&ticket.ID, &ticket.TicketNumber, &ticket.CategoryID, &ticket.CounterID,
		&ticket.Status, &ticket.Priority, &ticket.CreatedAt, &ticket.CalledAt,
		&ticket.CompletedAt, &ticket.WaitTime, &ticket.ServiceTime, &ticket.DailySequence, &ticket.QueueDate, &ticket.Notes,
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

func (r *ticketRepository) List(ctx context.Context, filters map[string]interface{}) ([]model.Ticket, error) {
	result := r.ticketQry.ListTickets(ctx, filters)
	rows, err := r.pool.Query(ctx, result.Query, result.Args...)
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
		var coNumber, coName *string
		var notes *string

		err := row.Scan(
			&t.ID, &t.TicketNumber, &catID, &t.CounterID, &t.Status, &t.Priority,
			&t.CreatedAt, &t.CalledAt, &t.CompletedAt, &t.WaitTime, &t.ServiceTime, &t.DailySequence, &t.QueueDate, &notes,
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
		if notes != nil {
			t.Notes = *notes
		}
		if coNumber != nil && *coNumber != "" {
			t.Counter.Number = *coNumber
			t.Counter.Name = *coName
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

func (r *ticketRepository) GetTodayCount(ctx context.Context) (int, error) {
	sql := r.ticketQry.GetTodayTicketCount(ctx)
	var count int
	err := r.pool.QueryRow(ctx, sql).Scan(&count)
	return count, err
}

func (r *ticketRepository) GetTodayCountByCategory(ctx context.Context, categoryID int) (int, error) {
	sql := r.ticketQry.GetTodayTicketCountByCategory(ctx)
	var count int
	err := r.pool.QueryRow(ctx, sql, categoryID).Scan(&count)
	return count, err
}

func (r *ticketRepository) GenerateNumber(ctx context.Context, categoryID int, prefix string) (string, int, error) {
	sql := r.ticketQry.GenerateTicketNumber(ctx)
	var number int
	err := r.pool.QueryRow(ctx, sql, categoryID).Scan(&number)
	if err != nil {
		return "", 0, err
	}
	return fmt.Sprintf("%s%03d", prefix, number), number, nil
}

func (r *ticketRepository) GetWaitingPreview(ctx context.Context, limit int) ([]model.Ticket, error) {
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

func (r *ticketRepository) GetWaitingPreviewByCategories(ctx context.Context, categoryIDs []int, limit int) ([]model.Ticket, error) {
	if len(categoryIDs) == 0 {
		return []model.Ticket{}, nil
	}
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

func (r *ticketRepository) GetTodayCompletedByCategories(ctx context.Context, categoryIDs []int) ([]model.Ticket, error) {
	if len(categoryIDs) == 0 {
		return []model.Ticket{}, nil
	}
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

		var catID, catIDFromJoin int
		var catName, catPrefix, catColor string
		var coID *int
		var coNumber, coName *string

		err := row.Scan(
			&t.ID, &t.TicketNumber, &catID, &t.CounterID,
			&t.Status, &t.Priority, &t.CreatedAt, &t.CalledAt,
			&t.CompletedAt, &t.WaitTime, &t.ServiceTime, &t.DailySequence, &t.QueueDate, &t.Notes,
			&catIDFromJoin, &catName, &catPrefix, &catColor,
			&coID, &coNumber, &coName,
		)
		if err != nil {
			return model.Ticket{}, err
		}

		t.Category.ID = catIDFromJoin
		t.Category.Name = catName
		t.Category.Prefix = catPrefix
		t.Category.ColorCode = catColor
		if coID != nil {
			t.Counter.ID = *coID
			if coNumber != nil {
				t.Counter.Number = *coNumber
			}
			if coName != nil {
				t.Counter.Name = *coName
			}
		}

		return t, nil
	})

	if err != nil {
		log.Error().Err(err).Str("layer", "repository").Str("func", "GetTodayCompletedByCategories").Msg("Failed to collect rows")
		return nil, err
	}

	return tickets, nil
}

func (r *ticketRepository) GetLastCalledByCategoryID(ctx context.Context, categoryID int) (string, error) {
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
