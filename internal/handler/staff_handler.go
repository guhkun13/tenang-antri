package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"tenangantri/internal/middleware"
	"tenangantri/internal/service"
	"tenangantri/internal/websocket"
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

	data, err := h.staffService.GetDashboardData(c.Request.Context(), userID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"Error": "Failed to load dashboard data"})
		return
	}

	if data == nil || data.User == nil || data.Counter == nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"Error": "No counter assigned"})
		return
	}

	c.HTML(http.StatusOK, "pages/staff/dashboard.html", gin.H{
		"User":             data.User,
		"Counter":          data.Counter,
		"CurrentTicket":    data.CurrentTicket,
		"WaitingTickets":   data.WaitingTickets,
		"QueueStats":       data.QueueStats,
		"CompletedTickets": data.CompletedTickets,
		"CategoryIDs":      data.CategoryIDs,
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
		c.JSON(http.StatusConflict, gin.H{"error": "Selesaikan tiket saat ini sebelum memanggil berikutnya"})
		return
	}

	h.hub.BroadcastTicketUpdate(ticket)

	c.JSON(http.StatusOK, ticket)
}

// CallAgain calls the current ticket again
func (h *StaffHandler) CallAgain(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	ticket, err := h.staffService.CallAgain(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call ticket again"})
		return
	}

	if ticket == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Tidak ada ticket yang sedang dilayani"})
		return
	}

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

	data, err := h.staffService.GetQueueStatus(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get queue status"})
		return
	}

	if data == nil || data.Counter == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No counter assigned"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"counter":         data.Counter,
		"current_ticket":  data.CurrentTicket,
		"waiting_tickets": data.WaitingTickets,
		"queue_stats":     data.QueueStats,
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

// GetTicketDetail gets detailed information about a ticket
func (h *StaffHandler) GetTicketDetail(c *gin.Context) {
	ticketID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	ticket, err := h.staffService.GetTicketDetail(c.Request.Context(), ticketID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get ticket details"})
		return
	}

	if ticket == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

// TicketsPage shows the tickets management page for staff
func (h *StaffHandler) TicketsPage(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	// Parse query parameters
	filters := make(map[string]interface{})

	if dateFrom := c.Query("date_from"); dateFrom != "" {
		filters["date_from"] = dateFrom
	}
	if dateTo := c.Query("date_to"); dateTo != "" {
		filters["date_to"] = dateTo
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	// Parse pagination parameters
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || limit < 1 {
		limit = 20
	}

	// Parse sorting parameters
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	filters["page"] = page
	filters["limit"] = limit
	filters["offset"] = (page - 1) * limit
	filters["sort_by"] = sortBy
	filters["sort_order"] = sortOrder

	result, err := h.staffService.GetAllTickets(c.Request.Context(), userID, filters)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"Error": "Failed to load tickets"})
		return
	}

	// Calculate pagination info
	totalPages := (result.TotalCount + limit - 1) / limit
	if totalPages < 1 {
		totalPages = 1
	}

	hasPrev := page > 1
	hasNext := page < totalPages

	c.HTML(http.StatusOK, "pages/staff/tickets.html", gin.H{
		"Tickets":    result.Tickets,
		"Stats":      result.Stats,
		"TotalCount": result.TotalCount,
		"Page":       page,
		"Limit":      limit,
		"TotalPages": totalPages,
		"HasPrev":    hasPrev,
		"HasNext":    hasNext,
		"Filters":    filters,
		"SortBy":     sortBy,
		"SortOrder":  sortOrder,
	})
}

// CancelTicket cancels a ticket
func (h *StaffHandler) CancelTicket(c *gin.Context) {
	ticketID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	err = h.staffService.CancelTicket(c.Request.Context(), ticketID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel ticket"})
		return
	}

	h.hub.Broadcast("ticket_cancelled", gin.H{"message": "Ticket cancelled successfully"})

	c.JSON(http.StatusOK, gin.H{"message": "Ticket cancelled successfully"})
}

// ResetYesterdayTickets resets all yesterday's waiting tickets
func (h *StaffHandler) ResetYesterdayTickets(c *gin.Context) {
	count, err := h.staffService.ResetYesterdayTickets(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset yesterday tickets"})
		return
	}

	message := fmt.Sprintf("%d tiket kemarin berhasil dibatalkan", count)
	h.hub.Broadcast("yesterday_tickets_reset", gin.H{"message": message})

	c.JSON(http.StatusOK, gin.H{"message": message})
}
