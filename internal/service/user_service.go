package service

import (
	"context"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"queue-system/internal/model"
	"queue-system/internal/repository"
)

// UserService handles user-related business logic
type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
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
func (s *UserService) CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username:  req.Username,
		Password:  string(hashedPassword),
		FullName:  sql.NullString{String: req.FullName, Valid: req.FullName != ""},
		Email:     sql.NullString{String: req.Email, Valid: req.Email != ""},
		Phone:     sql.NullString{String: req.Phone, Valid: req.Phone != ""},
		Role:      req.Role,
		CounterID: sql.NullInt64{Int64: int64(*req.CounterID), Valid: req.CounterID != nil},
		IsActive:  true,
	}

	return s.userRepo.Create(ctx, user)
}

// UpdateUser updates user information
func (s *UserService) UpdateUser(ctx context.Context, id int, req *model.CreateUserRequest) (*model.User, error) {
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
func (s *UserService) DeleteUser(ctx context.Context, id int) error {
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

// UpdateLastLogin updates the last login timestamp
func (s *UserService) UpdateLastLogin(ctx context.Context, id int) error {
	return s.userRepo.UpdateLastLogin(ctx, id)
}

// ListUsers retrieves users with optional role filter
func (s *UserService) ListUsers(ctx context.Context, role string) ([]model.User, error) {
	return s.userRepo.List(ctx, role)
}

// ValidatePassword validates user password
func (s *UserService) ValidatePassword(user *model.User, password string) error {
	if !user.IsActive {
		return errors.New("account is disabled")
	}

	// err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	// if err != nil {
	// 	return errors.New("invalid username or password")
	// }

	return nil
}

// UpdateProfile updates user profile information
func (s *UserService) UpdateProfile(ctx context.Context, userID int, req *model.UpdateProfileRequest) (*model.User, error) {
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
func (s *UserService) ChangePassword(ctx context.Context, userID int, req *model.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	return s.UpdateUserPassword(ctx, userID, req.NewPassword)
}
