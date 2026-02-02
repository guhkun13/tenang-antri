package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"queue-system/internal/model"
	"queue-system/internal/service"
	"queue-system/internal/websocket"
)

// AdminHandler handles admin-related requests
type AdminHandler struct {
	adminService *service.AdminService
	hub          *websocket.Hub
}

func NewAdminHandler(adminService *service.AdminService, hub *websocket.Hub) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
		hub:          hub,
	}
}

// Dashboard shows admin dashboard
func (h *AdminHandler) Dashboard(c *gin.Context) {
	stats, counters, categories, err := h.adminService.GetDashboardData(c.Request.Context())
	if err != nil {
		log.Error().Err(err).Str("layer", "handler").Msg("Failed to get dashboard data")
	}

	c.HTML(http.StatusOK, "pages/admin/dashboard.html", gin.H{
		"Stats":      stats,
		"Counters":   counters,
		"Categories": categories,
	})
}

// GetStats gets dashboard statistics
func (h *AdminHandler) GetStats(c *gin.Context) {
	stats, err := h.adminService.GetStats(c.Request.Context())
	if err != nil {
		log.Error().Err(err).Str("layer", "handler").Msg("Failed to get stats")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// User Management

// ListUsers shows users page
func (h *AdminHandler) ListUsers(c *gin.Context) {
	role := c.Query("role")
	users, err := h.adminService.ListUsers(c.Request.Context(), role)
	if err != nil {
		log.Error().Err(err).Str("layer", "handler").Str("func", "ListUsers").Msg("Failed to load users")
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"Error": "Failed to load users"})
		return
	}

	counters, _ := h.adminService.ListCounters(c.Request.Context())

	c.HTML(http.StatusOK, "pages/admin/users.html", gin.H{
		"Users":    users,
		"Counters": counters,
	})
}

// CreateUser creates a new user
func (h *AdminHandler) CreateUser(c *gin.Context) {
	var req model.CreateUserRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.adminService.CreateUser(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// GetUser gets a user by ID
func (h *AdminHandler) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.adminService.GetUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser updates a user
func (h *AdminHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.adminService.UpdateUserProfile(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser deletes a user
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.adminService.DeleteUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// ResetUserPassword resets user password
func (h *AdminHandler) ResetUserPassword(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	password, err := h.adminService.ResetUserPassword(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully", "password": password})
}

// Category Management

// ListCategories shows categories page
func (h *AdminHandler) ListCategories(c *gin.Context) {
	categories, err := h.adminService.ListCategories(c.Request.Context(), false)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"Error": "Failed to load categories"})
		return
	}

	c.HTML(http.StatusOK, "pages/admin/categories.html", gin.H{
		"Categories": categories,
	})
}

// CreateCategory creates a new category
func (h *AdminHandler) CreateCategory(c *gin.Context) {
	var req model.CreateCategoryRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := h.adminService.CreateCategory(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	h.hub.Broadcast("category_created", category)
	c.JSON(http.StatusCreated, category)
}

// UpdateCategory updates a category
func (h *AdminHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var req model.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := h.adminService.UpdateCategory(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	h.hub.Broadcast("category_updated", category)
	c.JSON(http.StatusOK, category)
}

// DeleteCategory deletes a category
func (h *AdminHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	err = h.adminService.DeleteCategory(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	h.hub.Broadcast("category_deleted", id)
	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}

// Counter Management

// ListCounters shows counters page
func (h *AdminHandler) ListCounters(c *gin.Context) {
	counters, err := h.adminService.ListCounters(c.Request.Context())
	if err != nil {
		log.Error().Err(err).Str("layer", "handler").Str("func", "ListCounters").Msg("Failed to list counters")
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"Error": "Failed to load counters"})
		return
	}

	categories, _ := h.adminService.ListCategories(c.Request.Context(), true)

	log.Info().Str("layer", "handler").Str("func", "ListCounters").Msg("Counters loaded successfully")

	c.HTML(http.StatusOK, "pages/admin/counters.html", gin.H{
		"Counters":   counters,
		"Categories": categories,
	})
}

// CreateCounter creates a new counter
func (h *AdminHandler) CreateCounter(c *gin.Context) {
	var req model.CreateCounterRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	counter, err := h.adminService.CreateCounter(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create counter"})
		return
	}

	h.hub.Broadcast("counter_created", counter)
	c.JSON(http.StatusCreated, counter)
}

// UpdateCounter updates a counter
func (h *AdminHandler) UpdateCounter(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid counter ID"})
		return
	}

	var req model.CreateCounterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	counter, err := h.adminService.UpdateCounter(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update counter"})
		return
	}

	h.hub.Broadcast("counter_updated", counter)
	c.JSON(http.StatusOK, counter)
}

// DeleteCounter deletes a counter
func (h *AdminHandler) DeleteCounter(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid counter ID"})
		return
	}

	err = h.adminService.DeleteCounter(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete counter"})
		return
	}

	h.hub.Broadcast("counter_deleted", id)
	c.JSON(http.StatusOK, gin.H{"message": "Counter deleted successfully"})
}

// GetCounter gets a counter by ID
func (h *AdminHandler) GetCounter(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid counter ID"})
		return
	}

	counter, err := h.adminService.GetCounter(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Counter not found"})
		return
	}

	c.JSON(http.StatusOK, counter)
}

