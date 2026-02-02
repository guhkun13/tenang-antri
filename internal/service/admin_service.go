package service

import (
	"context"
	"database/sql"

	"github.com/rs/zerolog/log"

	"queue-system/internal/model"
	"queue-system/internal/repository"
)

// AdminService handles admin-specific business logic
type AdminService struct {
	userRepo     *repository.UserRepository
	counterRepo  *repository.CounterRepository
	categoryRepo *repository.CategoryRepository
	ticketRepo   *repository.TicketRepository
	statsRepo    *repository.StatsRepository
}

func NewAdminService(userRepo *repository.UserRepository, counterRepo *repository.CounterRepository, categoryRepo *repository.CategoryRepository, ticketRepo *repository.TicketRepository, statsRepo *repository.StatsRepository) *AdminService {
	return &AdminService{
		userRepo:     userRepo,
		counterRepo:  counterRepo,
		categoryRepo: categoryRepo,
		ticketRepo:   ticketRepo,
		statsRepo:    statsRepo,
	}
}

// GetDashboardData gets admin dashboard data
func (s *AdminService) GetDashboardData(ctx context.Context) (*model.DashboardStats, []model.Counter, []model.Category, error) {
	stats, err := s.statsRepo.GetDashboardStats(ctx)
	if err != nil {
		log.Error().Err(err).Str("layer", "service").Msg("Failed to get dashboard stats")
		stats = &model.DashboardStats{TicketsByStatus: make(map[string]int)}
	}

	counters, err := s.counterRepo.List(ctx)
	if err != nil {
		log.Error().Err(err).Str("layer", "service").Str("func", "GetDashboardData").Msg("Failed to get counters")
	}

	categories, err := s.categoryRepo.List(ctx, true)
	if err != nil {
		log.Error().Err(err).Str("layer", "service").Str("func", "GetDashboardData").Msg("Failed to get categories")
	}

	log.Info().Str("layer", "service").Str("func", "GetDashboardData").Msg("Dashboard loaded successfully")

	return stats, counters, categories, nil
}

// GetStats gets dashboard statistics
func (s *AdminService) GetStats(ctx context.Context) (*model.DashboardStats, error) {
	return s.statsRepo.GetDashboardStats(ctx)
}

// User Management methods

// CreateUser creates a new user
func (s *AdminService) CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	user := &model.User{
		Username:  req.Username,
		Password:  req.Password,
		FullName:  sql.NullString{String: req.FullName, Valid: req.FullName != ""},
		Email:     sql.NullString{String: req.Email, Valid: req.Email != ""},
		Phone:     sql.NullString{String: req.Phone, Valid: req.Phone != ""},
		Role:      req.Role,
		CounterID: sql.NullInt64{Int64: int64(*req.CounterID), Valid: req.CounterID != nil},
	}

	return s.userRepo.Create(ctx, user)
}

// GetUser gets a user by ID
func (s *AdminService) GetUser(ctx context.Context, id int) (*model.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

// UpdateUserProfile updates user profile (without password)
func (s *AdminService) UpdateUserProfile(ctx context.Context, id int, req *model.UpdateUserRequest) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.FullName = sql.NullString{String: req.FullName, Valid: req.FullName != ""}
	user.Email = sql.NullString{String: req.Email, Valid: req.Email != ""}
	user.Phone = sql.NullString{String: req.Phone, Valid: req.Phone != ""}
	user.Role = req.Role
	user.CounterID = sql.NullInt64{Int64: int64(*req.CounterID), Valid: req.CounterID != nil}

	return s.userRepo.Update(ctx, user)
}

// UpdateUser updates a user
func (s *AdminService) UpdateUser(ctx context.Context, id int, req *model.CreateUserRequest) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.FullName = sql.NullString{String: req.FullName, Valid: req.FullName != ""}
	user.Email = sql.NullString{String: req.Email, Valid: req.Email != ""}
	user.Phone = sql.NullString{String: req.Phone, Valid: req.Phone != ""}
	user.Role = req.Role
	user.CounterID = sql.NullInt64{Int64: int64(*req.CounterID), Valid: req.CounterID != nil}

	return s.userRepo.Update(ctx, user)
}

// DeleteUser deletes a user
func (s *AdminService) DeleteUser(ctx context.Context, id int) error {
	return s.userRepo.Delete(ctx, id)
}

// ResetUserPassword resets a user's password
func (s *AdminService) ResetUserPassword(ctx context.Context, id int) (string, error) {
	userService := NewUserService(s.userRepo)
	return userService.ResetUserPassword(ctx, id)
}

// ListUsers lists users with optional role filter
func (s *AdminService) ListUsers(ctx context.Context, role string) ([]model.User, error) {
	return s.userRepo.List(ctx, role)
}

// Category Management methods

// CreateCategory creates a new category
func (s *AdminService) CreateCategory(ctx context.Context, req *model.CreateCategoryRequest) (*model.Category, error) {
	category := &model.Category{
		Name:        req.Name,
		Prefix:      req.Prefix,
		Priority:    req.Priority,
		ColorCode:   req.ColorCode,
		Description: req.Description,
		Icon:        req.Icon,
		IsActive:    true,
	}

	return s.categoryRepo.Create(ctx, category)
}

