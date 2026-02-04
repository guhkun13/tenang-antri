package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"tenangantri/internal/service"
)

// TrackingHandler handles ticket tracking requests
type TrackingHandler struct {
	trackingService *service.TrackingService
}

func NewTrackingHandler(trackingService *service.TrackingService) *TrackingHandler {
	return &TrackingHandler{
		trackingService: trackingService,
	}
}

// ShowTrackingPage renders the main tracking page
func (h *TrackingHandler) ShowTrackingPage(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/track/index.html", gin.H{})
}

// GetTrackingInfo returns tracking information for a ticket (HTMX endpoint)
func (h *TrackingHandler) GetTrackingInfo(c *gin.Context) {
	ticketNumber := strings.ToUpper(strings.TrimSpace(c.Param("ticket_number")))

	if ticketNumber == "" {
		c.HTML(http.StatusBadRequest, "pages/track/_tracking_info.html", gin.H{
			"Error": "Please enter a ticket number",
		})
		return
	}

	trackingInfo, err := h.trackingService.GetTicketTrackingInfo(c.Request.Context(), ticketNumber)
	if err != nil {
		log.Error().Err(err).Str("ticket_number", ticketNumber).Msg("Failed to get tracking info")
		c.HTML(http.StatusOK, "pages/track/_tracking_info.html", gin.H{
			"Error": "Ticket not found. Please check your ticket number and try again.",
		})
		return
	}

	c.HTML(http.StatusOK, "pages/track/_tracking_info.html", gin.H{
		"TrackingInfo": trackingInfo,
	})
}
