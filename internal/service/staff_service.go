package service

import (
	"context"

	"github.com/rs/zerolog/log"

	"tenangantri/internal/dto"
	"tenangantri/internal/model"
	"tenangantri/internal/repository"
)

// StaffService handles staff-specific business logic
type StaffService struct {
	userRepo     repository.UserRepository
	counterRepo  repository.CounterRepository
	ticketRepo   repository.TicketRepository
	statsRepo    repository.StatsRepository
	categoryRepo repository.CategoryRepository
}

func NewStaffService(userRepo repository.UserRepository,
	counterRepo repository.CounterRepository,
	ticketRepo repository.TicketRepository,
	statsRepo repository.StatsRepository,
	categoryRepo repository.CategoryRepository) *StaffService {
	return &StaffService{
		userRepo:     userRepo,
		counterRepo:  counterRepo,
		ticketRepo:   ticketRepo,
		statsRepo:    statsRepo,
		categoryRepo: categoryRepo,
	}
}

// GetDashboardData gets staff dashboard data
func (s *StaffService) GetDashboardData(ctx context.Context, userID int) (*dto.StaffDashboardResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		log.Error().Err(err).Str("layer", "service").Str("func", "GetDashboardData").Msg("Failed to load user")
		return nil, err
	}

	if user.CounterID.Int64 == 0 {
		return nil, nil
	}

	counter, err := s.counterRepo.GetByID(ctx, int(user.CounterID.Int64))
	if err != nil {
		log.Error().Err(err).Str("layer", "service").Str("func", "GetDashboardData").Msg("Failed to load counter")
		return nil, err
	}

	if counter.CategoryID != nil {
		category, err := s.categoryRepo.GetByID(ctx, *counter.CategoryID)
		if err != nil {
			log.Error().Err(err).Str("layer", "service").Str("func", "GetDashboardData").Msg("Failed to load counter category")
			return nil, err
		}
		counter.Category = category
	}

	// Get current serving ticket
	currentTicket, err := s.ticketRepo.GetCurrentForCounter(ctx, counter.ID)
	if err != nil {
		log.Error().Err(err).Str("layer", "service").Str("func", "GetDashboardData").Msg("Failed to load current ticket")
		return nil, err
	}
	log.Info().Interface("currentTicket", currentTicket).Msg("Current ticket")

	if currentTicket != nil && currentTicket.CategoryID != nil {
		currentCategory, err := s.categoryRepo.GetByID(ctx, *currentTicket.CategoryID)
		if err != nil {
			log.Error().Err(err).Str("layer", "service").Str("func", "GetDashboardData").Msg("Failed to load current ticket category")
			return nil, err
		}
		currentTicket.Category = currentCategory
	}

	log.Info().Interface("currentTicket", currentTicket).Msg("Current ticket")

	// Get waiting tickets preview
	var categoryIDs []int
	if counter.CategoryID != nil {
		categoryIDs = append(categoryIDs, *counter.CategoryID)
	}

	waitingTickets, err := s.ticketRepo.GetWaitingPreviewByCategories(ctx, categoryIDs, 5)
	if err != nil {
		log.Error().Err(err).Str("layer", "service").Str("func", "GetDashboardData").Msg("Failed to load waiting tickets preview")
		return nil, err
	}

	log.Info().Interface("waitingTickets", waitingTickets).Msg("Waiting tickets preview")

	// Get queue stats for counter's categories only
	queueStats, err := s.statsRepo.GetQueueLengthByCategories(ctx, categoryIDs)
	if err != nil {
		log.Error().Err(err).Str("layer", "service").Str("func", "GetDashboardData").Msg("Failed to load queue stats")
		return nil, err
	}

	// Get today's completed tickets
	completedTickets, err := s.ticketRepo.GetTodayCompletedByCategories(ctx, categoryIDs)
	if err != nil {
		log.Error().Err(err).Str("layer", "service").Str("func", "GetDashboardData").Msg("Failed to load completed tickets")
		return nil, err
	}

	response := &dto.StaffDashboardResponse{
		User:             user,
		Counter:          counter,
		CurrentTicket:    currentTicket,
		WaitingTickets:   waitingTickets,
		QueueStats:       queueStats,
		CompletedTickets: completedTickets,
		CategoryIDs:      categoryIDs,
	}

	return response, nil
}

