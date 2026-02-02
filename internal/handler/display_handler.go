package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"queue-system/internal/model"
	"queue-system/internal/service"
)

// DisplayHandler handles display-related requests
type DisplayHandler struct {
	displayService *service.DisplayService
}

func NewDisplayHandler(displayService *service.DisplayService) *DisplayHandler {
	return &DisplayHandler{
		displayService: displayService,
	}
}

// ShowDisplay shows the main display page
func (h *DisplayHandler) ShowDisplay(c *gin.Context) {
	tickets, categories, counters, err := h.displayService.GetDisplayData(c.Request.Context())
	if err != nil {
		log.Error().Err(err).Str("layer", "handler").Str("func", "ShowDisplay").Msg("Failed to get display data")
		tickets = []model.DisplayTicket{}
		categories = []model.Category{}
		counters = []model.Counter{}
	}

	c.HTML(http.StatusOK, "pages/display/index.html", gin.H{
		"Tickets":    tickets,
		"Categories": categories,
		"Counters":   counters,
	})
}

// GetCurrentlyServing gets currently serving tickets
func (h *DisplayHandler) GetCurrentlyServing(c *gin.Context) {
	tickets, err := h.displayService.GetCurrentlyServing(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get currently serving tickets"})
		return
	}

	c.JSON(http.StatusOK, tickets)
}

// GetQueueStats gets queue statistics
func (h *DisplayHandler) GetQueueStats(c *gin.Context) {
	stats, err := h.displayService.GetQueueStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get queue stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetWaitingByCategory gets waiting tickets count by category
func (h *DisplayHandler) GetWaitingByCategory(c *gin.Context) {
	queueStats, err := h.displayService.GetWaitingByCategory(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get queue stats"})
		return
	}

	c.JSON(http.StatusOK, queueStats)
}

// ShowCategoryDisplay shows category-specific display
func (h *DisplayHandler) ShowCategoryDisplay(c *gin.Context) {
	categoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"Error": "Invalid category ID"})
		return
	}

	category, tickets, err := h.displayService.GetCategoryDisplayData(c.Request.Context(), categoryID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get category display data")
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"Error": "Failed to load display data"})
		return
	}

	c.HTML(http.StatusOK, "display/category.html", gin.H{
		"Category": category,
		"Tickets":  tickets,
	})
}
