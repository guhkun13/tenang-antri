package dto

import (
	"encoding/json"
	"strconv"
)

// LoginRequest represents login credentials
type LoginRequest struct {
	Username string `json:"username" form:"username" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}

// LoginResponse represents login response
type LoginResponse struct {
	Token string `json:"token"`
}

// CreateTicketRequest represents ticket creation request
type CreateTicketRequest struct {
	CategoryID int `json:"category_id" form:"category_id" validate:"required"`
	Priority   int `json:"priority" form:"priority"`
}

// CallNextRequest represents call next ticket request
type CallNextRequest struct {
	CounterID int `json:"counter_id" form:"counter_id" validate:"required"`
}

// UpdateTicketStatusRequest represents ticket status update request
type UpdateTicketStatusRequest struct {
	Status string `json:"status" form:"status" validate:"required,oneof=waiting serving completed no_show cancelled"`
	Notes  string `json:"notes" form:"notes"`
}

// CreateCategoryRequest represents category creation request
type CreateCategoryRequest struct {
	Name        string      `json:"name" form:"name" validate:"required"`
	Prefix      string      `json:"prefix" form:"prefix" validate:"required"`
	Priority    interface{} `json:"priority" form:"priority"`
	ColorCode   string      `json:"color_code" form:"color_code"`
	Description string      `json:"description" form:"description"`
	Icon        string      `json:"icon" form:"icon"`
}

// UnmarshalJSON for CreateCategoryRequest to handle string priority
func (r *CreateCategoryRequest) UnmarshalJSON(data []byte) error {
	type Alias CreateCategoryRequest
	aux := &struct {
		Name        string      `json:"name"`
		Prefix      string      `json:"prefix"`
		Priority    interface{} `json:"priority"`
		ColorCode   string      `json:"color_code"`
		Description string      `json:"description"`
		Icon        string      `json:"icon"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	r.Name = aux.Name
	r.Prefix = aux.Prefix
	r.ColorCode = aux.ColorCode
	r.Description = aux.Description
	r.Icon = aux.Icon

	// Handle priority conversion
	switch v := aux.Priority.(type) {
	case int:
		r.Priority = v
	case float64:
		r.Priority = int(v)
	case string:
		if v == "" {
			r.Priority = 0
		} else {
			if parsed, err := strconv.Atoi(v); err == nil {
				r.Priority = parsed
			} else {
				r.Priority = 0
			}
		}
	default:
		r.Priority = 0
	}

	return nil
}

// UpdateCategoryStatusRequest represents category status update request
type UpdateCategoryStatusRequest struct {
	IsActive bool `json:"is_active"`
}

// CreateCounterRequest represents counter creation request
type CreateCounterRequest struct {
	Number      string `json:"number" form:"number" validate:"required"`
	Name        string `json:"name" form:"name" validate:"required"`
	Location    string `json:"location" form:"location"`
	CategoryIDs []int  `json:"category_ids" form:"category_ids"`
}

// CreateUserRequest represents user creation request
type CreateUserRequest struct {
	Username  string `json:"username" form:"username" validate:"required"`
	Password  string `json:"password" form:"password" validate:"required,min=6"`
	FullName  string `json:"full_name" form:"full_name" validate:"required"`
	Email     string `json:"email" form:"email" validate:"email"`
	Phone     string `json:"phone" form:"phone"`
	Role      string `json:"role" form:"role" validate:"required,oneof=admin staff"`
	CounterID *int   `json:"counter_id" form:"counter_id"`
}

// UpdateProfileRequest represents profile update request
type UpdateProfileRequest struct {
	FullName string `json:"full_name" form:"full_name"`
	Email    string `json:"email" form:"email"`
	Phone    string `json:"phone" form:"phone"`
}

// ChangePasswordRequest represents password change request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" form:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" form:"new_password" validate:"required,min=6"`
}

// UpdateUserRequest represents user update request (without password)
type UpdateUserRequest struct {
	FullName  string `json:"full_name" form:"full_name"`
	Email     string `json:"email" form:"email"`
	Phone     string `json:"phone" form:"phone"`
	Role      string `json:"role" form:"role" validate:"required,oneof=admin staff"`
	CounterID *int   `json:"counter_id" form:"counter_id"`
}
