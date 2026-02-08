package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"

	"tenangantri/internal/model"
	"tenangantri/internal/query"
)

type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByUsernameWithPassword(ctx context.Context, username string) (*UserWithPassword, error)
	GetByIDWithPassword(ctx context.Context, id int) (*UserWithPassword, error)
	GetByID(ctx context.Context, id int) (*model.User, error)
	Create(ctx context.Context, user *model.User) (*model.User, error)
	Update(ctx context.Context, user *model.User) (*model.User, error)
	Delete(ctx context.Context, id int) error
	UpdatePassword(ctx context.Context, id int, password string) error
	UpdateLastLogin(ctx context.Context, id int) error
	List(ctx context.Context, role string) ([]model.User, error)
}

type userRepository struct {
	pool    DB
	userQry *query.UserQueries
}

func NewUserRepository(pool DB) UserRepository {
	return &userRepository{
		pool:    pool,
		userQry: query.NewUserQueries(),
	}
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	queryStr := r.userQry.GetUserByUsername(ctx)
	row := r.pool.QueryRow(ctx, queryStr, username)

	user := &model.User{}
	err := row.Scan(
		&user.ID, &user.Username, &user.FullName, &user.Email, &user.Phone,
		&user.Role, &user.IsActive, &user.CounterID, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pgx.ErrNoRows
		}
		log.Error().Err(err).Str("layer", "repository").Str("username", username).Msg("Failed to scan user")
		return nil, err
	}
	return user, nil
}

type UserWithPassword struct {
	ID       int    `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
}

func (r *userRepository) GetByUsernameWithPassword(ctx context.Context, username string) (*UserWithPassword, error) {
	sql := r.userQry.GetUserPasswordByUsername(ctx)
	rows, err := r.pool.Query(ctx, sql, username)
	if err != nil {
		log.Error().Err(err).Str("layer", "repository").Msg("Failed to query user by username with password")
		return nil, err
	}
	defer rows.Close()

	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[UserWithPassword])
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, err
		}
		log.Error().Err(err).Str("layer", "repository").Str("function", "GetByUsernameWithPassword").Msg("Failed to collect user")
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByIDWithPassword(ctx context.Context, id int) (*UserWithPassword, error) {
	sql := `SELECT id, username, password FROM users WHERE id = $1`
	rows, err := r.pool.Query(ctx, sql, id)
	if err != nil {
		log.Error().Err(err).Str("layer", "repository").Msg("Failed to query user by id with password")
		return nil, err
	}
	defer rows.Close()

	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[UserWithPassword])
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, err
		}
		log.Error().Err(err).Str("layer", "repository").Str("function", "GetByIDWithPassword").Msg("Failed to collect user")
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*model.User, error) {
	queryStr := r.userQry.GetUserByID(ctx)
	row := r.pool.QueryRow(ctx, queryStr, id)

	user := &model.User{}
	err := row.Scan(
		&user.ID, &user.Username, &user.FullName, &user.Email, &user.Phone,
		&user.Role, &user.IsActive, &user.CounterID, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		log.Error().Err(err).Str("layer", "repository").Int("id", id).Msg("Failed to scan user")
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
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

func (r *userRepository) Update(ctx context.Context, user *model.User) (*model.User, error) {
	sql := r.userQry.UpdateUser(ctx)
	_, err := r.pool.Exec(ctx, sql, user.FullName.String, user.Email.String, user.Phone.String, user.Role, user.CounterID, user.IsActive, user.ID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Delete(ctx context.Context, id int) error {
	sql := r.userQry.DeleteUser(ctx)
	_, err := r.pool.Exec(ctx, sql, id)
	return err
}

func (r *userRepository) UpdatePassword(ctx context.Context, id int, password string) error {
	sql := r.userQry.UpdateUserPassword(ctx)
	_, err := r.pool.Exec(ctx, sql, password, id)
	return err
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, id int) error {
	sql := r.userQry.UpdateLastLogin(ctx)
	_, err := r.pool.Exec(ctx, sql, id)
	return err
}

func (r *userRepository) List(ctx context.Context, role string) ([]model.User, error) {
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
