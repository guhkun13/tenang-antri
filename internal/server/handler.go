package server

import (
	"tenangantri/internal/config"
	"tenangantri/internal/handler"
	"tenangantri/internal/middleware"
	"tenangantri/internal/repository"
	"tenangantri/internal/service"
	"tenangantri/internal/websocket"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Handlers struct {
	Hub             *websocket.Hub
	AuthHandler     *handler.AuthHandler
	AdminHandler    *handler.AdminHandler
	StaffHandler    *handler.StaffHandler
	KioskHandler    *handler.KioskHandler
	DisplayHandler  *handler.DisplayHandler
	TrackingHandler *handler.TrackingHandler
}

func BuildHandlers(cfg *config.Config, pool *pgxpool.Pool) *Handlers {
	userRepo := repository.NewUserRepository(pool)
	userCounterRepo := repository.NewUserCounterRepository(pool)
	counterRepo := repository.NewCounterRepository(pool)
	counterCategoryRepo := repository.NewCounterCategoryRepository(pool)
	categoryRepo := repository.NewCategoryRepository(pool)
	ticketRepo := repository.NewTicketRepository(pool)
	statsRepo := repository.NewStatsRepository(pool)

	userService := service.NewUserService(userRepo, userCounterRepo)
	adminService := service.NewAdminService(userRepo, userCounterRepo, counterRepo, counterCategoryRepo, categoryRepo, ticketRepo, statsRepo)
	staffService := service.NewStaffService(userRepo, userCounterRepo, counterRepo, counterCategoryRepo, ticketRepo, statsRepo, categoryRepo)
	kioskService := service.NewKioskService(categoryRepo, ticketRepo, statsRepo)
	displayService := service.NewDisplayService(statsRepo, categoryRepo, counterRepo)
	trackingService := service.NewTrackingService(ticketRepo, categoryRepo, counterRepo)

	middleware.InitAuth(&cfg.JWT)

	hub := websocket.NewHub()
	go hub.Run()

	authHandler := handler.NewAuthHandler(userService, &cfg.JWT)
	adminHandler := handler.NewAdminHandler(adminService, hub)
	staffHandler := handler.NewStaffHandler(staffService, hub)
	kioskHandler := handler.NewKioskHandler(kioskService, hub)
	displayHandler := handler.NewDisplayHandler(displayService)
	trackingHandler := handler.NewTrackingHandler(trackingService)

	return &Handlers{
		Hub:             hub,
		AuthHandler:     authHandler,
		AdminHandler:    adminHandler,
		StaffHandler:    staffHandler,
		KioskHandler:    kioskHandler,
		DisplayHandler:  displayHandler,
		TrackingHandler: trackingHandler,
	}
}
