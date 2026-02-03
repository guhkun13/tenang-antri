package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"tenangantri/internal/dto"
	"tenangantri/internal/model"
	"tenangantri/internal/service"
	"tenangantri/internal/websocket"
)

// KioskHandler handles kiosk-related requests
type KioskHandler struct {
	kioskService *service.KioskService
	hub          *websocket.Hub
}

func NewKioskHandler(kioskService *service.KioskService, hub *websocket.Hub) *KioskHandler {
	return &KioskHandler{
		kioskService: kioskService,
		hub:          hub,
	}
}

// ShowKiosk shows kiosk main page
func (h *KioskHandler) ShowKiosk(c *gin.Context) {
	categories, err := h.kioskService.GetCategories(c.Request.Context())
	log.Info().Interface("categories", categories).Msg("Categories")
	if err != nil {
		log.Error().Err(err).Msg("Failed to list categories")
		categories = []model.Category{}
	}

	c.HTML(http.StatusOK, "pages/kiosk/index.html", gin.H{
		"Categories": categories,
	})
}

// GenerateTicket generates a new ticket from kiosk
func (h *KioskHandler) GenerateTicket(c *gin.Context) {
	var req dto.CreateTicketRequest
	if err := c.ShouldBind(&req); err != nil {
		if c.GetHeader("HX-Request") != "" {
			c.HTML(http.StatusBadRequest, "pages/kiosk/ticket_error.html", gin.H{
				"Error": "Please select a service category",
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	ticket, queuePosition, estimatedWaitTime, err := h.kioskService.GenerateTicket(c.Request.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate ticket")
		if c.GetHeader("HX-Request") != "" {
			c.HTML(http.StatusInternalServerError, "pages/kiosk/ticket_error.html", gin.H{
				"Error": "Failed to generate ticket. Please try again.",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate ticket"})
		}
		return
	}

	// Broadcast updates
	stats, _, _ := h.kioskService.GetQueueInfo(c.Request.Context())
	if stats != nil {
		h.hub.BroadcastStatsUpdate(stats)
	}
	h.hub.BroadcastTicketUpdate(ticket)

	// Check if HTMX request
	if c.GetHeader("HX-Request") != "" {
		c.HTML(http.StatusOK, "pages/kiosk/ticket_preview.html", gin.H{
			"Ticket":            ticket,
			"QueuePosition":     queuePosition,
			"EstimatedWaitTime": estimatedWaitTime,
		})
	} else {
		c.JSON(http.StatusCreated, gin.H{
			"ticket":              ticket,
			"queue_position":      queuePosition,
			"estimated_wait_time": estimatedWaitTime,
		})
	}
}

// GetTicketStatus gets ticket status
func (h *KioskHandler) GetTicketStatus(c *gin.Context) {
	ticketNumber := c.Param("number")

	// Find ticket by number
	// Note: This would require adding a method to service
	// For now, we'll return a placeholder
	c.JSON(http.StatusOK, gin.H{
		"ticket_number": ticketNumber,
		"status":        "waiting",
		"message":       "Your ticket is in queue",
	})
}

// PrintTicket shows printable ticket view
func (h *KioskHandler) PrintTicket(c *gin.Context) {
	ticketNumber := c.Param("number")

	// In a real implementation, this would trigger a print job
	// For now, we just return a printable view
	c.HTML(http.StatusOK, "pages/kiosk/print_ticket.html", gin.H{
		"TicketNumber": ticketNumber,
		"Date":         time.Now().Format("2006-01-02 15:04:05"),
	})
}

// GetQueueInfo gets queue information for kiosk
func (h *KioskHandler) GetQueueInfo(c *gin.Context) {
	stats, categories, err := h.kioskService.GetQueueInfo(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get queue info"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"waiting_tickets": stats.WaitingTickets,
		"avg_wait_time":   stats.AvgWaitTime,
		"categories":      categories,
	})
}
