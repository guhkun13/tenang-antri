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

	ticket.CategoryID = sql.NullInt64{Int64: int64(catID), Valid: catID > 0}
	if coID.Valid {
		ticket.CounterID = sql.NullInt64{Int64: coID.Int64, Valid: true}
	}
	if calledAt.Valid {
		ticket.CalledAt = calledAt
	}
	if completedAt.Valid {
		ticket.CompletedAt = completedAt
	}
	if waitTime.Valid {
		ticket.WaitTime = waitTime
	}
	if serviceTime.Valid {
		ticket.ServiceTime = serviceTime
	}

	return ticket, nil
}

func (r *ticketRepository) GetWithDetails(ctx context.Context, id int) (*model.Ticket, error) {
	queryStr := r.ticketQry.GetTicketWithDetails(ctx)
	row := r.pool.QueryRow(ctx, queryStr, id)

	ticket := &model.Ticket{}

	var catID int
	var catName, catPrefix, catColor sql.NullString
	var coID sql.NullInt64
	var coNumber, coName sql.NullString
	var waitTime, serviceTime sql.NullInt64
	var calledAt, completedAt sql.NullTime

	err := row.Scan(
		&ticket.ID, &ticket.TicketNumber, &catID, &coID,
		&ticket.Status, &ticket.Priority, &ticket.CreatedAt, &calledAt,
		&completedAt, &waitTime, &serviceTime, &ticket.DailySequence, &ticket.QueueDate, &ticket.Notes,
		&catID, &catName, &catPrefix, &catColor,
		&coID, &coNumber, &coName,
	)
	if err != nil {
		return nil, err
	}

	ticket.CategoryID = sql.NullInt64{Int64: int64(catID), Valid: catID > 0}
	if catID > 0 {
		ticket.Category = &model.Category{
			ID:        catID,
			Name:      catName.String,
			Prefix:    catPrefix.String,
			ColorCode: catColor.String,
		}
	}
	if coID.Valid {
		ticket.CounterID = sql.NullInt64{Int64: coID.Int64, Valid: true}
		ticket.Counter = &model.Counter{
			ID:     int(coID.Int64),
			Number: coNumber.String,
			Name:   sql.NullString{String: coName.String, Valid: coName.Valid},
		}
	}
	if calledAt.Valid {
		ticket.CalledAt = calledAt
	}
	if completedAt.Valid {
		ticket.CompletedAt = completedAt
	}
	if waitTime.Valid {
		ticket.WaitTime = waitTime
	}
	if serviceTime.Valid {
		ticket.ServiceTime = serviceTime
	}

	return ticket, nil
}

func (r *ticketRepository) GetByTicketNumber(ctx context.Context, ticketNumber string) (*model.Ticket, error) {
	queryStr := r.ticketQry.GetTicketByNumber(ctx)
	row := r.pool.QueryRow(ctx, queryStr, ticketNumber)

	ticket := &model.Ticket{}

	var catID int
	var catName, catPrefix, catColor sql.NullString
	var coID sql.NullInt64
	var coNumber, coName sql.NullString
	var waitTime, serviceTime sql.NullInt64
	var calledAt, completedAt sql.NullTime

	err := row.Scan(
		&ticket.ID, &ticket.TicketNumber, &catID, &coID,
		&ticket.Status, &ticket.Priority, &ticket.CreatedAt, &calledAt,
		&completedAt, &waitTime, &serviceTime, &ticket.DailySequence, &ticket.QueueDate, &ticket.Notes,
		&catID, &catName, &catPrefix, &catColor,
		&coID, &coNumber, &coName,
	)
	if err != nil {
		return nil, err
	}

	ticket.CategoryID = sql.NullInt64{Int64: int64(catID), Valid: catID > 0}
	if catID > 0 {
		ticket.Category = &model.Category{
			ID:        catID,
			Name:      catName.String,
			Prefix:    catPrefix.String,
			ColorCode: catColor.String,
		}
	}
	if coID.Valid {
		ticket.CounterID = sql.NullInt64{Int64: coID.Int64, Valid: true}
		ticket.Counter = &model.Counter{
			ID:     int(coID.Int64),
			Number: coNumber.String,
			Name:   sql.NullString{String: coName.String, Valid: coName.Valid},
		}
	}
	if calledAt.Valid {
		ticket.CalledAt = calledAt
	}
	if completedAt.Valid {
		ticket.CompletedAt = completedAt
	}
	if waitTime.Valid {
		ticket.WaitTime = waitTime
	}
	if serviceTime.Valid {
		ticket.ServiceTime = serviceTime
	}

	return ticket, nil
}

func (r *ticketRepository) Create(ctx context.Context, ticket *model.Ticket) (*model.Ticket, error) {
	queryStr := r.ticketQry.CreateTicket(ctx)
	var id int
	var createdAt time.Time
	err := r.pool.QueryRow(ctx, queryStr, ticket.TicketNumber, ticket.CategoryID, ticket.Status, ticket.Priority, ticket.Notes, ticket.DailySequence, ticket.QueueDate).Scan(&id, &createdAt)
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
	err := row.Scan(
		&ticket.ID, &ticket.TicketNumber, &ticket.CategoryID,
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
		var notes sql.NullString

		err := row.Scan(
			&t.ID, &t.TicketNumber, &t.CategoryID, &t.CounterID, &t.Status, &t.Priority,
			&t.CreatedAt, &t.CalledAt, &t.CompletedAt, &t.WaitTime, &t.ServiceTime, &t.DailySequence, &t.QueueDate, &notes,
		)
		if err != nil {
			return model.Ticket{}, err
		}

		if notes.Valid {
			t.Notes = notes
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
	queryStr := r.ticketQry.GetTodayCompletedTicketsByCategories(ctx, categoryIDs)
	args := make([]any, len(categoryIDs))
	for i, id := range categoryIDs {
		args[i] = id
	}
	rows, err := r.pool.Query(ctx, queryStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tickets, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (model.Ticket, error) {
		var t model.Ticket
		var catID int
		var catName, catPrefix, catColor sql.NullString
		var coID sql.NullInt64
		var coNumber, coName sql.NullString
		var waitTime, serviceTime sql.NullInt64
		var calledAt, completedAt sql.NullTime

		err := row.Scan(
			&t.ID, &t.TicketNumber, &catID, &coID,
			&t.Status, &t.Priority, &t.CreatedAt, &calledAt,
			&completedAt, &waitTime, &serviceTime, &t.DailySequence, &t.QueueDate, &t.Notes,
			&catID, &catName, &catPrefix, &catColor,
			&coID, &coNumber, &coName,
		)
		if err != nil {
			return model.Ticket{}, err
		}

		t.CategoryID = sql.NullInt64{Int64: int64(catID), Valid: catID > 0}
		if catID > 0 {
			t.Category = &model.Category{
				ID:        catID,
				Name:      catName.String,
				Prefix:    catPrefix.String,
				ColorCode: catColor.String,
			}
		}
		if coID.Valid {
			t.CounterID = sql.NullInt64{Int64: coID.Int64, Valid: true}
			t.Counter = &model.Counter{
				ID:     int(coID.Int64),
				Number: coNumber.String,
				Name:   sql.NullString{String: coName.String, Valid: coName.Valid},
			}
		}
		if calledAt.Valid {
			t.CalledAt = calledAt
		}
		if completedAt.Valid {
			t.CompletedAt = completedAt
		}
		if waitTime.Valid {
			t.WaitTime = waitTime
		}
		if serviceTime.Valid {
			t.ServiceTime = serviceTime
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
