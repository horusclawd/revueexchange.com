package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/revueexchange/api/internal/config"
	authmware "github.com/revueexchange/api/internal/middleware"
	"github.com/revueexchange/api/internal/service"
)

// Handler holds all handlers
type Handler struct {
	AuthService   *service.AuthService
	UserService   *service.UserService
	BountyService *service.BountyService
	PointsService *service.PointsService
}

// SetupRouter sets up the HTTP router
func SetupRouter(services *service.Services, cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// CORS
	r.Use(middleware.AllowContentType("application/json"))

	// Create handler
	h := &Handler{
		AuthService:   services.AuthService,
		UserService:   services.UserService,
		BountyService: services.BountyService,
		PointsService: services.PointsService,
	}

	// Public routes
	r.Group(func(r chi.Router) {
		r.Post("/api/v1/auth/register", h.Register)
		r.Post("/api/v1/auth/login", h.Login)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(authmware.AuthMiddleware(cfg.JWTSecret))
		r.Get("/api/v1/auth/me", h.Me)
		r.Put("/api/v1/users/{id}", h.UpdateUser)

		// Bounties
		r.Get("/api/v1/bounties", h.ListBounties)
		r.Post("/api/v1/bounties", h.CreateBounty)
		r.Get("/api/v1/bounties/{id}", h.GetBounty)

		// Points
		r.Get("/api/v1/points/balance", h.GetBalance)
		r.Get("/api/v1/points/transactions", h.GetTransactions)
	})

	// Health check
	r.Get("/health", h.HealthCheck)

	return r
}
