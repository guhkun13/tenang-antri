package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"queue-system/internal/middleware"
	"queue-system/internal/service"
	"queue-system/internal/websocket"
)

// StaffHandler handles staff-related requests
type StaffHandler struct {
	staffService *service.StaffService
	hub          *websocket.Hub
}

func NewStaffHandler(staffService *service.StaffService, hub *websocket.Hub) *StaffHandler {
	return &StaffHandler{
		staffService: staffService,
		hub:          hub,
	}
}

// Dashboard shows staff dashboard
func (h *StaffHandler) Dashboard(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	user, counter, currentTicket, waitingTickets, queueStats, completedTickets, categoryIDs, err := h.staffService.GetDashboardData(c.Request.Context(), userID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"Error": "Failed to load dashboard data"})
		return
	}

	if user == nil || counter == nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"Error": "No counter assigned"})
		return
	}

	c.HTML(http.StatusOK, "pages/staff/dashboard.html", gin.H{
		"User":             user,
		"Counter":          counter,
		"CurrentTicket":    currentTicket,
		"WaitingTickets":   waitingTickets,
		"QueueStats":       queueStats,
		"CompletedTickets": completedTickets,
		"CategoryIDs":      categoryIDs,
	})
}

// CallNext calls the next ticket
func (h *StaffHandler) CallNext(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	ticket, err := h.staffService.CallNext(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call next ticket"})
		return
	}

	if ticket == nil {
		c.JSON(http.StatusOK, gin.H{"message": "No tickets in queue"})
		return
	}

	// Broadcast updates would require stats service
	// For now, broadcast ticket update
	h.hub.BroadcastTicketUpdate(ticket)

	c.JSON(http.StatusOK, ticket)
}

// CompleteTicket completes the current ticket
func (h *StaffHandler) CompleteTicket(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	err := h.staffService.CompleteTicket(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete ticket"})
		return
	}

	// Broadcast updates
	h.hub.Broadcast("ticket_completed", gin.H{"message": "Ticket completed successfully"})

	c.JSON(http.StatusOK, gin.H{"message": "Ticket completed successfully"})
}

// MarkNoShow marks current ticket as no-show
func (h *StaffHandler) MarkNoShow(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	err := h.staffService.MarkNoShow(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark no-show"})
		return
	}

	// Broadcast updates
	h.hub.Broadcast("ticket_no_show", gin.H{"message": "Ticket marked as no-show"})

	c.JSON(http.StatusOK, gin.H{"message": "Ticket marked as no-show"})
}

// PauseCounter pauses the counter
func (h *StaffHandler) PauseCounter(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	err := h.staffService.PauseCounter(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to pause counter"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Counter paused"})
}

// ResumeCounter resumes the counter
func (h *StaffHandler) ResumeCounter(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	err := h.staffService.ResumeCounter(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resume counter"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Counter resumed"})
}

// GetQueueStatus gets queue status
func (h *StaffHandler) GetQueueStatus(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	counter, currentTicket, waitingTickets, queueStats, err := h.staffService.GetQueueStatus(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get queue status"})
		return
	}

	if counter == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No counter assigned"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"counter":         counter,
		"current_ticket":  currentTicket,
		"waiting_tickets": waitingTickets,
		"queue_stats":     queueStats,
	})
}

// GetCurrentTicket gets current ticket
func (h *StaffHandler) GetCurrentTicket(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	ticket, err := h.staffService.GetCurrentTicket(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get current ticket"})
		return
	}

	if ticket == nil {
		c.JSON(http.StatusOK, gin.H{"ticket": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ticket": ticket})
}

// TransferTicket transfers a ticket to another counter
func (h *StaffHandler) TransferTicket(c *gin.Context) {
	ticketID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	var req struct {
		CounterID int `json:"counter_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticket, err := h.staffService.TransferTicket(c.Request.Context(), ticketID, req.CounterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to transfer ticket"})
		return
	}

	// Broadcast updates
	h.hub.BroadcastTicketUpdate(ticket)

	c.JSON(http.StatusOK, ticket)
}
