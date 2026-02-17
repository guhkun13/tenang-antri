package service

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"

	"tenangantri/internal/dto"
	"tenangantri/internal/model"
	"tenangantri/internal/repository"
)

// getPriorityFromInterface converts priority interface{} to int
func getPriorityFromInterface(priority interface{}) int {
	switch v := priority.(type) {
	case int:
		return v
	case float64:
		return int(v)
	case string:
		if v == "" {
			return 0
		}
		if parsed, err := strconv.Atoi(v); err == nil {
			return parsed
		}
		return 0
	default:
		return 0
	}
}

// AdminService handles admin-specific business logic
type AdminService struct {
	userRepo            repository.UserRepository
	userCounterRepo     repository.UserCounterRepository
	counterRepo         repository.CounterRepository
	counterCategoryRepo repository.CounterCategoryRepository
	categoryRepo        repository.CategoryRepository
	ticketRepo          repository.TicketRepository
	statsRepo           repository.StatsRepository
}

func NewAdminService(userRepo repository.UserRepository,
	userCounterRepo repository.UserCounterRepository,
	counterRepo repository.CounterRepository,
	counterCategoryRepo repository.CounterCategoryRepository,
	categoryRepo repository.CategoryRepository,
	ticketRepo repository.TicketRepository,
	statsRepo repository.StatsRepository) *AdminService {
	return &AdminService{
		userRepo:            userRepo,
		userCounterRepo:     userCounterRepo,
		counterRepo:         counterRepo,
		counterCategoryRepo: counterCategoryRepo,
		categoryRepo:        categoryRepo,
		ticketRepo:          ticketRepo,
		statsRepo:           statsRepo,
	}
}

// GetDashboardData gets admin dashboard data
func (s *AdminService) GetDashboardData(ctx context.Context) (*dto.DashboardStats, []model.Counter, []model.Category, error) {
	stats, err := s.statsRepo.GetDashboardStats(ctx)
	if err != nil {
		log.Error().Err(err).Str("layer", "service").Msg("Failed to get dashboard stats")
		stats = &dto.DashboardStats{TicketsByStatus: make(map[string]int)}
	}

	log.Info().Interface("stats", stats).Msg("Dashboard stats")

	counters, err := s.counterRepo.List(ctx)
	if err != nil {
		log.Error().Err(err).Str("layer", "service").Str("func", "GetDashboardData").Msg("Failed to get counters")
	}

	categories, err := s.categoryRepo.List(ctx, true, false)
	if err != nil {
		log.Error().Err(err).Str("layer", "service").Str("func", "GetDashboardData").Msg("Failed to get categories")
	}

	log.Info().Str("layer", "service").Str("func", "GetDashboardData").Msg("Dashboard loaded successfully")

	return stats, counters, categories, nil
}

// GetStats gets dashboard statistics
func (s *AdminService) GetStats(ctx context.Context) (*dto.DashboardStats, error) {
	return s.statsRepo.GetDashboardStats(ctx)
}

// User Management methods

// CreateUser creates a new user
func (s *AdminService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*model.User, error) {
	// Delegate to UserService to ensure password is properly hashed
	userService := NewUserService(s.userRepo, s.userCounterRepo)
	return userService.CreateUser(ctx, req)
}

// GetUser gets a user by ID
func (s *AdminService) GetUser(ctx context.Context, id int) (*model.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

// UpdateUserProfile updates user profile (without password)
func (s *AdminService) UpdateUserProfile(ctx context.Context, id int, req *dto.UpdateUserRequest) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.FullName = sql.NullString{String: req.FullName, Valid: req.FullName != ""}
	user.Email = sql.NullString{String: req.Email, Valid: req.Email != ""}
	user.Phone = sql.NullString{String: req.Phone, Valid: req.Phone != ""}
	user.Role = req.Role

	updatedUser, err := s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	// Update user-counter association
	if req.CounterID != nil {
		// Delete existing association
		_ = s.userCounterRepo.DeleteByUserID(ctx, id)
		// Create new association
		_, _ = s.userCounterRepo.Create(ctx, id, *req.CounterID)
	} else {
		// Remove association if counter_id is not provided
		_ = s.userCounterRepo.DeleteByUserID(ctx, id)
	}

	return updatedUser, nil
}

