package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"

	"tenangantri/internal/model"
	"tenangantri/internal/query"
)

type CategoryRepository interface {
	GetByID(ctx context.Context, id int) (*model.Category, error)
	Create(ctx context.Context, category *model.Category) (*model.Category, error)
	Update(ctx context.Context, category *model.Category) (*model.Category, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, activeOnly bool) ([]model.Category, error)
}

type categoryRepository struct {
	pool        DB
	categoryQry *query.CategoryQueries
}

func NewCategoryRepository(pool DB) CategoryRepository {
	return &categoryRepository{
		pool:        pool,
		categoryQry: query.NewCategoryQueries(),
	}
}

func (r *categoryRepository) GetByID(ctx context.Context, id int) (*model.Category, error) {
	queryStr := r.categoryQry.GetCategoryByID(ctx)
	row := r.pool.QueryRow(ctx, queryStr, id)

	cat := &model.Category{}
	err := row.Scan(
		&cat.ID, &cat.Name, &cat.Prefix, &cat.Priority,
		&cat.ColorCode, &cat.Description, &cat.Icon, &cat.IsActive,
		&cat.CreatedAt, &cat.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		log.Error().Err(err).Int("id", id).Msg("Failed to scan category")
		return nil, err
	}
	return cat, nil
}

func (r *categoryRepository) Create(ctx context.Context, category *model.Category) (*model.Category, error) {
	sql := r.categoryQry.CreateCategory(ctx)
	var id int
	var createdAt, updatedAt time.Time
	err := r.pool.QueryRow(ctx, sql, category.Name, category.Prefix, category.Priority, category.ColorCode, category.Description, category.Icon, category.IsActive).Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	category.ID = id
	category.CreatedAt = createdAt
	category.UpdatedAt = updatedAt
	return category, nil
}

func (r *categoryRepository) Update(ctx context.Context, category *model.Category) (*model.Category, error) {
	sql := r.categoryQry.UpdateCategory(ctx)
	_, err := r.pool.Exec(ctx, sql, category.Name, category.Prefix, category.Priority, category.ColorCode, category.Description, category.Icon, category.IsActive, category.ID)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (r *categoryRepository) Delete(ctx context.Context, id int) error {
	sql := r.categoryQry.DeleteCategory(ctx)
	_, err := r.pool.Exec(ctx, sql, id)
	return err
}

func (r *categoryRepository) List(ctx context.Context, activeOnly bool) ([]model.Category, error) {
	sql := r.categoryQry.ListCategories(ctx, activeOnly)
	rows, err := r.pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[model.Category])
}
