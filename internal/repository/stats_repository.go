package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"tenangantri/internal/model"
	"tenangantri/internal/query"
)

type StatsRepository struct {
	pool     *pgxpool.Pool
	statsQry *query.StatsQueries
}

func NewStatsRepository(pool *pgxpool.Pool) *StatsRepository {
	return &StatsRepository{
		pool:     pool,
		statsQry: query.NewStatsQueries(),
	}
}

func (r *StatsRepository) GetDashboardStats(ctx context.Context) (*model.DashboardStats, error) {
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

	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		stats.TicketsByStatus[status] = count
	}

	stats.QueueLengthByCategory, err = r.GetQueueLengthByCategory(ctx)
	if err != nil {
		return nil, err
	}

	stats.HourlyDistribution, err = r.GetHourlyDistribution(ctx)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (r *StatsRepository) GetQueueLengthByCategory(ctx context.Context) ([]model.CategoryQueueStats, error) {
	sql := r.statsQry.GetQueueLengthByCategory(ctx)
	rows, err := r.pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.CategoryQueueStats
	for rows.Next() {
		var stats model.CategoryQueueStats
		if err := rows.Scan(&stats.CategoryID, &stats.CategoryName, &stats.Prefix, &stats.ColorCode, &stats.WaitingCount); err != nil {
			return nil, err
		}
		result = append(result, stats)
	}

	return result, nil
}

func (r *StatsRepository) GetQueueLengthByCategories(ctx context.Context, categoryIDs []int) ([]model.CategoryQueueStats, error) {
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

	var result []model.CategoryQueueStats
	for rows.Next() {
		var stats model.CategoryQueueStats
		err := rows.Scan(&stats.CategoryID, &stats.CategoryName, &stats.Prefix, &stats.ColorCode, &stats.WaitingCount)
		if err != nil {
			log.Error().Err(err).Str("layer", "repository").Str("func", "GetQueueLengthByCategories").Msg("Failed to scan row")
			return nil, err
		}
		result = append(result, stats)
	}

	return result, nil
}

func (r *StatsRepository) GetHourlyDistribution(ctx context.Context) ([]model.HourlyStats, error) {
	sql := r.statsQry.GetHourlyDistribution(ctx)
	rows, err := r.pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.HourlyStats
	for rows.Next() {
		var stats model.HourlyStats
		if err := rows.Scan(&stats.Hour, &stats.Count); err != nil {
			return nil, err
		}
		result = append(result, stats)
	}

	return result, nil
}

func (r *StatsRepository) GetCurrentlyServingTickets(ctx context.Context) ([]model.DisplayTicket, error) {
	sql := r.statsQry.GetCurrentlyServingTickets(ctx)
	rows, err := r.pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.DisplayTicket
	for rows.Next() {
		var dt model.DisplayTicket
		if err := rows.Scan(&dt.TicketNumber, &dt.CounterNumber, &dt.CategoryPrefix, &dt.ColorCode, &dt.Status); err != nil {
			return nil, err
		}
		result = append(result, dt)
	}

	return result, nil
}
