package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"tenangantri/internal/model"
	"tenangantri/internal/query"
)

type CounterRepository struct {
	pool       *pgxpool.Pool
	counterQry *query.CounterQueries
}

func NewCounterRepository(pool *pgxpool.Pool) *CounterRepository {
	return &CounterRepository{
		pool:       pool,
		counterQry: query.NewCounterQueries(),
	}
}

func (r *CounterRepository) GetByID(ctx context.Context, id int) (*model.Counter, error) {
	sql := r.counterQry.GetCounterByID(ctx)
	row := r.pool.QueryRow(ctx, sql, id)

	counter := &model.Counter{}
	err := row.Scan(
		&counter.ID, &counter.Number, &counter.Name, &counter.Location,
		&counter.Status, &counter.IsActive, &counter.CategoryID,
		&counter.CurrentStaffID, &counter.CreatedAt, &counter.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return counter, nil
}

func (r *CounterRepository) Create(ctx context.Context, counter *model.Counter) (*model.Counter, error) {
	sql := r.counterQry.CreateCounter(ctx)
	var id int
	var createdAt, updatedAt time.Time
	err := r.pool.QueryRow(ctx, sql, counter.Number, counter.Name, counter.Location, counter.Status, counter.IsActive, counter.CategoryID).Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	counter.ID = id
	counter.CreatedAt = createdAt
	counter.UpdatedAt = updatedAt
	return counter, nil
}

func (r *CounterRepository) Update(ctx context.Context, counter *model.Counter) (*model.Counter, error) {
	sql := r.counterQry.UpdateCounter(ctx)
	_, err := r.pool.Exec(ctx, sql, counter.Number, counter.Name, counter.Location, counter.Status, counter.IsActive, counter.CategoryID, counter.ID)
	if err != nil {
		return nil, err
	}
	return counter, nil
}

func (r *CounterRepository) UpdateStatus(ctx context.Context, id int, status string) error {
	sql := r.counterQry.UpdateCounterStatus(ctx)
	_, err := r.pool.Exec(ctx, sql, status, id)
	return err
}

func (r *CounterRepository) UpdateStaff(ctx context.Context, counterID int, staffID *int) error {
	sql := r.counterQry.UpdateCounterStaff(ctx)
	_, err := r.pool.Exec(ctx, sql, staffID, counterID)
	return err
}

func (r *CounterRepository) Delete(ctx context.Context, id int) error {
	sql := r.counterQry.DeleteCounter(ctx)
	_, err := r.pool.Exec(ctx, sql, id)
	return err
}

func (r *CounterRepository) List(ctx context.Context) ([]model.Counter, error) {
	sql := r.counterQry.ListCounters(ctx)
	rows, err := r.pool.Query(ctx, sql)
	if err != nil {
		log.Error().Err(err).Str("layer", "repository").Str("func", "List").Msg("Failed to list counters")
		return nil, err
	}
	defer rows.Close()

	var counters []model.Counter
	for rows.Next() {
		counter := model.Counter{}
		err := rows.Scan(
			&counter.ID, &counter.Number, &counter.Name, &counter.Location,
			&counter.Status, &counter.IsActive, &counter.CategoryID,
			&counter.CurrentStaffID, &counter.CreatedAt, &counter.UpdatedAt,
		)
		if err != nil {
			log.Error().Err(err).Str("layer", "repository").Str("func", "List").Msg("Failed to scan counter")
			return nil, err
		}
		counters = append(counters, counter)
	}

	if err := rows.Err(); err != nil {
		log.Error().Err(err).Str("layer", "repository").Str("func", "List").Msg("Error iterating counters")
		return nil, err
	}

	return counters, nil
}
