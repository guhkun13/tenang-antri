package server

import (
	"net/http"
	"tenangantri/internal/middleware"
	"tenangantri/internal/websocket"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func NewRouter(handlers *Handlers) *gin.Engine {
	authHandler := handlers.AuthHandler
	adminHandler := handlers.AdminHandler
	staffHandler := handlers.StaffHandler
	kioskHandler := handlers.KioskHandler
	displayHandler := handlers.DisplayHandler
	trackingHandler := handlers.TrackingHandler
	hub := handlers.Hub

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware())

	// Add custom template functions
	funcMap := BuildFuncMap()

	// Load templates manually to preserve relative paths
	tmpl, err := LoadTemplate(funcMap)

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load templates")
	}

	// Set the template set to Gin
	r.SetHTMLTemplate(tmpl)

	// Static files
	r.Static("/static", "./web/static")
	r.Static("/templates", "./web/templates")

	// Public routes
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "/kiosk")
	})

	r.GET("/login", authHandler.ShowLogin)
	r.POST("/login", authHandler.Login)
	r.GET("/logout", authHandler.Logout)

	// Kiosk routes (public)
	kiosk := r.Group("/kiosk")
	{
		kiosk.GET("/", kioskHandler.ShowKiosk)
		kiosk.POST("/ticket", kioskHandler.GenerateTicket)
		kiosk.GET("/ticket/:number", kioskHandler.GetTicketStatus)
		kiosk.GET("/ticket/:number/print", kioskHandler.PrintTicket)
		kiosk.GET("/queue-info", kioskHandler.GetQueueInfo)
	}

	// Display routes (public)
	display := r.Group("/display")
	{
		display.GET("/", displayHandler.ShowDisplay)
		display.GET("/serving", displayHandler.GetCurrentlyServing)
		display.GET("/stats", displayHandler.GetQueueStats)
		display.GET("/waiting", displayHandler.GetWaitingByCategory)
		display.GET("/category/:id", displayHandler.ShowCategoryDisplay)
	}

	// Tracking routes (public)
	track := r.Group("/track")
	{
		track.GET("/", trackingHandler.ShowTrackingPage)
		track.GET("/info/:ticket_number", trackingHandler.GetTrackingInfo)
	}

	// WebSocket endpoint
	r.GET("/ws", func(c *gin.Context) {
		websocket.ServeWs(hub, c.Writer, c.Request)
	})

	// Protected routes
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// Profile routes
		protected.GET("/profile", authHandler.ShowProfile)
		protected.GET("/api/profile", authHandler.GetProfile)
		protected.PUT("/api/profile", authHandler.UpdateProfile)
		protected.POST("/api/change-password", authHandler.ChangePassword)

		// Staff routes
		staff := protected.Group("/staff")
		staff.Use(middleware.RoleMiddleware("staff", "admin"))
		{
			staff.GET("/dashboard", staffHandler.Dashboard)
			staff.POST("/call-next", staffHandler.CallNext)
			staff.POST("/complete", staffHandler.CompleteTicket)
			staff.POST("/no-show", staffHandler.MarkNoShow)
			staff.POST("/pause", staffHandler.PauseCounter)
			staff.POST("/resume", staffHandler.ResumeCounter)
			staff.GET("/queue-status", staffHandler.GetQueueStatus)
			staff.GET("/current-ticket", staffHandler.GetCurrentTicket)
			staff.POST("/transfer/:id", staffHandler.TransferTicket)
		}

		// Admin routes
		admin := protected.Group("/admin")
		admin.Use(middleware.RoleMiddleware("admin"))
		{
			// Dashboard
			admin.GET("/dashboard", adminHandler.Dashboard)
			admin.GET("/api/stats", adminHandler.GetStats)

			// Users
			admin.GET("/users", adminHandler.ListUsers)
			admin.GET("/api/users/:id", adminHandler.GetUser)
			admin.POST("/api/users", adminHandler.CreateUser)
			admin.PUT("/api/users/:id", adminHandler.UpdateUser)
			admin.DELETE("/api/users/:id", adminHandler.DeleteUser)
			admin.POST("/api/users/:id/reset-password", adminHandler.ResetUserPassword)

			// Categories
			admin.GET("/categories", adminHandler.ListCategories)
			admin.GET("/api/categories/:id", adminHandler.GetCategory)
			admin.POST("/api/categories", adminHandler.CreateCategory)
			admin.PUT("/api/categories/:id", adminHandler.UpdateCategory)
			admin.PUT("/api/categories/:id/status", adminHandler.UpdateCategoryStatus)
			admin.DELETE("/api/categories/:id", adminHandler.DeleteCategory)

			// Counters
			admin.GET("/counters", adminHandler.ListCounters)
			admin.GET("/api/counters/:id", adminHandler.GetCounter)
			admin.POST("/api/counters", adminHandler.CreateCounter)
			admin.PUT("/api/counters/:id", adminHandler.UpdateCounter)
			admin.PUT("/api/counters/:id/status", adminHandler.UpdateCounterStatus)
			admin.DELETE("/api/counters/:id", adminHandler.DeleteCounter)

			// Tickets
			admin.GET("/tickets", adminHandler.ListTickets)
			admin.GET("/api/tickets/:id", adminHandler.GetTicket)
			admin.POST("/api/tickets", adminHandler.CreateTicket)
			admin.PUT("/api/tickets/:id/status", adminHandler.UpdateTicketStatus)
			admin.POST("/api/tickets/:id/cancel", adminHandler.CancelTicket)

			// Reports
			admin.GET("/reports", adminHandler.Reports)
			admin.GET("/api/reports/data", adminHandler.GetReportData)
			admin.GET("/api/export/tickets", adminHandler.ExportTickets)
			admin.GET("/api/export/pdf", adminHandler.ExportPDF)
		}
	}

	return r
}