// UpdateUser updates a user
func (s *AdminService) UpdateUser(ctx context.Context, id int, req *dto.UpdateUserRequest) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.FullName = sql.NullString{String: req.FullName, Valid: req.FullName != ""}
	user.Email = sql.NullString{String: req.Email, Valid: req.Email != ""}
	user.Phone = sql.NullString{String: req.Phone, Valid: req.Phone != ""}
	user.Role = req.Role

	updatedUser, err := s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	// Update user-counter association
	if req.CounterID != nil {
		// Delete existing association
		_ = s.userCounterRepo.DeleteByUserID(ctx, id)
		// Create new association
		_, _ = s.userCounterRepo.Create(ctx, id, *req.CounterID)
	} else {
		// Remove association if counter_id is not provided
		_ = s.userCounterRepo.DeleteByUserID(ctx, id)
	}

	return updatedUser, nil
}

// DeleteUser deletes a user
func (s *AdminService) DeleteUser(ctx context.Context, id int) error {
	// Delete user-counter association first
	_ = s.userCounterRepo.DeleteByUserID(ctx, id)
	return s.userRepo.Delete(ctx, id)
}

// ResetUserPassword resets a user's password
func (s *AdminService) ResetUserPassword(ctx context.Context, id int) (string, error) {
	userService := NewUserService(s.userRepo, s.userCounterRepo)
	return userService.ResetUserPassword(ctx, id)
}

// ListUsers lists users with optional role filter
func (s *AdminService) ListUsers(ctx context.Context, role string) ([]model.User, error) {
	return s.userRepo.List(ctx, role)
}

// Category Management methods

// CreateCategory creates a new category
func (s *AdminService) CreateCategory(ctx context.Context, req *dto.CreateCategoryRequest) (*model.Category, error) {
	category := &model.Category{
		Name:        req.Name,
		Prefix:      req.Prefix,
		Priority:    getPriorityFromInterface(req.Priority),
		ColorCode:   req.ColorCode,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
		Icon:        sql.NullString{String: req.Icon, Valid: req.Icon != ""},
		IsActive:    true,
	}

	return s.categoryRepo.Create(ctx, category)
}

// UpdateCategory updates a category
func (s *AdminService) UpdateCategory(ctx context.Context, id int, req *dto.CreateCategoryRequest) (*model.Category, error) {
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	category.Name = req.Name
	category.Prefix = req.Prefix
	category.Priority = getPriorityFromInterface(req.Priority)
	category.ColorCode = req.ColorCode
	category.Description = sql.NullString{String: req.Description, Valid: req.Description != ""}
	category.Icon = sql.NullString{String: req.Icon, Valid: req.Icon != ""}

	return s.categoryRepo.Update(ctx, category)
}

// UpdateCategoryStatus updates only the status of a category
func (s *AdminService) UpdateCategoryStatus(ctx context.Context, id int, isActive bool) (*model.Category, error) {
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	category.IsActive = isActive
	return s.categoryRepo.Update(ctx, category)
}

// DeleteCategory deletes a category
func (s *AdminService) DeleteCategory(ctx context.Context, id int) error {
	// Delete counter-category associations first
	_ = s.counterCategoryRepo.DeleteByCategoryID(ctx, id)
	return s.categoryRepo.Delete(ctx, id)
}

// GetCategory gets a category by ID
func (s *AdminService) GetCategory(ctx context.Context, id int) (*model.Category, error) {
	return s.categoryRepo.GetByID(ctx, id)
}

// ListCategories lists categories with optional active filter
func (s *AdminService) ListCategories(ctx context.Context, activeOnly bool) ([]model.Category, error) {
	return s.categoryRepo.List(ctx, activeOnly, false)
}

// Counter Management methods

// CreateCounter creates a new counter with categories
func (s *AdminService) CreateCounter(ctx context.Context, req *dto.CreateCounterRequest) (*model.Counter, error) {
	counter := &model.Counter{
		Number:   req.Number,
		Name:     sql.NullString{String: req.Name, Valid: req.Name != ""},
		Location: sql.NullString{String: req.Location, Valid: req.Location != ""},
		Status:   model.CounterStatusOffline,
	}

	createdCounter, err := s.counterRepo.Create(ctx, counter)
	if err != nil {
		return nil, err
	}

	// Add category associations if provided
	for _, categoryID := range req.CategoryIDs {
		_, _ = s.counterCategoryRepo.Create(ctx, createdCounter.ID, categoryID)
	}

	return createdCounter, nil
}

