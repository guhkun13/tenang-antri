package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"queue-system/internal/model"
	"queue-system/internal/query"
)

// StatsRepository handles statistics data operations
type StatsRepository struct {
	statsQueries *query.StatsQueries
}

func NewStatsRepository(pool *pgxpool.Pool) *StatsRepository {
	return &StatsRepository{
		statsQueries: query.NewStatsQueries(pool),
	}
}

// GetDashboardStats retrieves comprehensive dashboard statistics
func (r *StatsRepository) GetDashboardStats(ctx context.Context) (*model.DashboardStats, error) {
	stats := &model.DashboardStats{
		TicketsByStatus: make(map[string]int),
	}

	// Total tickets today
	if err := r.statsQueries.GetTotalTicketsToday(ctx).Scan(&stats.TotalTicketsToday); err != nil {
		return nil, err
	}

	// Currently serving
	if err := r.statsQueries.GetCurrentlyServingCount(ctx).Scan(&stats.CurrentlyServing); err != nil {
		return nil, err
	}

	// Waiting tickets
	if err := r.statsQueries.GetWaitingTicketsCount(ctx).Scan(&stats.WaitingTickets); err != nil {
		return nil, err
	}

	// Active counters
	if err := r.statsQueries.GetActiveCountersCount(ctx).Scan(&stats.ActiveCounters); err != nil {
		return nil, err
	}

	// Paused counters
	if err := r.statsQueries.GetPausedCountersCount(ctx).Scan(&stats.PausedCounters); err != nil {
		return nil, err
	}

	// Average wait time today
	if err := r.statsQueries.GetAvgWaitTimeToday(ctx).Scan(&stats.AvgWaitTime); err != nil {
		return nil, err
	}

	// Average service time today
	if err := r.statsQueries.GetAvgServiceTimeToday(ctx).Scan(&stats.AvgServiceTime); err != nil {
		return nil, err
	}

	// Tickets by status today
	rows, err := r.statsQueries.GetTicketsByStatusToday(ctx)
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

	// Queue length by category
	stats.QueueLengthByCategory, err = r.GetQueueLengthByCategory(ctx)
	if err != nil {
		return nil, err
	}

	// Hourly distribution
	stats.HourlyDistribution, err = r.GetHourlyDistribution(ctx)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

// GetQueueLengthByCategory retrieves queue length for each category
func (r *StatsRepository) GetQueueLengthByCategory(ctx context.Context) ([]model.CategoryQueueStats, error) {
	rows, err := r.statsQueries.GetQueueLengthByCategory(ctx)
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

// GetQueueLengthByCategories retrieves queue length for specific categories
func (r *StatsRepository) GetQueueLengthByCategories(ctx context.Context, categoryIDs []int) ([]model.CategoryQueueStats, error) {
	if len(categoryIDs) == 0 {
		return []model.CategoryQueueStats{}, nil
	}

	rows, err := r.statsQueries.GetQueueLengthByCategories(ctx, categoryIDs)
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

// GetHourlyDistribution retrieves hourly ticket distribution
func (r *StatsRepository) GetHourlyDistribution(ctx context.Context) ([]model.HourlyStats, error) {
	rows, err := r.statsQueries.GetHourlyDistribution(ctx)
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

// GetCurrentlyServingTickets retrieves currently serving tickets for display
func (r *StatsRepository) GetCurrentlyServingTickets(ctx context.Context) ([]model.DisplayTicket, error) {
	rows, err := r.statsQueries.GetCurrentlyServingTickets(ctx)
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
