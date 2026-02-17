package service

import (
	"context"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"tenangantri/internal/dto"
	"tenangantri/internal/model"
	"tenangantri/internal/repository"
)

// UserService handles user-related business logic
type UserService struct {
	userRepo        repository.UserRepository
	userCounterRepo repository.UserCounterRepository
}

func NewUserService(userRepo repository.UserRepository, userCounterRepo repository.UserCounterRepository) *UserService {
	return &UserService{
		userRepo:        userRepo,
		userCounterRepo: userCounterRepo,
	}
}

// GetUserByUsername retrieves a user by username
func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	return s.userRepo.GetByUsername(ctx, username)
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(ctx context.Context, id int) (*model.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

// CreateUser creates a new user with hashed password
func (s *UserService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*model.User, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		FullName: sql.NullString{String: req.FullName, Valid: req.FullName != ""},
		Email:    sql.NullString{String: req.Email, Valid: req.Email != ""},
		Phone:    sql.NullString{String: req.Phone, Valid: req.Phone != ""},
		Role:     req.Role,
		IsActive: true,
	}

	createdUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	// Create user-counter association if counter_id is provided
	if req.CounterID != nil {
		_, err = s.userCounterRepo.Create(ctx, createdUser.ID, *req.CounterID)
		if err != nil {
			// Log error but don't fail user creation
			// Could roll back user creation here if strict consistency is needed
		}
	}

	return createdUser, nil
}

// UpdateUser updates user information
func (s *UserService) UpdateUser(ctx context.Context, id int, req *dto.UpdateUserRequest) (*model.User, error) {
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
func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	// Delete user-counter association first
	_ = s.userCounterRepo.DeleteByUserID(ctx, id)
	return s.userRepo.Delete(ctx, id)
}

// UpdateUserPassword updates a user's password
func (s *UserService) UpdateUserPassword(ctx context.Context, id int, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(ctx, id, string(hashedPassword))
}

// ResetUserPassword resets a user's password to a default
func (s *UserService) ResetUserPassword(ctx context.Context, id int) (string, error) {
	newPassword := "password123" // Generate random password in production
	err := s.UpdateUserPassword(ctx, id, newPassword)
	return newPassword, err
}

// UpdateLastLogin updates last login timestamp
func (s *UserService) UpdateLastLogin(ctx context.Context, id int) error {
	return s.userRepo.UpdateLastLogin(ctx, id)
}

// ListUsers retrieves users with optional role filter
func (s *UserService) ListUsers(ctx context.Context, role string) ([]model.User, error) {
	return s.userRepo.List(ctx, role)
}

// ValidatePassword validates user password
func (s *UserService) ValidatePassword(ctx context.Context, username string, password string) error {
	userWithPass, err := s.userRepo.GetByUsernameWithPassword(ctx, username)
	if err != nil {
		return errors.New("invalid username or password")
	}

	user, err := s.userRepo.GetByID(ctx, userWithPass.ID)
	if err != nil {
		return errors.New("invalid username or password")
	}

	if !user.IsActive {
		return errors.New("account is disabled")
	}

	err = bcrypt.CompareHashAndPassword([]byte(userWithPass.Password), []byte(password))
	if err != nil {
		return errors.New("invalid username or password")
	}

	return nil
}

// UpdateProfile updates user profile information
func (s *UserService) UpdateProfile(ctx context.Context, userID int, req *dto.UpdateProfileRequest) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.FullName = sql.NullString{String: req.FullName, Valid: req.FullName != ""}
	user.Email = sql.NullString{String: req.Email, Valid: req.Email != ""}
	user.Phone = sql.NullString{String: req.Phone, Valid: req.Phone != ""}

	return s.userRepo.Update(ctx, user)
}

// ChangePassword changes a user's password
func (s *UserService) ChangePassword(ctx context.Context, userID int, req *dto.ChangePasswordRequest) error {
	userWithPass, err := s.userRepo.GetByIDWithPassword(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(userWithPass.Password), []byte(req.CurrentPassword))
	if err != nil {
		return errors.New("current password is incorrect")
	}

	return s.UpdateUserPassword(ctx, userID, req.NewPassword)
}

// GetUserCounterID retrieves the counter ID assigned to a user
func (s *UserService) GetUserCounterID(ctx context.Context, userID int) (sql.NullInt64, error) {
	return s.userCounterRepo.GetCounterIDByUserID(ctx, userID)
}
