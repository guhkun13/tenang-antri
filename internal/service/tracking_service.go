package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"

	"tenangantri/internal/dto"
	"tenangantri/internal/model"
	"tenangantri/internal/repository"
)

// TrackingService handles ticket tracking business logic
type TrackingService struct {
	ticketRepo   repository.TicketRepository
	categoryRepo repository.CategoryRepository
	counterRepo  repository.CounterRepository
}

func NewTrackingService(
	ticketRepo repository.TicketRepository,
	categoryRepo repository.CategoryRepository,
	counterRepo repository.CounterRepository,
) *TrackingService {
	return &TrackingService{
		ticketRepo:   ticketRepo,
		categoryRepo: categoryRepo,
		counterRepo:  counterRepo,
	}
}

// GetTicketTrackingInfo retrieves comprehensive tracking information for a ticket
func (s *TrackingService) GetTicketTrackingInfo(ctx context.Context, ticketNumber string) (*dto.TrackingInfo, error) {
	// Get ticket by number
	ticket, err := s.ticketRepo.GetByTicketNumber(ctx, ticketNumber)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("ticket not found")
		}
		log.Error().Err(err).Str("ticket_number", ticketNumber).Msg("Failed to get ticket by number")
		return nil, fmt.Errorf("failed to retrieve ticket")
	}

	// Build tracking info
	trackingInfo := &dto.TrackingInfo{
		TicketNumber:     ticket.TicketNumber,
		Status:           ticket.Status,
		CreatedAt:        ticket.CreatedAt,
		QueuePosition:    0,
		OperationalHours: "Mon-Fri: 08:00 - 17:00",
	}

	// Add category information
	if ticket.CategoryID.Valid {
		catID := int(ticket.CategoryID.Int64)
		category, err := s.categoryRepo.GetByID(ctx, catID)
		if err == nil && category != nil {
			trackingInfo.CategoryName = category.Name
			trackingInfo.CategoryColor = category.ColorCode
		}

		lastCalled, err := s.ticketRepo.GetLastCalledByCategoryID(ctx, catID)
		if err != nil {
			log.Error().Err(err).Int("category_id", catID).Msg("Failed to get last called ticket")
		} else {
			trackingInfo.LastCalledTicketNumber = lastCalled
		}
	}

	if ticket.CounterID.Valid {
		counterID := int(ticket.CounterID.Int64)
		counter, err := s.counterRepo.GetByID(ctx, counterID)
		if err == nil && counter != nil {
			trackingInfo.CounterNumber = counter.Number
			if counter.Name.Valid {
				trackingInfo.CounterName = counter.Name.String
			}

			switch counter.Status {
			case "active":
				trackingInfo.CounterStatus = "Open"
			case "paused":
				trackingInfo.CounterStatus = "Paused"
			default:
				trackingInfo.CounterStatus = "Closed"
			}

			currentServing, err := s.ticketRepo.GetCurrentForCounter(ctx, counter.ID)
			if err == nil && currentServing != nil {
				trackingInfo.IsCounterServing = true
				trackingInfo.CounterCurrentServingTicket = currentServing.TicketNumber
			}
		}
	}

	// Calculate queue position for waiting tickets
	if ticket.Status == "waiting" {
		position, err := s.CalculateQueuePosition(ctx, ticket)
		if err != nil {
			log.Error().Err(err).Msg("Failed to calculate queue position")
		} else {
			trackingInfo.QueuePosition = position
		}
	}

	return trackingInfo, nil
}

// CalculateQueuePosition determines the position of the ticket in the queue
func (s *TrackingService) CalculateQueuePosition(ctx context.Context, ticket *model.Ticket) (int, error) {
	if ticket.Status != "waiting" {
		return 0, nil
	}

	// Count tickets with same category that are waiting and were created before this ticket
	filters := map[string]interface{}{
		"status":      "waiting",
		"category_id": int(ticket.CategoryID.Int64),
	}

	tickets, err := s.ticketRepo.List(ctx, filters)
	if err != nil {
		return 0, err
	}

	position := 0
	for _, t := range tickets {
		if t.CreatedAt.Before(ticket.CreatedAt) {
			position++
		}
	}

	return position, nil
}

// EstimateWaitTime calculates the estimated wait time in minutes
func (s *TrackingService) EstimateWaitTime(ctx context.Context, ticket *model.Ticket, position int) (int, error) {
	if ticket.Status != "waiting" {
		return 0, nil
	}

	// Default average service time per ticket (in minutes)
	// In a real implementation, this could be calculated from historical data
	const avgServiceTimeMinutes = 5

	// Estimate is position * average service time
	estimatedWait := position * avgServiceTimeMinutes

	return estimatedWait, nil
}
