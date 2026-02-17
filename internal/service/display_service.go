package service

import (
	"context"

	"github.com/rs/zerolog/log"

	"tenangantri/internal/dto"
	"tenangantri/internal/model"
	"tenangantri/internal/repository"
)

// DisplayService handles display-related business logic
type DisplayService struct {
	statsRepo    repository.StatsRepository
	categoryRepo repository.CategoryRepository
	counterRepo  repository.CounterRepository
}

func NewDisplayService(statsRepo repository.StatsRepository, categoryRepo repository.CategoryRepository, counterRepo repository.CounterRepository) *DisplayService {
	return &DisplayService{
		statsRepo:    statsRepo,
		categoryRepo: categoryRepo,
		counterRepo:  counterRepo,
	}
}

// GetDisplayData gets data for the main display
func (s *DisplayService) GetDisplayData(ctx context.Context) ([]dto.DisplayTicket, []model.Category, []model.Counter, error) {
	tickets, err := s.statsRepo.GetCurrentlyServingTickets(ctx)
	if err != nil {
		log.Error().Err(err).Str("layer", "service").Str("func", "GetDisplayData").Msg("Failed to get currently serving tickets")
		tickets = []dto.DisplayTicket{}
	}

	categories, err := s.categoryRepo.List(ctx, true, true)
	if err != nil {
		log.Error().Err(err).Str("layer", "service").Str("func", "GetDisplayData").Msg("Failed to get categories")
		categories = []model.Category{}
	}

	counters, err := s.counterRepo.List(ctx)
	if err != nil {
		log.Error().Err(err).Str("layer", "service").Str("func", "GetDisplayData").Msg("Failed to get counters")
		counters = []model.Counter{}
	}

	return tickets, categories, counters, nil
}

// GetCurrentlyServing gets currently serving tickets
func (s *DisplayService) GetCurrentlyServing(ctx context.Context) ([]dto.DisplayTicket, error) {
	return s.statsRepo.GetCurrentlyServingTickets(ctx)
}

// GetQueueStats gets queue statistics for display
func (s *DisplayService) GetQueueStats(ctx context.Context) (*dto.DashboardStats, error) {
	return s.statsRepo.GetDashboardStats(ctx)
}

// GetWaitingByCategory gets waiting tickets count by category
func (s *DisplayService) GetWaitingByCategory(ctx context.Context) ([]dto.CategoryQueueStats, error) {
	return s.statsRepo.GetQueueLengthByCategory(ctx)
}

// GetCategoryDisplayData gets data for category-specific display
func (s *DisplayService) GetCategoryDisplayData(ctx context.Context, categoryID int) (*model.Category, []dto.DisplayTicket, error) {
	categories, err := s.categoryRepo.List(ctx, true, true)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get categories")
		return nil, nil, err
	}

	var selectedCategory *model.Category
	for _, cat := range categories {
		if cat.ID == categoryID {
			selectedCategory = &cat
			break
		}
	}

	tickets, err := s.statsRepo.GetCurrentlyServingTickets(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get currently serving tickets")
		return selectedCategory, []dto.DisplayTicket{}, nil
	}

	return selectedCategory, tickets, nil
}
