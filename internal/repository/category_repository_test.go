package repository

import (
	"context"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"

	"tenangantri/internal/query"
)

func TestCategoryRepository_GetByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	categoryQry := query.NewCategoryQueries()
	repo := &categoryRepository{
		pool:        mock,
		categoryQry: categoryQry,
	}

	catID := 1
	now := time.Now()
	rows := pgxmock.NewRows([]string{"id", "name", "prefix", "priority", "color_code", "description", "icon", "is_active", "created_at", "updated_at"}).
		AddRow(catID, "General", "A", 1, "#000000", "General Service", "box", true, now, now)

	mock.ExpectQuery("SELECT id, name, prefix").
		WithArgs(catID).
		WillReturnRows(rows)

	ctx := context.Background()
	cat, err := repo.GetByID(ctx, catID)

	assert.NoError(t, err)
	assert.NotNil(t, cat)
	assert.Equal(t, "General", cat.Name)
	assert.Equal(t, "A", cat.Prefix)

	assert.NoError(t, mock.ExpectationsWereMet())
}
