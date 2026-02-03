package service

import (
	"context"

	"github.com/rs/zerolog/log"

	"queue-system/internal/dto"
	"queue-system/internal/model"
	"queue-system/internal/repository"
)

// KioskService handles kiosk-related business logic
type KioskService struct {
	categoryRepo *repository.CategoryRepository
	ticketRepo   *repository.TicketRepository
	statsRepo    *repository.StatsRepository
}

func NewKioskService(categoryRepo *repository.CategoryRepository, ticketRepo *repository.TicketRepository, statsRepo *repository.StatsRepository) *KioskService {
	return &KioskService{
		categoryRepo: categoryRepo,
		ticketRepo:   ticketRepo,
		statsRepo:    statsRepo,
	}
}

// GetCategories gets active categories for kiosk
func (s *KioskService) GetCategories(ctx context.Context) ([]model.Category, error) {
	categories, err := s.categoryRepo.List(ctx, true)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list categories")
		return []model.Category{}, err
	}

	return categories, nil
}

// GenerateTicket generates a new ticket from kiosk
func (s *KioskService) GenerateTicket(ctx context.Context, req *dto.CreateTicketRequest) (*model.Ticket, int, int, error) {
	// Get category to validate and get prefix
	category, err := s.categoryRepo.GetByID(ctx, req.CategoryID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get category by ID")
		return nil, 0, 0, err
	}

	// Generate ticket number
	ticketNumber, err := s.ticketRepo.GenerateNumber(ctx, category.Prefix)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate ticket number")
		return nil, 0, 0, err
	}

	ticket := &model.Ticket{
		TicketNumber: ticketNumber,
		CategoryID:   req.CategoryID,
		Status:       "waiting",
		Priority:     req.Priority,
	}

	createdTicket, err := s.ticketRepo.Create(ctx, ticket)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create ticket")
		return nil, 0, 0, err
	}

	// Get ticket details with category
	ticketWithDetails, err := s.ticketRepo.GetWithDetails(ctx, createdTicket.ID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get ticket details")
		return nil, 0, 0, err
	}

	// Get queue position
	waitingCount, err := s.ticketRepo.GetTodayCountByCategory(ctx, category.ID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get waiting count")
		waitingCount = 0
	}
	queuePosition := waitingCount

	// Get estimated wait time (simple calculation: avg wait time * position)
	stats, err := s.statsRepo.GetDashboardStats(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get stats")
		estimatedWaitTime := 0
		return ticketWithDetails, queuePosition, estimatedWaitTime, nil
	}

	estimatedWaitTime := stats.AvgWaitTime * queuePosition / 60 // in minutes

	return ticketWithDetails, queuePosition, estimatedWaitTime, nil
}

// GetQueueInfo gets queue information for kiosk display
func (s *KioskService) GetQueueInfo(ctx context.Context) (*model.DashboardStats, []model.Category, error) {
	stats, err := s.statsRepo.GetDashboardStats(ctx)
	if err != nil {
		return nil, nil, err
	}

	categories, err := s.categoryRepo.List(ctx, true)
	if err != nil {
		return nil, nil, err
	}

	return stats, categories, nil
}
