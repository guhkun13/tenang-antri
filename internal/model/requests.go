package model

import (
	"encoding/json"
	"strconv"
)

// LoginResponse represents login response
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// CallNextRequest represents call next ticket request
type CallNextRequest struct {
	CounterID int `json:"counter_id" form:"counter_id" validate:"required"`
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

// CreateCounterRequest represents counter creation request
type CreateCounterRequest struct {
	Number      string `json:"number" form:"number" validate:"required"`
	Name        string `json:"name" form:"name" validate:"required"`
	Location    string `json:"location" form:"location"`
	CategoryIDs []int  `json:"category_ids" form:"category_ids"`
}