// CallNext calls the next ticket for a staff member
func (s *StaffService) CallNext(ctx context.Context, userID int) (*model.Ticket, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.CounterID.Int64 == 0 {
		return nil, nil // No counter assigned
	}

	counter, err := s.counterRepo.GetByID(ctx, int(user.CounterID.Int64))
	if err != nil {
		return nil, err
	}

	// Check if counter is offline
	if counter.Status == model.CounterStatusOffline {
		return nil, nil // Counter is offline
	}

	// Business Rule: One ticket at a time - must complete current before calling next
	currentTicket, _ := s.ticketRepo.GetCurrentForCounter(ctx, counter.ID)
	if currentTicket != nil {
		return nil, nil // Already serving a ticket
	}

	// Check if counter is paused
	if counter.Status == model.CounterStatusPaused {
		return nil, nil // Counter is paused
	}

	// Get category for this counter
	var categoryIDs []int
	if counter.CategoryID != nil {
		categoryIDs = append(categoryIDs, *counter.CategoryID)
	}

	if len(categoryIDs) == 0 {
		return nil, nil // No category assigned
	}

	// Update counter status to serving
	_ = s.counterRepo.UpdateStatus(ctx, counter.ID, model.CounterStatusServing)

	// Get next ticket
	nextTicket, err := s.ticketRepo.GetNextTicket(ctx, categoryIDs)
	if err != nil {
		return nil, err
	}

	if nextTicket == nil {
		// No tickets in queue, set back to idle
		_ = s.counterRepo.UpdateStatus(ctx, counter.ID, model.CounterStatusIdle)
		return nil, nil
	}

	// Assign ticket to counter
	err = s.ticketRepo.AssignToCounter(ctx, nextTicket.ID, counter.ID)
	if err != nil {
		return nil, err
	}

	// Get full ticket details
	return s.ticketRepo.GetWithDetails(ctx, nextTicket.ID)
}

// CompleteTicket completes the current ticket and sets counter to IDLE
func (s *StaffService) CompleteTicket(ctx context.Context, userID int) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.CounterID.Int64 == 0 {
		return nil // No counter assigned
	}

	counterID := int(user.CounterID.Int64)

	currentTicket, err := s.ticketRepo.GetCurrentForCounter(ctx, counterID)
	if err != nil {
		return err
	}

	if currentTicket == nil {
		return nil // No ticket being served
	}

	// Update ticket to completed
	err = s.ticketRepo.UpdateStatus(ctx, currentTicket.ID, "completed")
	if err != nil {
		return err
	}

	// Set counter back to IDLE
	return s.counterRepo.UpdateStatus(ctx, counterID, model.CounterStatusIdle)
}

// MarkNoShow marks current ticket as no-show
func (s *StaffService) MarkNoShow(ctx context.Context, userID int) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.CounterID.Int64 == 0 {
		return nil // No counter assigned
	}

	currentTicket, err := s.ticketRepo.GetCurrentForCounter(ctx, int(user.CounterID.Int64))
	if err != nil {
		return err
	}

	if currentTicket == nil {
		return nil // No ticket being served
	}

	return s.ticketRepo.UpdateStatus(ctx, currentTicket.ID, "no_show")
}

// PauseCounter pauses the counter (staff on break)
func (s *StaffService) PauseCounter(ctx context.Context, userID int) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.CounterID.Int64 == 0 {
		return nil // No counter assigned
	}

	return s.counterRepo.UpdateStatus(ctx, int(user.CounterID.Int64), model.CounterStatusPaused)
}

// ResumeCounter resumes the counter from paused
func (s *StaffService) ResumeCounter(ctx context.Context, userID int) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.CounterID.Int64 == 0 {
		return nil // No counter assigned
	}

	return s.counterRepo.UpdateStatus(ctx, int(user.CounterID.Int64), model.CounterStatusIdle)
}

// GetQueueStatus gets queue status for staff
func (s *StaffService) GetQueueStatus(ctx context.Context, userID int) (*dto.StaffQueueStatusResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.CounterID.Int64 == 0 {
		return nil, nil
	}

	counter, err := s.counterRepo.GetByID(ctx, int(user.CounterID.Int64))
	if err != nil {
		return nil, err
	}

	// Get waiting tickets
	waitingTickets, err := s.ticketRepo.GetWaitingPreview(ctx, 10)
	if err != nil {
		return nil, err
	}

	// Get current ticket
	currentTicket, err := s.ticketRepo.GetCurrentForCounter(ctx, counter.ID)
	if err != nil {
		return nil, err
	}

	// Get queue stats
	queueStats, err := s.statsRepo.GetQueueLengthByCategory(ctx)
	if err != nil {
		return nil, err
	}

	response := &dto.StaffQueueStatusResponse{
		Counter:        counter,
		CurrentTicket:  currentTicket,
		WaitingTickets: waitingTickets,
		QueueStats:     queueStats,
	}

	return response, nil
}

// GetCurrentTicket gets the current ticket for staff
func (s *StaffService) GetCurrentTicket(ctx context.Context, userID int) (*model.Ticket, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.CounterID.Int64 == 0 {
		return nil, nil // No counter assigned
	}

	return s.ticketRepo.GetCurrentForCounter(ctx, int(user.CounterID.Int64))
}

// TransferTicket transfers a ticket to another counter
func (s *StaffService) TransferTicket(ctx context.Context, ticketID, counterID int) (*model.Ticket, error) {
	err := s.ticketRepo.AssignToCounter(ctx, ticketID, counterID)
	if err != nil {
		return nil, err
	}

	return s.ticketRepo.GetWithDetails(ctx, ticketID)
}
