package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/revueexchange/api/internal/config"
	"github.com/revueexchange/api/internal/handler"
	"github.com/revueexchange/api/internal/service"
	"github.com/revueexchange/api/pkg/logger"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Setup logger
	log := logger.Setup(cfg.LogLevel)
	log.Info().Str("environment", cfg.Environment).Msg("Starting RevUExchange API")

	// Initialize repository (database connection)
	repo, err := config.InitDB(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer repo.Close()

	// Initialize services
	services := service.NewServices(repo, cfg)

	// Setup router
	router := handler.SetupRouter(services, cfg)

	// Create server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Info().Int("port", cfg.Port).Msg("Starting HTTP server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("HTTP server failed")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exited")
}
