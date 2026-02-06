package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"tenangantri/internal/dto"
	"tenangantri/internal/model"
)

func TestKioskService_GenerateTicket(t *testing.T) {
	mockCatRepo := new(MockCategoryRepository)
	mockTicketRepo := new(MockTicketRepository)
	mockStatsRepo := new(MockStatsRepository)

	service := NewKioskService(mockCatRepo, mockTicketRepo, mockStatsRepo)

	ctx := context.Background()
	catID := 1
	category := &model.Category{
		ID:     catID,
		Prefix: "A",
		Name:   "General",
	}

	req := &dto.CreateTicketRequest{
		CategoryID: catID,
		Priority:   1,
	}

	mockCatRepo.On("GetByID", ctx, catID).Return(category, nil)
	mockTicketRepo.On("GenerateNumber", ctx, catID, "A").Return("A001", 1, nil)
	mockTicketRepo.On("Create", ctx, mock.AnythingOfType("*model.Ticket")).Return(&model.Ticket{
		ID:           1,
		TicketNumber: "A001",
	}, nil)
	mockTicketRepo.On("GetWithDetails", ctx, 1).Return(&model.Ticket{
		ID:           1,
		TicketNumber: "A001",
	}, nil)
	mockTicketRepo.On("GetTodayCountByCategory", ctx, catID).Return(5, nil)
	mockStatsRepo.On("GetDashboardStats", ctx).Return(&model.DashboardStats{
		AvgWaitTime: 600, // 10 mins
	}, nil)

	ticket, position, waitTime, err := service.GenerateTicket(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, ticket)
	assert.Equal(t, "A001", ticket.TicketNumber)
	assert.Equal(t, 5, position)
	assert.Equal(t, 50, waitTime) // (600 * 5) / 60 = 50

	mockCatRepo.AssertExpectations(t)
	mockTicketRepo.AssertExpectations(t)
	mockStatsRepo.AssertExpectations(t)
}
