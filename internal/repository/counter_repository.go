package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"

	"tenangantri/internal/model"
	"tenangantri/internal/query"
)

type CounterRepository interface {
	GetByID(ctx context.Context, id int) (*model.Counter, error)
	Create(ctx context.Context, counter *model.Counter) (*model.Counter, error)
	Update(ctx context.Context, counter *model.Counter) (*model.Counter, error)
	UpdateStatus(ctx context.Context, id int, status string) error
	UpdateStaff(ctx context.Context, counterID int, staffID *int) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context) ([]model.Counter, error)
}

type counterRepository struct {
	pool       DB
	counterQry *query.CounterQueries
}

func NewCounterRepository(pool DB) CounterRepository {
	return &counterRepository{
		pool:       pool,
		counterQry: query.NewCounterQueries(),
	}
}

func (r *counterRepository) GetByID(ctx context.Context, id int) (*model.Counter, error) {
	queryStr := r.counterQry.GetCounterByID(ctx)
	row := r.pool.QueryRow(ctx, queryStr, id)

	counter := &model.Counter{}
	var catID, staffID sql.NullInt64
	err := row.Scan(
		&counter.ID, &counter.Number, &counter.Name, &counter.Location,
		&counter.Status, &catID,
		&counter.CreatedAt, &counter.UpdatedAt, &staffID,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		log.Error().Err(err).Int("id", id).Msg("Failed to scan counter")
		return nil, err
	}

	if catID.Valid {
		counter.CategoryID = sql.NullInt64{Int64: catID.Int64, Valid: true}
	}
	if staffID.Valid {
		counter.CurrentStaffID = sql.NullInt64{Int64: staffID.Int64, Valid: true}
	}

	return counter, nil
}

func (r *counterRepository) Create(ctx context.Context, counter *model.Counter) (*model.Counter, error) {
	queryStr := r.counterQry.CreateCounter(ctx)
	var id int
	var createdAt, updatedAt time.Time
	err := r.pool.QueryRow(ctx, queryStr, counter.Number, counter.Name, counter.Location, counter.Status, counter.CategoryID).Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	counter.ID = id
	counter.CreatedAt = createdAt
	counter.UpdatedAt = updatedAt
	return counter, nil
}

func (r *counterRepository) Update(ctx context.Context, counter *model.Counter) (*model.Counter, error) {
	queryStr := r.counterQry.UpdateCounter(ctx)
	_, err := r.pool.Exec(ctx, queryStr, counter.Number, counter.Name, counter.Location, counter.Status, counter.CategoryID, counter.ID)
	if err != nil {
		return nil, err
	}
	return r.GetByID(ctx, counter.ID)
}

func (r *counterRepository) UpdateStatus(ctx context.Context, id int, status string) error {
	queryStr := r.counterQry.UpdateCounterStatus(ctx)
	_, err := r.pool.Exec(ctx, queryStr, status, id)
	return err
}

func (r *counterRepository) UpdateStaff(ctx context.Context, counterID int, staffID *int) error {
	queryStr := r.counterQry.UpdateCounterStaff(ctx)
	_, err := r.pool.Exec(ctx, queryStr, staffID, counterID)
	return err
}

func (r *counterRepository) Delete(ctx context.Context, id int) error {
	queryStr := r.counterQry.DeleteCounter(ctx)
	_, err := r.pool.Exec(ctx, queryStr, id)
	return err
}

func (r *counterRepository) List(ctx context.Context) ([]model.Counter, error) {
	queryStr := r.counterQry.ListCounters(ctx)
	rows, err := r.pool.Query(ctx, queryStr)
	if err != nil {
		log.Error().Err(err).Str("layer", "repository").Str("func", "List").Msg("Failed to list counters")
		return nil, err
	}
	defer rows.Close()

	counters, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Counter])
	if err != nil {
		log.Error().Err(err).Str("layer", "repository").Str("func", "List").Msg("Failed to collect rows")
		return nil, err
	}

	return counters, nil
}