// UpdateCounter updates a counter and its categories
func (s *AdminService) UpdateCounter(ctx context.Context, id int, req *dto.CreateCounterRequest) (*model.Counter, error) {
	counter, err := s.counterRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	counter.Number = req.Number
	counter.Name = sql.NullString{String: req.Name, Valid: req.Name != ""}
	counter.Location = sql.NullString{String: req.Location, Valid: req.Location != ""}

	updatedCounter, err := s.counterRepo.Update(ctx, counter)
	if err != nil {
		return nil, err
	}

	// Update category associations - clear existing and add new
	_ = s.counterCategoryRepo.DeleteByCounterID(ctx, id)
	for _, categoryID := range req.CategoryIDs {
		_, _ = s.counterCategoryRepo.Create(ctx, id, categoryID)
	}

	return updatedCounter, nil
}

// DeleteCounter deletes a counter
func (s *AdminService) DeleteCounter(ctx context.Context, id int) error {
	// Delete counter-category associations first
	_ = s.counterCategoryRepo.DeleteByCounterID(ctx, id)
	// Delete user-counter associations
	_ = s.userCounterRepo.DeleteByCounterID(ctx, id)
	return s.counterRepo.Delete(ctx, id)
}

// ListCounters lists all counters
func (s *AdminService) ListCounters(ctx context.Context) ([]model.Counter, error) {
	counters, err := s.counterRepo.List(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to load counters")
		return nil, err
	}

	log.Info().Interface("counters", counters).Msg("Counters loaded successfully")

	return counters, nil
}

// GetCounterCategories gets categories served by a counter
func (s *AdminService) GetCounterCategories(ctx context.Context, counterID int) ([]int, error) {
	return s.counterCategoryRepo.GetCategoryIDsByCounterID(ctx, counterID)
}

// GetCounterWithCategories gets a counter with its assigned categories
func (s *AdminService) GetCounterWithCategories(ctx context.Context, counterID int) (*model.Counter, []int, error) {
	counter, err := s.counterRepo.GetByID(ctx, counterID)
	if err != nil {
		return nil, nil, err
	}

	categoryIDs, err := s.counterCategoryRepo.GetCategoryIDsByCounterID(ctx, counterID)
	if err != nil {
		return counter, nil, err
	}

	return counter, categoryIDs, nil
}

// AssignCategoriesToCounter assigns multiple categories to a counter
func (s *AdminService) AssignCategoriesToCounter(ctx context.Context, counterID int, categoryIDs []int) error {
	// Delete existing assignments
	if err := s.counterCategoryRepo.DeleteByCounterID(ctx, counterID); err != nil {
		return err
	}

	// Add new assignments
	for _, categoryID := range categoryIDs {
		if _, err := s.counterCategoryRepo.Create(ctx, counterID, categoryID); err != nil {
			return err
		}
	}

	return nil
}

// UpdateCounterStatus updates counter status
func (s *AdminService) UpdateCounterStatus(ctx context.Context, id int, status string) (*model.Counter, error) {
	counter, err := s.counterRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	counter.Status = status
	return s.counterRepo.Update(ctx, counter)
}

// GetCounter gets a counter by ID
func (s *AdminService) GetCounter(ctx context.Context, id int) (*model.Counter, error) {
	return s.counterRepo.GetByID(ctx, id)
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
func (s *AdminService) CreateTicket(ctx context.Context, req *dto.CreateTicketRequest) (*model.Ticket, error) {
	category, err := s.categoryRepo.GetByID(ctx, req.CategoryID)
	if err != nil {
		return nil, err
	}

	ticketNumber, dailySequence, err := s.ticketRepo.GenerateNumber(ctx, category.ID, category.Prefix)
	if err != nil {
		return nil, err
	}

	ticket := &model.Ticket{
		TicketNumber:  ticketNumber,
		DailySequence: dailySequence,
		QueueDate:     time.Now(),
		CategoryID:    sql.NullInt64{Int64: int64(req.CategoryID), Valid: true},
		Status:        "waiting",
		Priority:      req.Priority,
	}

	createdTicket, err := s.ticketRepo.Create(ctx, ticket)
	if err != nil {
		return nil, err
	}

	return s.ticketRepo.GetWithDetails(ctx, createdTicket.ID)
}

// UpdateTicketStatus updates ticket status
func (s *AdminService) UpdateTicketStatus(ctx context.Context, id int, req *dto.UpdateTicketStatusRequest) (*model.Ticket, error) {
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