// Ticket Management

// ListTickets shows tickets page
func (h *AdminHandler) ListTickets(c *gin.Context) {
	filters := make(map[string]interface{})

	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if categoryID := c.Query("category_id"); categoryID != "" {
		if id, err := strconv.Atoi(categoryID); err == nil {
			filters["category_id"] = id
		}
	}
	if counterID := c.Query("counter_id"); counterID != "" {
		if id, err := strconv.Atoi(counterID); err == nil {
			filters["counter_id"] = id
		}
	}
	if dateFrom := c.Query("date_from"); dateFrom != "" {
		filters["date_from"] = dateFrom
	}
	if dateTo := c.Query("date_to"); dateTo != "" {
		filters["date_to"] = dateTo
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		log.Error().Err(err).Str("layer", "handler").Str("func", "ListTickets").Msg("Failed to parse page")
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"Error": "Failed to load tickets"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if err != nil {
		log.Error().Err(err).Str("layer", "handler").Str("func", "ListTickets").Msg("Failed to parse limit")
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"Error": "Failed to load tickets"})
		return
	}
	filters["limit"] = limit
	filters["offset"] = (page - 1) * limit

	tickets, err := h.adminService.ListTickets(c.Request.Context(), filters)
	if err != nil {
		log.Error().Err(err).Str("layer", "handler").Str("func", "ListTickets").Msg("Failed to list tickets")
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"Error": "Failed to load tickets"})
		return
	}

	categories, err := h.adminService.ListCategories(c.Request.Context(), false)
	if err != nil {
		log.Error().Err(err).Str("layer", "handler").Str("func", "ListTickets").Msg("Failed to list categories")
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"Error": "Failed to load tickets"})
		return
	}

	counters, err := h.adminService.ListCounters(c.Request.Context())
	if err != nil {
		log.Error().Err(err).Str("layer", "handler").Str("func", "ListTickets").Msg("Failed to list counters")
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"Error": "Failed to load tickets"})
		return
	}

	c.HTML(http.StatusOK, "pages/admin/tickets.html", gin.H{
		"Tickets":    tickets,
		"Categories": categories,
		"Counters":   counters,
		"Filters":    filters,
	})
}

// GetTicket gets a ticket by ID
func (h *AdminHandler) GetTicket(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	ticket, err := h.adminService.GetTicket(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

// CreateTicket creates a new ticket
func (h *AdminHandler) CreateTicket(c *gin.Context) {
	var req model.CreateTicketRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticket, err := h.adminService.CreateTicket(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ticket"})
		return
	}

	// Broadcast update
	stats, _ := h.adminService.GetStats(c.Request.Context())
	h.hub.BroadcastStatsUpdate(stats)
	h.hub.BroadcastTicketUpdate(ticket)

	c.JSON(http.StatusCreated, ticket)
}

// UpdateTicketStatus updates ticket status
func (h *AdminHandler) UpdateTicketStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	var req model.UpdateTicketStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticket, err := h.adminService.UpdateTicketStatus(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ticket status"})
		return
	}

	// Broadcast updates
	stats, _ := h.adminService.GetStats(c.Request.Context())
	h.hub.BroadcastStatsUpdate(stats)
	h.hub.BroadcastTicketUpdate(ticket)

	c.JSON(http.StatusOK, ticket)
}

// CancelTicket cancels a ticket
func (h *AdminHandler) CancelTicket(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	ticket, err := h.adminService.CancelTicket(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel ticket"})
		return
	}

	// Broadcast updates
	stats, _ := h.adminService.GetStats(c.Request.Context())
	h.hub.BroadcastStatsUpdate(stats)
	h.hub.BroadcastTicketUpdate(ticket)

	c.JSON(http.StatusOK, gin.H{"message": "Ticket cancelled successfully"})
}

// Reports

// Reports shows reports page
func (h *AdminHandler) Reports(c *gin.Context) {
	dateFrom := c.DefaultQuery("date_from", time.Now().AddDate(0, 0, -7).Format("2006-01-02"))
	dateTo := c.DefaultQuery("date_to", time.Now().Format("2006-01-02"))

	c.HTML(http.StatusOK, "pages/admin/reports.html", gin.H{
		"DateFrom": dateFrom,
		"DateTo":   dateTo,
	})
}

// ExportTickets exports tickets to CSV/Excel
func (h *AdminHandler) ExportTickets(c *gin.Context) {
	// Implementation for exporting tickets to CSV/Excel
	c.JSON(http.StatusOK, gin.H{"message": "Export functionality - TODO"})
}
