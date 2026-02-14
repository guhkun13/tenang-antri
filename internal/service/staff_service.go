package service

import (
	"context"
	"database/sql"

	"github.com/rs/zerolog/log"

	"tenangantri/internal/dto"
	"tenangantri/internal/model"
	"tenangantri/internal/repository"
)

// TicketListResult holds the paginated ticket list with metadata
type TicketListResult struct {
	Tickets    []model.Ticket
	Stats      map[string]int
	TotalCount int
}

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

	if counter.CategoryID.Valid {
		category, err := s.categoryRepo.GetByID(ctx, int(counter.CategoryID.Int64))
		if err != nil {
			log.Error().Err(err).Str("layer", "service").Str("func", "GetDashboardData").Msg("Failed to load counter category")
			return nil, err
		}
		counter.CategoryID = sql.NullInt64{Int64: int64(category.ID), Valid: true}
	}

	// Get current serving ticket
	currentTicket, err := s.ticketRepo.GetCurrentForCounter(ctx, counter.ID)
	if err != nil {
		log.Error().Err(err).Str("layer", "service").Str("func", "GetDashboardData").Msg("Failed to load current ticket")
		return nil, err
	}
	log.Info().Interface("currentTicket", currentTicket).Msg("Current ticket")

	if currentTicket != nil && currentTicket.CategoryID.Valid {
		currentCategory, err := s.categoryRepo.GetByID(ctx, int(currentTicket.CategoryID.Int64))
		if err != nil {
			log.Error().Err(err).Str("layer", "service").Str("func", "GetDashboardData").Msg("Failed to load current ticket category")
			return nil, err
		}
		currentTicket.Category = currentCategory
	}

	log.Info().Interface("currentTicket", currentTicket).Msg("Current ticket")

	// Get waiting tickets preview
	var categoryIDs []int
	if counter.CategoryID.Valid {
		categoryIDs = append(categoryIDs, int(counter.CategoryID.Int64))
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
	if counter.CategoryID.Valid {
		categoryIDs = append(categoryIDs, int(counter.CategoryID.Int64))
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

// GetTicketDetail gets detailed information about a ticket including timing metrics
func (s *StaffService) GetTicketDetail(ctx context.Context, ticketID int) (*model.Ticket, error) {
	ticket, err := s.ticketRepo.GetWithDetails(ctx, ticketID)
	if err != nil {
		return nil, err
	}
	return ticket, nil
}

// GetAllTickets gets all tickets for staff view based on their counter's categories with filters, pagination, and sorting
func (s *StaffService) GetAllTickets(ctx context.Context, userID int, filters map[string]interface{}) (*TicketListResult, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.CounterID.Int64 == 0 {
		return &TicketListResult{
			Tickets:    []model.Ticket{},
			Stats:      map[string]int{"Total": 0, "Waiting": 0, "Serving": 0, "Completed": 0},
			TotalCount: 0,
		}, nil
	}

	counter, err := s.counterRepo.GetByID(ctx, int(user.CounterID.Int64))
	if err != nil {
		return nil, err
	}

	var categoryIDs []int
	if counter.CategoryID.Valid {
		categoryIDs = append(categoryIDs, int(counter.CategoryID.Int64))
	}

	var tickets []model.Ticket
	var totalCount int

	if len(categoryIDs) > 0 {
		tickets, totalCount, err = s.ticketRepo.GetTicketsByCategoriesWithFilters(ctx, categoryIDs, filters)
		if err != nil {
			return nil, err
		}
	} else {
		tickets = []model.Ticket{}
		totalCount = 0
	}

	// Calculate stats from today's tickets
	stats := map[string]int{
		"Total":     0,
		"Waiting":   0,
		"Serving":   0,
		"Completed": 0,
	}

	// Get all today's tickets for stats (without pagination)
	if len(categoryIDs) > 0 {
		todayTickets, _ := s.ticketRepo.GetTodayByCategories(ctx, categoryIDs)
		stats["Total"] = len(todayTickets)
		for _, t := range todayTickets {
			switch t.Status {
			case "waiting":
				stats["Waiting"]++
			case "serving":
				stats["Serving"]++
			case "completed":
				stats["Completed"]++
			}
		}
	}

	return &TicketListResult{
		Tickets:    tickets,
		Stats:      stats,
		TotalCount: totalCount,
	}, nil
}

// CancelTicket cancels a ticket
func (s *StaffService) CancelTicket(ctx context.Context, ticketID int) error {
	return s.ticketRepo.UpdateStatus(ctx, ticketID, "cancelled")
}

// ResetYesterdayTickets resets all yesterday's waiting tickets
func (s *StaffService) ResetYesterdayTickets(ctx context.Context) (int, error) {
	return s.ticketRepo.CancelYesterdayWaiting(ctx)
}
