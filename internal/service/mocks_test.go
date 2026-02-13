package service

import (
	"context"

	"github.com/stretchr/testify/mock"

	"tenangantri/internal/model"
	"tenangantri/internal/repository"
)

type MockTicketRepository struct {
	mock.Mock
}

func (m *MockTicketRepository) GetByID(ctx context.Context, id int) (*model.Ticket, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Ticket), args.Error(1)
}

func (m *MockTicketRepository) GetWithDetails(ctx context.Context, id int) (*model.Ticket, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Ticket), args.Error(1)
}

func (m *MockTicketRepository) GetByTicketNumber(ctx context.Context, ticketNumber string) (*model.Ticket, error) {
	args := m.Called(ctx, ticketNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Ticket), args.Error(1)
}

func (m *MockTicketRepository) Create(ctx context.Context, ticket *model.Ticket) (*model.Ticket, error) {
	args := m.Called(ctx, ticket)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Ticket), args.Error(1)
}

func (m *MockTicketRepository) UpdateStatus(ctx context.Context, id int, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockTicketRepository) AssignToCounter(ctx context.Context, ticketID, counterID int) error {
	args := m.Called(ctx, ticketID, counterID)
	return args.Error(0)
}

func (m *MockTicketRepository) GetNextTicket(ctx context.Context, categoryIDs []int) (*model.Ticket, error) {
	args := m.Called(ctx, categoryIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Ticket), args.Error(1)
}

func (m *MockTicketRepository) GetCurrentForCounter(ctx context.Context, counterID int) (*model.Ticket, error) {
	args := m.Called(ctx, counterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Ticket), args.Error(1)
}

func (m *MockTicketRepository) List(ctx context.Context, filters map[string]interface{}) ([]model.Ticket, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).([]model.Ticket), args.Error(1)
}

func (m *MockTicketRepository) GetTodayCount(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

func (m *MockTicketRepository) GetTodayCountByCategory(ctx context.Context, categoryID int) (int, error) {
	args := m.Called(ctx, categoryID)
	return args.Int(0), args.Error(1)
}

func (m *MockTicketRepository) GenerateNumber(ctx context.Context, categoryID int, prefix string) (string, int, error) {
	args := m.Called(ctx, categoryID, prefix)
	return args.String(0), args.Int(1), args.Error(2)
}

func (m *MockTicketRepository) GetWaitingPreview(ctx context.Context, limit int) ([]model.Ticket, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]model.Ticket), args.Error(1)
}

func (m *MockTicketRepository) GetWaitingPreviewByCategories(ctx context.Context, categoryIDs []int, limit int) ([]model.Ticket, error) {
	args := m.Called(ctx, categoryIDs, limit)
	return args.Get(0).([]model.Ticket), args.Error(1)
}

func (m *MockTicketRepository) GetTodayCompletedByCategories(ctx context.Context, categoryIDs []int) ([]model.Ticket, error) {
	args := m.Called(ctx, categoryIDs)
	return args.Get(0).([]model.Ticket), args.Error(1)
}

func (m *MockTicketRepository) GetLastCalledByCategoryID(ctx context.Context, categoryID int) (string, error) {
	args := m.Called(ctx, categoryID)
	return args.String(0), args.Error(1)
}

func (m *MockTicketRepository) GetTodayByCategories(ctx context.Context, categoryIDs []int) ([]model.Ticket, error) {
	args := m.Called(ctx, categoryIDs)
	return args.Get(0).([]model.Ticket), args.Error(1)
}

func (m *MockTicketRepository) GetAllTodayTickets(ctx context.Context) ([]model.Ticket, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Ticket), args.Error(1)
}

func (m *MockTicketRepository) GetAllTicketsByCategories(ctx context.Context, categoryIDs []int) ([]model.Ticket, error) {
	args := m.Called(ctx, categoryIDs)
	return args.Get(0).([]model.Ticket), args.Error(1)
}

func (m *MockTicketRepository) CancelYesterdayWaiting(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) GetByID(ctx context.Context, id int) (*model.Category, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Category), args.Error(1)
}

func (m *MockCategoryRepository) List(ctx context.Context, activeOnly bool, withCountersOnly bool) ([]model.Category, error) {
	args := m.Called(ctx, activeOnly, withCountersOnly)
	return args.Get(0).([]model.Category), args.Error(1)
}

func (m *MockCategoryRepository) Create(ctx context.Context, category *model.Category) (*model.Category, error) {
	args := m.Called(ctx, category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Category), args.Error(1)
}

func (m *MockCategoryRepository) Update(ctx context.Context, category *model.Category) (*model.Category, error) {
	args := m.Called(ctx, category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Category), args.Error(1)
}

func (m *MockCategoryRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockStatsRepository struct {
	mock.Mock
}

func (m *MockStatsRepository) GetDashboardStats(ctx context.Context) (*model.DashboardStats, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.DashboardStats), args.Error(1)
}

func (m *MockStatsRepository) GetQueueLengthByCategory(ctx context.Context) ([]model.CategoryQueueStats, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.CategoryQueueStats), args.Error(1)
}

func (m *MockStatsRepository) GetQueueLengthByCategories(ctx context.Context, categoryIDs []int) ([]model.CategoryQueueStats, error) {
	args := m.Called(ctx, categoryIDs)
	return args.Get(0).([]model.CategoryQueueStats), args.Error(1)
}

func (m *MockStatsRepository) GetHourlyDistribution(ctx context.Context) ([]model.HourlyStats, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.HourlyStats), args.Error(1)
}

func (m *MockStatsRepository) GetCurrentlyServingTickets(ctx context.Context) ([]model.DisplayTicket, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.DisplayTicket), args.Error(1)
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsernameWithPassword(ctx context.Context, username string) (*repository.UserWithPassword, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.UserWithPassword), args.Error(1)
}

func (m *MockUserRepository) GetByIDWithPassword(ctx context.Context, id int) (*repository.UserWithPassword, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.UserWithPassword), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *model.User) (*model.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) UpdatePassword(ctx context.Context, id int, password string) error {
	args := m.Called(ctx, id, password)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateLastLogin(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context, role string) ([]model.User, error) {
	args := m.Called(ctx, role)
	return args.Get(0).([]model.User), args.Error(1)
}

type MockCounterRepository struct {
	mock.Mock
}

func (m *MockCounterRepository) GetByID(ctx context.Context, id int) (*model.Counter, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Counter), args.Error(1)
}

func (m *MockCounterRepository) Create(ctx context.Context, counter *model.Counter) (*model.Counter, error) {
	args := m.Called(ctx, counter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Counter), args.Error(1)
}

func (m *MockCounterRepository) Update(ctx context.Context, counter *model.Counter) (*model.Counter, error) {
	args := m.Called(ctx, counter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Counter), args.Error(1)
}

func (m *MockCounterRepository) UpdateStatus(ctx context.Context, id int, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockCounterRepository) UpdateStaff(ctx context.Context, counterID int, staffID *int) error {
	args := m.Called(ctx, counterID, staffID)
	return args.Error(0)
}

func (m *MockCounterRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCounterRepository) List(ctx context.Context) ([]model.Counter, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Counter), args.Error(1)
}
