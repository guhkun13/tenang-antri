package repository

import (
	"context"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"

	"tenangantri/internal/query"
)

func TestUserRepository_GetByUsername(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	userQry := query.NewUserQueries()
	repo := &userRepository{
		pool:    mock,
		userQry: userQry,
	}

	username := "admin"
	now := time.Now()
	rows := pgxmock.NewRows([]string{"id", "username", "full_name", "email", "phone", "role", "is_active", "created_at", "updated_at", "last_login"}).
		AddRow(1, username, "Admin User", "admin@example.com", "123456", "admin", true, now, now, nil)

	mock.ExpectQuery("SELECT id, username, full_name").
		WithArgs(username).
		WillReturnRows(rows)

	ctx := context.Background()
	user, err := repo.GetByUsername(ctx, username)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, username, user.Username)
	assert.Equal(t, "admin", user.Role)

	assert.NoError(t, mock.ExpectationsWereMet())
}
