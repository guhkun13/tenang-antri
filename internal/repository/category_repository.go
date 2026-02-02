package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"queue-system/internal/model"
	"queue-system/internal/query"
)

// CategoryRepository handles category data operations
type CategoryRepository struct {
	categoryQueries *query.CategoryQueries
}

func NewCategoryRepository(pool *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{
		categoryQueries: query.NewCategoryQueries(pool),
	}
}

// GetByID retrieves a category by ID
func (r *CategoryRepository) GetByID(ctx context.Context, id int) (*model.Category, error) {
	row := r.categoryQueries.GetCategoryByID(ctx, id)

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

// Create creates a new category
func (r *CategoryRepository) Create(ctx context.Context, category *model.Category) (*model.Category, error) {
	id, createdAt, updatedAt, err := r.categoryQueries.CreateCategory(
		ctx,
		category.Name,
		category.Prefix,
		category.Priority,
		category.ColorCode,
		category.Description,
		category.Icon,
		category.IsActive,
	)
	if err != nil {
		return nil, err
	}

	category.ID = id
	category.CreatedAt = createdAt
	category.UpdatedAt = updatedAt
	return category, nil
}

// Update updates category information
func (r *CategoryRepository) Update(ctx context.Context, category *model.Category) (*model.Category, error) {
	err := r.categoryQueries.UpdateCategory(
		ctx,
		category.Name,
		category.Prefix,
		category.Priority,
		category.ColorCode,
		category.Description,
		category.Icon,
		category.IsActive,
		category.ID,
	)
	if err != nil {
		return nil, err
	}

	return category, nil
}

// Delete deletes a category
func (r *CategoryRepository) Delete(ctx context.Context, id int) error {
	return r.categoryQueries.DeleteCategory(ctx, id)
}

// List retrieves categories with optional active filter
func (r *CategoryRepository) List(ctx context.Context, activeOnly bool) ([]model.Category, error) {
	rows, err := r.categoryQueries.ListCategories(ctx, activeOnly)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[model.Category])
}
