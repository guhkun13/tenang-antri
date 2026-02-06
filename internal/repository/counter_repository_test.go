package repository

import (
	"context"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"

	"tenangantri/internal/query"
)

func TestCounterRepository_GetByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	counterQry := query.NewCounterQueries()
	repo := &counterRepository{
		pool:       mock,
		counterQry: counterQry,
	}

	counterID := 1
	now := time.Now()
	rows := pgxmock.NewRows([]string{"id", "number", "name", "location", "status", "category_id", "current_staff_id", "created_at", "updated_at"}).
		AddRow(counterID, "1", "Counter 1", "Main Hall", "active", 1, nil, now, now)

	mock.ExpectQuery("SELECT id, number, name").
		WithArgs(counterID).
		WillReturnRows(rows)

	ctx := context.Background()
	counter, err := repo.GetByID(ctx, counterID)

	assert.NoError(t, err)
	assert.NotNil(t, counter)
	assert.Equal(t, "1", counter.Number)
	assert.Equal(t, "active", counter.Status)

	assert.NoError(t, mock.ExpectationsWereMet())
}
