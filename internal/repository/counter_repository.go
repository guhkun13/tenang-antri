package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"queue-system/internal/model"
	"queue-system/internal/query"
)

// CounterRepository handles counter data operations
type CounterRepository struct {
	counterQueries *query.CounterQueries
}

func NewCounterRepository(pool *pgxpool.Pool) *CounterRepository {
	return &CounterRepository{
		counterQueries: query.NewCounterQueries(pool),
	}
}

// GetByID retrieves a counter by ID
func (r *CounterRepository) GetByID(ctx context.Context, id int) (*model.Counter, error) {
	row := r.counterQueries.GetCounterByID(ctx, id)

	counter := &model.Counter{}
	err := row.Scan(
		&counter.ID, &counter.Number, &counter.Name, &counter.Location,
		&counter.Status, &counter.IsActive, &counter.CurrentStaffID,
		&counter.CreatedAt, &counter.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return counter, nil
}

// Create creates a new counter
func (r *CounterRepository) Create(ctx context.Context, counter *model.Counter) (*model.Counter, error) {
	id, createdAt, updatedAt, err := r.counterQueries.CreateCounter(
		ctx,
		counter.Number,
		counter.Name,
		counter.Location,
		counter.Status,
		counter.IsActive,
	)
	if err != nil {
		return nil, err
	}

	counter.ID = id
	counter.CreatedAt = createdAt
	counter.UpdatedAt = updatedAt
	return counter, nil
}

// Update updates counter information
func (r *CounterRepository) Update(ctx context.Context, counter *model.Counter) (*model.Counter, error) {
	err := r.counterQueries.UpdateCounter(
		ctx,
		counter.Number,
		counter.Name,
		counter.Location,
		counter.Status,
		counter.IsActive,
		counter.ID,
	)
	if err != nil {
		return nil, err
	}

	return counter, nil
}

// UpdateStatus updates only the status of a counter
func (r *CounterRepository) UpdateStatus(ctx context.Context, id int, status string) error {
	return r.counterQueries.UpdateCounterStatus(ctx, status, id)
}

// UpdateStaff updates the current staff assigned to a counter
func (r *CounterRepository) UpdateStaff(ctx context.Context, counterID int, staffID *int) error {
	return r.counterQueries.UpdateCounterStaff(ctx, counterID, staffID)
}

// Delete deletes a counter
func (r *CounterRepository) Delete(ctx context.Context, id int) error {
	return r.counterQueries.DeleteCounter(ctx, id)
}

// List retrieves all counters
func (r *CounterRepository) List(ctx context.Context) ([]model.Counter, error) {
	rows, err := r.counterQueries.ListCounters(ctx)
	if err != nil {
		log.Error().Err(err).Str("layer", "repository").Str("func", "List").Msg("Failed to list counters")
		return nil, err
	}
	defer rows.Close()

	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Counter])
	if err != nil {
		log.Error().Err(err).Str("layer", "repository").Str("func", "List").Msg("Failed to collect rows")
		return nil, err
	}

	return res, nil
}

// GetCategories retrieves categories for a specific counter
func (r *CounterRepository) GetCategories(ctx context.Context, counterID int) ([]model.Category, error) {
	rows, err := r.counterQueries.GetCounterCategories(ctx, counterID)
	if err != nil {
		log.Error().Err(err).Str("layer", "repository").Str("func", "GetCategories").Msg("Failed to get counter categories")
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[model.Category])
}

// AssignCategory assigns a category to a counter
func (r *CounterRepository) AssignCategory(ctx context.Context, counterID, categoryID int) error {
	return r.counterQueries.AssignCategoryToCounter(ctx, counterID, categoryID)
}

// RemoveCategory removes a category from a counter
func (r *CounterRepository) RemoveCategory(ctx context.Context, counterID, categoryID int) error {
	return r.counterQueries.RemoveCategoryFromCounter(ctx, counterID, categoryID)
}

// ClearCategories removes all categories from a counter
func (r *CounterRepository) ClearCategories(ctx context.Context, counterID int) error {
	return r.counterQueries.ClearCounterCategories(ctx, counterID)
}
