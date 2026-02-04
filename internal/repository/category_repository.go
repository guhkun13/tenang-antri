package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"tenangantri/internal/model"
	"tenangantri/internal/query"
)

type CategoryRepository struct {
	pool        *pgxpool.Pool
	categoryQry *query.CategoryQueries
}

func NewCategoryRepository(pool *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{
		pool:        pool,
		categoryQry: query.NewCategoryQueries(),
	}
}

func (r *CategoryRepository) GetByID(ctx context.Context, id int) (*model.Category, error) {
	sql := r.categoryQry.GetCategoryByID(ctx)
	row := r.pool.QueryRow(ctx, sql, id)

	category := &model.Category{}
	err := row.Scan(
		&category.ID, &category.Name, &category.Prefix, &category.Priority,
		&category.ColorCode, &category.Description, &category.Icon,
		&category.IsActive, &category.CreatedAt, &category.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (r *CategoryRepository) Create(ctx context.Context, category *model.Category) (*model.Category, error) {
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

func (r *CategoryRepository) Update(ctx context.Context, category *model.Category) (*model.Category, error) {
	sql := r.categoryQry.UpdateCategory(ctx)
	_, err := r.pool.Exec(ctx, sql, category.Name, category.Prefix, category.Priority, category.ColorCode, category.Description, category.Icon, category.IsActive, category.ID)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (r *CategoryRepository) Delete(ctx context.Context, id int) error {
	sql := r.categoryQry.DeleteCategory(ctx)
	_, err := r.pool.Exec(ctx, sql, id)
	return err
}

func (r *CategoryRepository) List(ctx context.Context, activeOnly bool) ([]model.Category, error) {
	sql := r.categoryQry.ListCategories(ctx, activeOnly)
	rows, err := r.pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[model.Category])
}
