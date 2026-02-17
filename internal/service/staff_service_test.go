package service

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	"tenangantri/internal/model"
)

func TestStaffService_CallNext(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockUserCounterRepo := new(MockUserCounterRepository)
	mockCounterRepo := new(MockCounterRepository)
	mockCounterCategoryRepo := new(MockCounterCategoryRepository)
	mockTicketRepo := new(MockTicketRepository)
	mockStatsRepo := new(MockStatsRepository)
	mockCatRepo := new(MockCategoryRepository)

	service := NewStaffService(mockUserRepo, mockUserCounterRepo, mockCounterRepo, mockCounterCategoryRepo, mockTicketRepo, mockStatsRepo, mockCatRepo)

	ctx := context.Background()
	staffID := 1
	counterID := 1
	categoryID := 1

	mockUserCounterRepo.On("GetCounterIDByUserID", ctx, staffID).Return(sql.NullInt64{Int64: int64(counterID), Valid: true}, nil)

	mockCounterRepo.On("GetByID", ctx, counterID).Return(&model.Counter{
		ID:     counterID,
		Status: "active",
	}, nil)

	mockCounterCategoryRepo.On("GetCategoryIDsByCounterID", ctx, counterID).Return([]int{categoryID}, nil)

	mockTicketRepo.On("GetCurrentForCounter", ctx, counterID).Return(nil, nil)
	mockTicketRepo.On("GetNextTicket", ctx, []int{categoryID}).Return(&model.Ticket{
		ID:           10,
		TicketNumber: "A010",
	}, nil)

	mockTicketRepo.On("AssignToCounter", ctx, 10, counterID).Return(nil)
	mockCounterRepo.On("UpdateStatus", ctx, counterID, "serving").Return(nil)
	mockTicketRepo.On("GetWithDetails", ctx, 10).Return(&model.Ticket{
		ID:           10,
		TicketNumber: "A010",
		Status:       "serving",
		CounterID:    sql.NullInt64{Int64: int64(counterID), Valid: true},
	}, nil)

	ticket, err := service.CallNext(ctx, staffID)

	assert.NoError(t, err)
	assert.NotNil(t, ticket)
	assert.Equal(t, "A010", ticket.TicketNumber)
	assert.Equal(t, "serving", ticket.Status)

	mockCounterRepo.AssertExpectations(t)
	mockTicketRepo.AssertExpectations(t)
}
