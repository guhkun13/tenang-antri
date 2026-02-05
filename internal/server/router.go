package server

import (
	"errors"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"tenangantri/internal/config"
	"tenangantri/internal/handler"
	"tenangantri/internal/middleware"
	"tenangantri/internal/model"
	"tenangantri/internal/repository"
	"tenangantri/internal/service"
	"tenangantri/internal/websocket"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func NewRouter(cfg *config.Config, pool *pgxpool.Pool) *gin.Engine {
	// Initialize repositories
	userRepo := repository.NewUserRepository(pool)
	counterRepo := repository.NewCounterRepository(pool)
	categoryRepo := repository.NewCategoryRepository(pool)
	ticketRepo := repository.NewTicketRepository(pool)
	statsRepo := repository.NewStatsRepository(pool)

	// Initialize services
	userService := service.NewUserService(userRepo)
	adminService := service.NewAdminService(userRepo, counterRepo, categoryRepo, ticketRepo, statsRepo)
	staffService := service.NewStaffService(userRepo, counterRepo, ticketRepo, statsRepo, categoryRepo)
	kioskService := service.NewKioskService(categoryRepo, ticketRepo, statsRepo)
	displayService := service.NewDisplayService(statsRepo, categoryRepo, counterRepo)
	trackingService := service.NewTrackingService(ticketRepo, categoryRepo, counterRepo)

	// Initialize JWT middleware
	middleware.InitAuth(&cfg.JWT)

	// Initialize WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// Initialize handlers
	authHandler := handler.NewAuthHandler(userService, &cfg.JWT)
	adminHandler := handler.NewAdminHandler(adminService, hub)
	staffHandler := handler.NewStaffHandler(staffService, hub)
	kioskHandler := handler.NewKioskHandler(kioskService, hub)
	displayHandler := handler.NewDisplayHandler(displayService)
	trackingHandler := handler.NewTrackingHandler(trackingService)

	// Setup router
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(loggerMiddleware())

	// Add custom template functions
	funcMap := template.FuncMap{
		"formatDate": func(t time.Time) string {
			return t.Format("2006-01-02 15:04:05")
		},
		"formatDuration": func(seconds int) string {
			if seconds < 60 {
				return "< 1 min"
			}
			minutes := seconds / 60
			if minutes < 60 {
				return "{{ . }} min"
			}
			// hours := minutes / 60
			// mins := minutes % 60
			return "{{ . }}h {{ . }}m"
		},
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"mul": func(a, b int) int {
			return a * b
		},
		"div": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a / b
		},
		"mod": func(a, b int) int {
			return a % b
		},
		"sum": func(items interface{}, field string) int {
			total := 0
			switch v := items.(type) {
			case []model.Category:
				for _, item := range v {
					if field == "Priority" {
						total += item.Priority
					}
				}
			}
			return total
		},
		"countActive": func(items []model.Category) int {
			count := 0
			for _, item := range items {
				if item.IsActive {
					count++
				}
			}
			return count
		},
		"countInactive": func(items []model.Category) int {
			count := 0
			for _, item := range items {
				if !item.IsActive {
					count++
				}
			}
			return count
		},
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, errors.New("invalid dict call: even number of arguments required")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, errors.New("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
		"js": func(s string) template.JS {
			return template.JS(s)
		},
		"uppercase": func(s string) string {
			return strings.ToUpper(s)
		},
		"upper": func(s string) string {
			return strings.ToUpper(s)
		},
		"now": func() time.Time {
			return time.Now()
		},
		"gt": func(a, b int) bool {
			return a > b
		},
	}

	// Create a new template set
	tmpl := template.New("").Funcs(funcMap)

	// Load templates manually to preserve relative paths
	err := filepath.Walk("web/templates", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".html") {
			// Get relative path as template name (e.g., "customer/index.html")
			name, err := filepath.Rel("web/templates", path)
			if err != nil {
				return err
			}
			// Normalise path separators to /
			name = filepath.ToSlash(name)

			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			_, err = tmpl.New(name).Parse(string(content))
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load templates")
	}

	// Set the template set to Gin
	r.SetHTMLTemplate(tmpl)

	// Static files
	r.Static("/static", "./web/static")

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

func loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		log.Info().
			Str("client_ip", clientIP).
			Str("method", method).
			Str("path", path).
			Int("status", statusCode).
			Dur("latency", latency).
			Msg("Request")
	}
}
