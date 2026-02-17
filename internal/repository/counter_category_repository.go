package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"

	"tenangantri/internal/model"
	"tenangantri/internal/query"
)

type CounterCategoryRepository interface {
	Create(ctx context.Context, counterID, categoryID int) (*model.CounterCategory, error)
	GetByID(ctx context.Context, id int) (*model.CounterCategory, error)
	GetByCounterID(ctx context.Context, counterID int) ([]model.CounterCategory, error)
	GetCategoryIDsByCounterID(ctx context.Context, counterID int) ([]int, error)
	GetByCategoryID(ctx context.Context, categoryID int) ([]model.CounterCategory, error)
	GetCounterIDsByCategoryID(ctx context.Context, categoryID int) ([]int, error)
	DeleteByID(ctx context.Context, id int) error
	DeleteByCounterID(ctx context.Context, counterID int) error
	DeleteByCategoryID(ctx context.Context, categoryID int) error
	DeleteByCounterAndCategory(ctx context.Context, counterID, categoryID int) error
	ListAll(ctx context.Context) ([]model.CounterCategory, error)
}

type counterCategoryRepository struct {
	pool DB
	qry  *query.CounterCategoryQueries
}

func NewCounterCategoryRepository(pool DB) CounterCategoryRepository {
	return &counterCategoryRepository{
		pool: pool,
		qry:  query.NewCounterCategoryQueries(),
	}
}

func (r *counterCategoryRepository) Create(ctx context.Context, counterID, categoryID int) (*model.CounterCategory, error) {
	queryStr := r.qry.Create(ctx)
	var id int
	var createdAt, updatedAt time.Time
	err := r.pool.QueryRow(ctx, queryStr, counterID, categoryID).Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		log.Error().Err(err).Str("layer", "repository").Int("counter_id", counterID).Int("category_id", categoryID).Msg("Failed to create counter_category")
		return nil, err
	}

	return &model.CounterCategory{
		ID:         id,
		CounterID:  counterID,
		CategoryID: categoryID,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}, nil
}

func (r *counterCategoryRepository) GetByID(ctx context.Context, id int) (*model.CounterCategory, error) {
	queryStr := r.qry.GetByID(ctx)
	row := r.pool.QueryRow(ctx, queryStr, id)

	cs := &model.CounterCategory{}
	err := row.Scan(&cs.ID, &cs.CounterID, &cs.CategoryID, &cs.CreatedAt, &cs.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		log.Error().Err(err).Str("layer", "repository").Int("id", id).Msg("Failed to get counter_category by id")
		return nil, err
	}
	return cs, nil
}

func (r *counterCategoryRepository) GetByCounterID(ctx context.Context, counterID int) ([]model.CounterCategory, error) {
	queryStr := r.qry.GetByCounterID(ctx)
	rows, err := r.pool.Query(ctx, queryStr, counterID)
	if err != nil {
		log.Error().Err(err).Str("layer", "repository").Int("counter_id", counterID).Msg("Failed to get counter_category by counter_id")
		return nil, err
	}
	defer rows.Close()

	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.CounterCategory])
	if err != nil {
		log.Error().Err(err).Str("layer", "repository").Msg("Failed to collect counter_category rows")
		return nil, err
	}

	return res, nil
}

func (r *counterCategoryRepository) GetCategoryIDsByCounterID(ctx context.Context, counterID int) ([]int, error) {
	queryStr := r.qry.GetCategoryIDsByCounterID(ctx)
	rows, err := r.pool.Query(ctx, queryStr, counterID)
	if err != nil {
		log.Error().Err(err).Str("layer", "repository").Int("counter_id", counterID).Msg("Failed to get category_ids by counter_id")
		return nil, err
	}
	defer rows.Close()

	var categoryIDs []int
	for rows.Next() {
		var catID int
		if err := rows.Scan(&catID); err != nil {
			continue
		}
		categoryIDs = append(categoryIDs, catID)
	}

	return categoryIDs, nil
}

func (r *counterCategoryRepository) GetByCategoryID(ctx context.Context, categoryID int) ([]model.CounterCategory, error) {
	queryStr := r.qry.GetByCategoryID(ctx)
	rows, err := r.pool.Query(ctx, queryStr, categoryID)
	if err != nil {
		log.Error().Err(err).Str("layer", "repository").Int("category_id", categoryID).Msg("Failed to get counter_category by category_id")
		return nil, err
	}
	defer rows.Close()

	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.CounterCategory])
	if err != nil {
		log.Error().Err(err).Str("layer", "repository").Msg("Failed to collect counter_category rows")
		return nil, err
	}

	return res, nil
}

func (r *counterCategoryRepository) GetCounterIDsByCategoryID(ctx context.Context, categoryID int) ([]int, error) {
	queryStr := r.qry.GetCounterIDsByCategoryID(ctx)
	rows, err := r.pool.Query(ctx, queryStr, categoryID)
	if err != nil {
		log.Error().Err(err).Str("layer", "repository").Int("category_id", categoryID).Msg("Failed to get counter_ids by category_id")
		return nil, err
	}
	defer rows.Close()

	var counterIDs []int
	for rows.Next() {
		var counterID int
		if err := rows.Scan(&counterID); err != nil {
			continue
		}
		counterIDs = append(counterIDs, counterID)
	}

	return counterIDs, nil
}

func (r *counterCategoryRepository) DeleteByID(ctx context.Context, id int) error {
	queryStr := r.qry.DeleteByID(ctx)
	_, err := r.pool.Exec(ctx, queryStr, id)
	if err != nil {
		log.Error().Err(err).Str("layer", "repository").Int("id", id).Msg("Failed to delete counter_category by id")
		return err
	}
	return nil
}

func (r *counterCategoryRepository) DeleteByCounterID(ctx context.Context, counterID int) error {
	queryStr := r.qry.DeleteByCounterID(ctx)
	_, err := r.pool.Exec(ctx, queryStr, counterID)
	if err != nil {
		log.Error().Err(err).Str("layer", "repository").Int("counter_id", counterID).Msg("Failed to delete counter_category by counter_id")
		return err
	}
	return nil
}

func (r *counterCategoryRepository) DeleteByCategoryID(ctx context.Context, categoryID int) error {
	queryStr := r.qry.DeleteByCategoryID(ctx)
	_, err := r.pool.Exec(ctx, queryStr, categoryID)
	if err != nil {
		log.Error().Err(err).Str("layer", "repository").Int("category_id", categoryID).Msg("Failed to delete counter_category by category_id")
		return err
	}
	return nil
}

func (r *counterCategoryRepository) DeleteByCounterAndCategory(ctx context.Context, counterID, categoryID int) error {
	queryStr := r.qry.DeleteByCounterAndCategory(ctx)
	_, err := r.pool.Exec(ctx, queryStr, counterID, categoryID)
	if err != nil {
		log.Error().Err(err).Str("layer", "repository").Int("counter_id", counterID).Int("category_id", categoryID).Msg("Failed to delete counter_category by counter and category")
		return err
	}
	return nil
}

func (r *counterCategoryRepository) ListAll(ctx context.Context) ([]model.CounterCategory, error) {
	queryStr := r.qry.ListAll(ctx)
	rows, err := r.pool.Query(ctx, queryStr)
	if err != nil {
		log.Error().Err(err).Str("layer", "repository").Msg("Failed to list counter_category")
		return nil, err
	}
	defer rows.Close()

	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.CounterCategory])
	if err != nil {
		log.Error().Err(err).Str("layer", "repository").Msg("Failed to collect counter_category rows")
		return nil, err
	}

	return res, nil
}
