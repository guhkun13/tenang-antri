package service

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"

	"tenangantri/internal/dto"
	"tenangantri/internal/model"
	"tenangantri/internal/repository"
)

// TicketService handles ticket-related business logic
type TicketService struct {
	ticketRepo   repository.TicketRepository
	categoryRepo repository.CategoryRepository
	statsRepo    repository.StatsRepository
}

func NewTicketService(ticketRepo repository.TicketRepository, categoryRepo repository.CategoryRepository, statsRepo repository.StatsRepository) *TicketService {
	return &TicketService{
		ticketRepo:   ticketRepo,
		categoryRepo: categoryRepo,
		statsRepo:    statsRepo,
	}
}

// CreateTicket creates a new ticket
func (s *TicketService) CreateTicket(ctx context.Context, req *dto.CreateTicketRequest) (*model.Ticket, error) {
	// Get category to validate and get prefix
	category, err := s.categoryRepo.GetByID(ctx, req.CategoryID)
	if err != nil {
		return nil, err
	}

	// Generate ticket number
	ticketNumber, dailySequence, err := s.ticketRepo.GenerateNumber(ctx, category.ID, category.Prefix)
	if err != nil {
		return nil, err
	}

	ticket := &model.Ticket{
		TicketNumber:  ticketNumber,
		DailySequence: dailySequence,
		QueueDate:     time.Now(),
		Category:      &model.Category{ID: req.CategoryID},
		Status:        "waiting",
		Priority:      req.Priority,
	}

	createdTicket, err := s.ticketRepo.Create(ctx, ticket)
	if err != nil {
		return nil, err
	}

	// Get full details
	return s.ticketRepo.GetWithDetails(ctx, createdTicket.ID)
}

// GetTicket retrieves a ticket by ID
func (s *TicketService) GetTicket(ctx context.Context, id int) (*model.Ticket, error) {
	return s.ticketRepo.GetWithDetails(ctx, id)
}

// UpdateTicketStatus updates the status of a ticket
func (s *TicketService) UpdateTicketStatus(ctx context.Context, id int, req *dto.UpdateTicketStatusRequest) (*model.Ticket, error) {
	err := s.ticketRepo.UpdateStatus(ctx, id, req.Status)
	if err != nil {
		return nil, err
	}

	return s.ticketRepo.GetWithDetails(ctx, id)
}

// CancelTicket cancels a ticket
func (s *TicketService) CancelTicket(ctx context.Context, id int) error {
	return s.ticketRepo.UpdateStatus(ctx, id, "cancelled")
}

// GetNextTicket gets the next ticket for the given categories
func (s *TicketService) GetNextTicket(ctx context.Context, categoryIDs []int) (*model.Ticket, error) {
	return s.ticketRepo.GetNextTicket(ctx, categoryIDs)
}

// AssignTicketToCounter assigns a ticket to a counter
func (s *TicketService) AssignTicketToCounter(ctx context.Context, ticketID, counterID int) (*model.Ticket, error) {
	err := s.ticketRepo.AssignToCounter(ctx, ticketID, counterID)
	if err != nil {
		return nil, err
	}

	return s.ticketRepo.GetWithDetails(ctx, ticketID)
}

// GetCurrentTicketForCounter gets the current ticket being served at a counter
func (s *TicketService) GetCurrentTicketForCounter(ctx context.Context, counterID int) (*model.Ticket, error) {
	ticket, err := s.ticketRepo.GetCurrentForCounter(ctx, counterID)
	if err != nil {
		return nil, err
	}

	if ticket == nil {
		return nil, nil
	}

	// Load category details
	if ticket.Category != nil {
		category, err := s.categoryRepo.GetByID(ctx, ticket.Category.ID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to load ticket category")
		} else {
			ticket.Category = category
		}
	}

	return ticket, nil
}

// ListTickets retrieves tickets with optional filters
func (s *TicketService) ListTickets(ctx context.Context, filters map[string]interface{}) ([]model.Ticket, error) {
	return s.ticketRepo.List(ctx, filters)
}

// GetWaitingTicketsPreview gets a preview of waiting tickets
func (s *TicketService) GetWaitingTicketsPreview(ctx context.Context, limit int) ([]model.Ticket, error) {
	return s.ticketRepo.GetWaitingPreview(ctx, limit)
}

// GetWaitingTicketsPreviewByCategories gets waiting tickets preview for specific categories
func (s *TicketService) GetWaitingTicketsPreviewByCategories(ctx context.Context, categoryIDs []int, limit int) ([]model.Ticket, error) {
	return s.ticketRepo.GetWaitingPreviewByCategories(ctx, categoryIDs, limit)
}

// GetTodayCompletedTicketsByCategories gets completed tickets today for specific categories
func (s *TicketService) GetTodayCompletedTicketsByCategories(ctx context.Context, categoryIDs []int) ([]model.Ticket, error) {
	return s.ticketRepo.GetTodayCompletedByCategories(ctx, categoryIDs)
}

// GetTodayTicketCount gets today's ticket count
func (s *TicketService) GetTodayTicketCount(ctx context.Context) (int, error) {
	return s.ticketRepo.GetTodayCount(ctx)
}

// GetTodayTicketCountByCategory gets today's ticket count for a specific category
func (s *TicketService) GetTodayTicketCountByCategory(ctx context.Context, categoryID int) (int, error) {
	return s.ticketRepo.GetTodayCountByCategory(ctx, categoryID)
}
