package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"queue-system/internal/config"
	"queue-system/internal/dto"
	"queue-system/internal/middleware"
	"queue-system/internal/model"
	"queue-system/internal/service"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	userService *service.UserService
	config      *config.JWTConfig
}

func NewAuthHandler(userService *service.UserService, cfg *config.JWTConfig) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		config:      cfg,
	}
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	err := c.ShouldBind(&req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to bind login request")
		c.HTML(http.StatusBadRequest, "pages/login.html", gin.H{
			"Error": "Invalid input: " + err.Error(),
		})
		return
	}

	user, err := h.userService.GetUserByUsername(c.Request.Context(), req.Username)
	if err != nil {
		log.Error().Err(err).Str("layer", "handler").Msg("Failed to get user by username")
		c.HTML(http.StatusUnauthorized, "pages/login.html", gin.H{
			"Error": "Invalid username or password",
		})
		return
	}

	if err := h.userService.ValidatePassword(user, req.Password); err != nil {
		log.Error().Err(err).Str("layer", "handler").Msg("Invalid username or password")
		c.HTML(http.StatusUnauthorized, "pages/login.html", gin.H{
			"Error": "Invalid username or password",
		})
		return
	}

	// Update last login
	err = h.userService.UpdateLastLogin(c.Request.Context(), user.ID)
	if err != nil {
		log.Error().Err(err).Str("layer", "handler").Msg("Failed to update last login")
	}

	// Generate token
	token, err := middleware.GenerateToken(user.ID, user.Username, user.Role, h.config.AccessTokenExpiry)
	if err != nil {
		log.Error().Err(err).Str("layer", "handler").Msg("Failed to generate token")
		c.HTML(http.StatusInternalServerError, "pages/login.html", gin.H{
			"Error": "Failed to generate token",
		})
		return
	}

	// Set cookie
	c.SetCookie("auth_token", token, int(h.config.AccessTokenExpiry.Seconds()), "/", "", false, true)

	// Redirect based on role
	if user.Role == "admin" {
		c.Redirect(http.StatusFound, "/admin/dashboard")
	} else {
		c.Redirect(http.StatusFound, "/staff/dashboard")
	}
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("auth_token", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login")
}

// ShowLogin shows the login page
func (h *AuthHandler) ShowLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/login.html", nil)
}

// GetProfile gets user profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	user, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get profile"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateProfile updates user profile
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.UpdateProfile(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// ChangePassword changes user password
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userService.ChangePassword(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

// ShowProfile shows the profile page
func (h *AuthHandler) ShowProfile(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	role := middleware.GetCurrentUserRole(c)

	user, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"Error": "Failed to load profile"})
		return
	}

	var counter *model.Counter
	if user.CounterID.Valid {
		// Get counter information would require counter service
		counter = &model.Counter{}
	}

	template := "pages/staff/profile.html"
	if role == "admin" {
		template = "pages/admin/profile.html"
	}

	c.HTML(http.StatusOK, template, gin.H{
		"User":    user,
		"Counter": counter,
	})
}
