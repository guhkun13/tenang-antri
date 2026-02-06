package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"

	"tenangantri/internal/model"
	"tenangantri/internal/query"
)

type StatsRepository interface {
	GetDashboardStats(ctx context.Context) (*model.DashboardStats, error)
	GetQueueLengthByCategory(ctx context.Context) ([]model.CategoryQueueStats, error)
	GetQueueLengthByCategories(ctx context.Context, categoryIDs []int) ([]model.CategoryQueueStats, error)
	GetHourlyDistribution(ctx context.Context) ([]model.HourlyStats, error)
	GetCurrentlyServingTickets(ctx context.Context) ([]model.DisplayTicket, error)
}

type statsRepository struct {
	pool     DB
	statsQry *query.StatsQueries
}

func NewStatsRepository(pool DB) StatsRepository {
	return &statsRepository{
		pool:     pool,
		statsQry: query.NewStatsQueries(),
	}
}

func (r *statsRepository) GetDashboardStats(ctx context.Context) (*model.DashboardStats, error) {
	stats := &model.DashboardStats{
		TicketsByStatus: make(map[string]int),
	}

	sql := r.statsQry.GetTotalTicketsToday(ctx)
	if err := r.pool.QueryRow(ctx, sql).Scan(&stats.TotalTicketsToday); err != nil {
		return nil, err
	}

	sql = r.statsQry.GetCurrentlyServingCount(ctx)
	if err := r.pool.QueryRow(ctx, sql).Scan(&stats.CurrentlyServing); err != nil {
		return nil, err
	}

	sql = r.statsQry.GetWaitingTicketsCount(ctx)
	if err := r.pool.QueryRow(ctx, sql).Scan(&stats.WaitingTickets); err != nil {
		return nil, err
	}

	sql = r.statsQry.GetActiveCountersCount(ctx)
	if err := r.pool.QueryRow(ctx, sql).Scan(&stats.ActiveCounters); err != nil {
		return nil, err
	}

	sql = r.statsQry.GetPausedCountersCount(ctx)
	if err := r.pool.QueryRow(ctx, sql).Scan(&stats.PausedCounters); err != nil {
		return nil, err
	}

	sql = r.statsQry.GetAvgWaitTimeToday(ctx)
	if err := r.pool.QueryRow(ctx, sql).Scan(&stats.AvgWaitTime); err != nil {
		return nil, err
	}

	sql = r.statsQry.GetAvgServiceTimeToday(ctx)
	if err := r.pool.QueryRow(ctx, sql).Scan(&stats.AvgServiceTime); err != nil {
		return nil, err
	}

	sql = r.statsQry.GetTicketsByStatusToday(ctx)
	rows, err := r.pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type statusCount struct {
		Status string
		Count  int
	}

	statusRows, err := pgx.CollectRows(rows, pgx.RowToStructByPos[statusCount])
	if err != nil {
		return nil, err
	}

	for _, sr := range statusRows {
		stats.TicketsByStatus[sr.Status] = sr.Count
	}

	stats.QueueLengthByCategory, err = r.GetQueueLengthByCategory(ctx)
	if err != nil {
		stats.QueueLengthByCategory = []model.CategoryQueueStats{}
	}

	stats.HourlyDistribution, err = r.GetHourlyDistribution(ctx)
	if err != nil {
		stats.HourlyDistribution = []model.HourlyStats{}
	}

	return stats, nil
}

func (r *statsRepository) GetQueueLengthByCategory(ctx context.Context) ([]model.CategoryQueueStats, error) {
	sql := r.statsQry.GetQueueLengthByCategory(ctx)
	rows, err := r.pool.Query(ctx, sql)
	if err != nil {
		return []model.CategoryQueueStats{}, err
	}
	defer rows.Close()

	ptrResult, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByPos[model.CategoryQueueStats])
	if err != nil {
		return []model.CategoryQueueStats{}, err
	}

	result := make([]model.CategoryQueueStats, len(ptrResult))
	for i, ptr := range ptrResult {
		result[i] = *ptr
	}

	return result, nil
}

func (r *statsRepository) GetQueueLengthByCategories(ctx context.Context, categoryIDs []int) ([]model.CategoryQueueStats, error) {
	if len(categoryIDs) == 0 {
		return []model.CategoryQueueStats{}, nil
	}

	sql := r.statsQry.GetQueueLengthByCategories(ctx, categoryIDs)
	args := make([]any, len(categoryIDs))
	for i, id := range categoryIDs {
		args[i] = id
	}
	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ptrResult, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByPos[model.CategoryQueueStats])
	if err != nil {
		log.Error().Err(err).Str("layer", "repository").Str("func", "GetQueueLengthByCategories").Msg("Failed to collect rows")
		return nil, err
	}

	result := make([]model.CategoryQueueStats, len(ptrResult))
	for i, ptr := range ptrResult {
		result[i] = *ptr
	}

	return result, nil
}

func (r *statsRepository) GetHourlyDistribution(ctx context.Context) ([]model.HourlyStats, error) {
	sql := r.statsQry.GetHourlyDistribution(ctx)
	rows, err := r.pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ptrResult, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByPos[model.HourlyStats])
	if err != nil {
		return nil, err
	}

	result := make([]model.HourlyStats, len(ptrResult))
	for i, ptr := range ptrResult {
		result[i] = *ptr
	}

	return result, nil
}

func (r *statsRepository) GetCurrentlyServingTickets(ctx context.Context) ([]model.DisplayTicket, error) {
	sql := r.statsQry.GetCurrentlyServingTickets(ctx)
	rows, err := r.pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ptrResult, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByPos[model.DisplayTicket])
	if err != nil {
		return nil, err
	}

	result := make([]model.DisplayTicket, len(ptrResult))
	for i, ptr := range ptrResult {
		result[i] = *ptr
	}

	return result, nil
}
