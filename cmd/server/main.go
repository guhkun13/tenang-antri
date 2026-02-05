package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"tenangantri/internal/config"
	"tenangantri/internal/migrate"
	"tenangantri/internal/server"
)

func main() {
	// Setup logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Check for migrate command
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "migrate":
			log.Info().Msg("Running migrations...")
			if err := migrate.Run(cfg.GetDatabaseURL(), "migrations"); err != nil {
				log.Fatal().Err(err).Msg("Migration failed")
			}
			log.Info().Msg("Migrations completed")
			return
		case "migrate-force":
			if len(os.Args) < 3 {
				log.Fatal().Msg("Usage: server migrate-force <version>")
			}
			versionStr := os.Args[2]
			var version int
			fmt.Sscanf(versionStr, "%d", &version)

			log.Info().Int("version", version).Msg("Forcing migration version...")
			if err := migrate.Force(cfg.GetDatabaseURL(), "migrations", version); err != nil {
				log.Fatal().Err(err).Msg("Force migration failed")
			}
			log.Info().Msg("Migration version forced")
			return
		}
	}

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Connect to database
	pool, err := pgxpool.New(context.Background(), cfg.GetDatabaseURL())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer pool.Close()

	// Test database connection
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal().Err(err).Msg("Failed to ping database")
	}
	log.Info().Msg("Connected to database")

	r := server.NewRouter(cfg, pool)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Info().Str("port", cfg.Server.Port).Msg("Starting server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exited")
}
