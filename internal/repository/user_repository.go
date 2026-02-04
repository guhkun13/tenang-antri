package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"tenangantri/internal/model"
	"tenangantri/internal/query"
)

type UserRepository struct {
	pool    *pgxpool.Pool
	userQry *query.UserQueries
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool:    pool,
		userQry: query.NewUserQueries(),
	}
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	sql := r.userQry.GetUserByUsername(ctx)
	row := r.pool.QueryRow(ctx, sql, username)

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

func (r *UserRepository) GetByID(ctx context.Context, id int) (*model.User, error) {
	sql := r.userQry.GetUserByID(ctx)
	row := r.pool.QueryRow(ctx, sql, id)

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

func (r *UserRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	sql := r.userQry.CreateUser(ctx)
	var id int
	var createdAt, updatedAt time.Time
	err := r.pool.QueryRow(ctx, sql, user.Username, user.Password, user.FullName.String, user.Email.String, user.Phone.String, user.Role, user.CounterID, user.IsActive).Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	user.ID = id
	user.CreatedAt = createdAt
	user.UpdatedAt = updatedAt
	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) (*model.User, error) {
	sql := r.userQry.UpdateUser(ctx)
	_, err := r.pool.Exec(ctx, sql, user.FullName.String, user.Email.String, user.Phone.String, user.Role, user.CounterID, user.IsActive, user.ID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Delete(ctx context.Context, id int) error {
	sql := r.userQry.DeleteUser(ctx)
	_, err := r.pool.Exec(ctx, sql, id)
	return err
}

func (r *UserRepository) UpdatePassword(ctx context.Context, id int, password string) error {
	sql := r.userQry.UpdateUserPassword(ctx)
	_, err := r.pool.Exec(ctx, sql, password, id)
	return err
}

func (r *UserRepository) UpdateLastLogin(ctx context.Context, id int) error {
	sql := r.userQry.UpdateLastLogin(ctx)
	_, err := r.pool.Exec(ctx, sql, id)
	return err
}

func (r *UserRepository) List(ctx context.Context, role string) ([]model.User, error) {
	sql := r.userQry.ListUsers(ctx, role)
	rows, err := r.pool.Query(ctx, sql)
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
