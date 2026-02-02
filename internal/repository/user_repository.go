package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"queue-system/internal/model"
	"queue-system/internal/query"
)

// UserRepository handles user data operations
type UserRepository struct {
	userQueries *query.UserQueries
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		userQueries: query.NewUserQueries(pool),
	}
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	row := r.userQueries.GetUserByUsername(ctx, username)

	user := &model.User{}
	err := row.Scan(
		&user.ID, &user.Username, &user.Password, &user.FullName,
		&user.Email, &user.Phone, &user.Role, &user.IsActive,
		&user.CounterID, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin,
	)
	if err != nil {
		log.Error().Err(err).Str("layer", "repository").Msg("Failed to get user by username")
		return nil, err
	}
	return user, nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id int) (*model.User, error) {
	row := r.userQueries.GetUserByID(ctx, id)

	user := &model.User{}
	err := row.Scan(
		&user.ID, &user.Username, &user.FullName, &user.Email,
		&user.Phone, &user.Role, &user.IsActive, &user.CounterID,
		&user.CreatedAt, &user.UpdatedAt, &user.LastLogin,
	)
	if err != nil {
		log.Error().Err(err).Str("layer", "repository").Msg("Failed to get user by id")
		return nil, err
	}
	return user, nil
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	id, createdAt, updatedAt, err := r.userQueries.CreateUser(
		ctx,
		user.Username,
		user.Password,
		user.FullName.String,
		user.Email.String,
		user.Phone.String,
		user.Role,
		func() *int {
			if user.CounterID.Valid {
				id := int(user.CounterID.Int64)
				return &id
			}
			return nil
		}(),
		user.IsActive,
	)
	if err != nil {
		return nil, err
	}

	user.ID = id
	user.CreatedAt = createdAt
	user.UpdatedAt = updatedAt
	return user, nil
}

// Update updates user information
func (r *UserRepository) Update(ctx context.Context, user *model.User) (*model.User, error) {
	err := r.userQueries.UpdateUser(
		ctx,
		user.FullName.String,
		user.Email.String,
		user.Phone.String,
		user.Role,
		func() *int {
			if user.CounterID.Valid {
				id := int(user.CounterID.Int64)
				return &id
			}
			return nil
		}(),
		user.IsActive,
		user.ID,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Delete deletes a user
func (r *UserRepository) Delete(ctx context.Context, id int) error {
	return r.userQueries.DeleteUser(ctx, id)
}

// UpdatePassword updates user password
func (r *UserRepository) UpdatePassword(ctx context.Context, id int, password string) error {
	return r.userQueries.UpdateUserPassword(ctx, id, password)
}

// UpdateLastLogin updates the last login timestamp
func (r *UserRepository) UpdateLastLogin(ctx context.Context, id int) error {
	return r.userQueries.UpdateLastLogin(ctx, id)
}

// List retrieves users with optional role filter
func (r *UserRepository) List(ctx context.Context, role string) ([]model.User, error) {
	rows, err := r.userQueries.ListUsers(ctx, role)
	if err != nil {
		log.Error().Err(err).Str("layer", "repository").Str("method", "List").Str("domain", "user").Msg("Failed to list users")
		return nil, err
	}
	defer rows.Close()

	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.User])
	if err != nil {
		log.Error().Err(err).Str("layer", "repository").Str("method", "List").Str("domain", "user").Msg("Failed to collect rows")
		return nil, err
	}

	return res, nil
}
