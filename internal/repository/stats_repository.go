package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"tenangantri/internal/dto"
	"tenangantri/internal/query"
)

type StatsRepository interface {
	GetDashboardStats(ctx context.Context) (*dto.DashboardStats, error)
	GetQueueLengthByCategory(ctx context.Context) ([]dto.CategoryQueueStats, error)
	GetQueueLengthByCategories(ctx context.Context, categoryIDs []int) ([]dto.CategoryQueueStats, error)
	GetHourlyDistribution(ctx context.Context) ([]dto.HourlyStats, error)
	GetCurrentlyServingTickets(ctx context.Context) ([]dto.DisplayTicket, error)
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

func (r *statsRepository) GetDashboardStats(ctx context.Context) (*dto.DashboardStats, error) {
	stats := &dto.DashboardStats{
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
		stats.QueueLengthByCategory = []dto.CategoryQueueStats{}
	}

	stats.HourlyDistribution, err = r.GetHourlyDistribution(ctx)
	if err != nil {
		stats.HourlyDistribution = []dto.HourlyStats{}
	}

	return stats, nil
}

func (r *statsRepository) GetQueueLengthByCategory(ctx context.Context) ([]dto.CategoryQueueStats, error) {
	sql := r.statsQry.GetQueueLengthByCategory(ctx)
	rows, err := r.pool.Query(ctx, sql)
	if err != nil {
		return []dto.CategoryQueueStats{}, err
	}
	defer rows.Close()

	var results []dto.CategoryQueueStats
	for rows.Next() {
		var stats dto.CategoryQueueStats
		var lastTicketNumber, counterNumber pgtype.Text
		err := rows.Scan(
			&stats.CategoryID, &stats.CategoryName, &stats.Prefix, &stats.ColorCode, &stats.WaitingCount,
			&lastTicketNumber, &counterNumber,
		)
		if err != nil {
			return []dto.CategoryQueueStats{}, err
		}
		if lastTicketNumber.Valid {
			stats.LastTicketNumber = lastTicketNumber.String
		}
		if counterNumber.Valid {
			stats.CounterNumber = counterNumber.String
		}
		results = append(results, stats)
	}

	return results, nil
}

func (r *statsRepository) GetQueueLengthByCategories(ctx context.Context, categoryIDs []int) ([]dto.CategoryQueueStats, error) {
	if len(categoryIDs) == 0 {
		return []dto.CategoryQueueStats{}, nil
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

	var results []dto.CategoryQueueStats
	for rows.Next() {
		var stats dto.CategoryQueueStats
		var lastTicketNumber, counterNumber pgtype.Text
		err := rows.Scan(
			&stats.CategoryID, &stats.CategoryName, &stats.Prefix, &stats.ColorCode, &stats.WaitingCount,
			&lastTicketNumber, &counterNumber,
		)
		if err != nil {
			return nil, err
		}
		if lastTicketNumber.Valid {
			stats.LastTicketNumber = lastTicketNumber.String
		}
		if counterNumber.Valid {
			stats.CounterNumber = counterNumber.String
		}
		results = append(results, stats)
	}

	return results, nil
}

func (r *statsRepository) GetHourlyDistribution(ctx context.Context) ([]dto.HourlyStats, error) {
	sql := r.statsQry.GetHourlyDistribution(ctx)
	rows, err := r.pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ptrResult, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByPos[dto.HourlyStats])
	if err != nil {
		return nil, err
	}

	result := make([]dto.HourlyStats, len(ptrResult))
	for i, ptr := range ptrResult {
		result[i] = *ptr
	}

	return result, nil
}

func (r *statsRepository) GetCurrentlyServingTickets(ctx context.Context) ([]dto.DisplayTicket, error) {
	sql := r.statsQry.GetCurrentlyServingTickets(ctx)
	rows, err := r.pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ptrResult, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByPos[dto.DisplayTicket])
	if err != nil {
		return nil, err
	}

	result := make([]dto.DisplayTicket, len(ptrResult))
	for i, ptr := range ptrResult {
		result[i] = *ptr
	}

	return result, nil
}
