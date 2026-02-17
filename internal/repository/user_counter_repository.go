package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"

	"tenangantri/internal/model"
	"tenangantri/internal/query"
)

type UserCounterRepository interface {
	Create(ctx context.Context, userID, counterID int) (*model.UserCounter, error)
	GetByUserID(ctx context.Context, userID int) (*model.UserCounter, error)
	GetCounterIDByUserID(ctx context.Context, userID int) (sql.NullInt64, error)
	GetByCounterID(ctx context.Context, counterID int) (*model.UserCounter, error)
	DeleteByUserID(ctx context.Context, userID int) error
	DeleteByCounterID(ctx context.Context, counterID int) error
	ListAll(ctx context.Context) ([]model.UserCounter, error)
}

type userCounterRepository struct {
	pool DB
	qry  *query.UserCounterQueries
}

func NewUserCounterRepository(pool DB) UserCounterRepository {
	return &userCounterRepository{
		pool: pool,
		qry:  query.NewUserCounterQueries(),
	}
}

func (r *userCounterRepository) Create(ctx context.Context, userID, counterID int) (*model.UserCounter, error) {
	queryStr := r.qry.Create(ctx)
	var id int
	var createdAt, updatedAt time.Time
	err := r.pool.QueryRow(ctx, queryStr, userID, counterID).Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		log.Error().Err(err).Str("layer", "repository").Int("user_id", userID).Int("counter_id", counterID).Msg("Failed to create user_counter")
		return nil, err
	}

	return &model.UserCounter{
		ID:        id,
		UserID:    userID,
		CounterID: counterID,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func (r *userCounterRepository) GetByUserID(ctx context.Context, userID int) (*model.UserCounter, error) {
	queryStr := r.qry.GetByUserID(ctx)
	row := r.pool.QueryRow(ctx, queryStr, userID)

	uc := &model.UserCounter{}
	err := row.Scan(&uc.ID, &uc.UserID, &uc.CounterID, &uc.CreatedAt, &uc.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		log.Error().Err(err).Str("layer", "repository").Int("user_id", userID).Msg("Failed to get user_counter by user_id")
		return nil, err
	}
	return uc, nil
}

func (r *userCounterRepository) GetCounterIDByUserID(ctx context.Context, userID int) (sql.NullInt64, error) {
	queryStr := r.qry.GetCounterIDByUserID(ctx)
	row := r.pool.QueryRow(ctx, queryStr, userID)

	var counterID sql.NullInt64
	err := row.Scan(&counterID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return sql.NullInt64{}, nil
		}
		log.Error().Err(err).Str("layer", "repository").Int("user_id", userID).Msg("Failed to get counter_id by user_id")
		return sql.NullInt64{}, err
	}
	return counterID, nil
}

func (r *userCounterRepository) GetByCounterID(ctx context.Context, counterID int) (*model.UserCounter, error) {
	queryStr := r.qry.GetByCounterID(ctx)
	row := r.pool.QueryRow(ctx, queryStr, counterID)

	uc := &model.UserCounter{}
	err := row.Scan(&uc.ID, &uc.UserID, &uc.CounterID, &uc.CreatedAt, &uc.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		log.Error().Err(err).Str("layer", "repository").Int("counter_id", counterID).Msg("Failed to get user_counter by counter_id")
		return nil, err
	}
	return uc, nil
}

func (r *userCounterRepository) DeleteByUserID(ctx context.Context, userID int) error {
	queryStr := r.qry.DeleteByUserID(ctx)
	_, err := r.pool.Exec(ctx, queryStr, userID)
	if err != nil {
		log.Error().Err(err).Str("layer", "repository").Int("user_id", userID).Msg("Failed to delete user_counter by user_id")
		return err
	}
	return nil
}

func (r *userCounterRepository) DeleteByCounterID(ctx context.Context, counterID int) error {
	queryStr := r.qry.DeleteByCounterID(ctx)
	_, err := r.pool.Exec(ctx, queryStr, counterID)
	if err != nil {
		log.Error().Err(err).Str("layer", "repository").Int("counter_id", counterID).Msg("Failed to delete user_counter by counter_id")
		return err
	}
	return nil
}

func (r *userCounterRepository) ListAll(ctx context.Context) ([]model.UserCounter, error) {
	queryStr := r.qry.ListAll(ctx)
	rows, err := r.pool.Query(ctx, queryStr)
	if err != nil {
		log.Error().Err(err).Str("layer", "repository").Msg("Failed to list user_counters")
		return nil, err
	}
	defer rows.Close()

	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.UserCounter])
	if err != nil {
		log.Error().Err(err).Str("layer", "repository").Msg("Failed to collect user_counters rows")
		return nil, err
	}

	return res, nil
}