// UpdateCategory updates a category
func (s *AdminService) UpdateCategory(ctx context.Context, id int, req *model.CreateCategoryRequest) (*model.Category, error) {
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	category.Name = req.Name
	category.Prefix = req.Prefix
	category.Priority = req.Priority
	category.ColorCode = req.ColorCode
	category.Description = req.Description
	category.Icon = req.Icon

	return s.categoryRepo.Update(ctx, category)
}

// DeleteCategory deletes a category
func (s *AdminService) DeleteCategory(ctx context.Context, id int) error {
	return s.categoryRepo.Delete(ctx, id)
}

// ListCategories lists categories with optional active filter
func (s *AdminService) ListCategories(ctx context.Context, activeOnly bool) ([]model.Category, error) {
	return s.categoryRepo.List(ctx, activeOnly)
}

// Counter Management methods

// CreateCounter creates a new counter
func (s *AdminService) CreateCounter(ctx context.Context, req *model.CreateCounterRequest) (*model.Counter, error) {
	counter := &model.Counter{
		Number:   req.Number,
		Name:     req.Name,
		Location: req.Location,
		Status:   "inactive",
		IsActive: true,
	}

	createdCounter, err := s.counterRepo.Create(ctx, counter)
	if err != nil {
		return nil, err
	}

	// Assign categories
	for _, catID := range req.CategoryIDs {
		s.counterRepo.AssignCategory(ctx, createdCounter.ID, catID)
	}

	// Load categories for the counter
	createdCounter.Categories, _ = s.counterRepo.GetCategories(ctx, createdCounter.ID)

	return createdCounter, nil
}

// UpdateCounter updates a counter
func (s *AdminService) UpdateCounter(ctx context.Context, id int, req *model.CreateCounterRequest) (*model.Counter, error) {
	counter, err := s.counterRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	counter.Number = req.Number
	counter.Name = req.Name
	counter.Location = req.Location

	updatedCounter, err := s.counterRepo.Update(ctx, counter)
	if err != nil {
		return nil, err
	}

	// Update categories
	s.counterRepo.ClearCategories(ctx, counter.ID)
	for _, catID := range req.CategoryIDs {
		s.counterRepo.AssignCategory(ctx, counter.ID, catID)
	}

	// Load categories for the counter
	updatedCounter.Categories, _ = s.counterRepo.GetCategories(ctx, counter.ID)

	return updatedCounter, nil
}

// DeleteCounter deletes a counter
func (s *AdminService) DeleteCounter(ctx context.Context, id int) error {
	return s.counterRepo.Delete(ctx, id)
}

// ListCounters lists all counters
func (s *AdminService) ListCounters(ctx context.Context) ([]model.Counter, error) {
	counters, err := s.counterRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	// Load categories for each counter
	for i := range counters {
		cats, _ := s.counterRepo.GetCategories(ctx, counters[i].ID)
		counters[i].Categories = cats
	}

	log.Info().
		Interface("counters", counters).
		Msg("Counters loaded successfully")

	return counters, nil
}

// GetCounter gets a counter by ID
func (s *AdminService) GetCounter(ctx context.Context, id int) (*model.Counter, error) {
	counter, err := s.counterRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Load categories for the counter
	categories, _ := s.counterRepo.GetCategories(ctx, counter.ID)
	counter.Categories = categories

	return counter, nil
}

// Ticket Management methods

// ListTickets lists tickets with optional filters
func (s *AdminService) ListTickets(ctx context.Context, filters map[string]interface{}) ([]model.Ticket, error) {
	return s.ticketRepo.List(ctx, filters)
}

// GetTicket gets a ticket by ID
func (s *AdminService) GetTicket(ctx context.Context, id int) (*model.Ticket, error) {
	return s.ticketRepo.GetWithDetails(ctx, id)
}

// CreateTicket creates a new ticket
func (s *AdminService) CreateTicket(ctx context.Context, req *model.CreateTicketRequest) (*model.Ticket, error) {
	return s.ticketRepo.Create(ctx, &model.Ticket{
		CategoryID: req.CategoryID,
		Status:     "waiting",
		Priority:   req.Priority,
	})
}

// UpdateTicketStatus updates ticket status
func (s *AdminService) UpdateTicketStatus(ctx context.Context, id int, req *model.UpdateTicketStatusRequest) (*model.Ticket, error) {
	err := s.ticketRepo.UpdateStatus(ctx, id, req.Status)
	if err != nil {
		return nil, err
	}
	return s.ticketRepo.GetWithDetails(ctx, id)
}

// CancelTicket cancels a ticket
func (s *AdminService) CancelTicket(ctx context.Context, id int) (*model.Ticket, error) {
	err := s.ticketRepo.UpdateStatus(ctx, id, "cancelled")
	if err != nil {
		return nil, err
	}
	return s.ticketRepo.GetWithDetails(ctx, id)
}
