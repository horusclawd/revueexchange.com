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
	AuthService        *service.AuthService
	UserService        *service.UserService
	ProductService    *service.ProductService
	BountyService     *service.BountyService
	ReviewService     *service.ReviewService
	PointsService     *service.PointsService
	PaymentService    *service.PaymentService
	SocialService     *service.SocialService
	BadgeService      *service.BadgeService
	GamificationService *service.GamificationService
	AnalyticsService  *service.AnalyticsService
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
		AuthService:          services.AuthService,
		UserService:          services.UserService,
		ProductService:      services.ProductService,
		BountyService:       services.BountyService,
		ReviewService:       services.ReviewService,
		PointsService:       services.PointsService,
		PaymentService:      services.PaymentService,
		SocialService:       services.SocialService,
		BadgeService:        services.BadgeService,
		GamificationService:  services.GamificationService,
		AnalyticsService:    services.AnalyticsService,
	}

	// Public routes
	r.Group(func(r chi.Router) {
		r.Post("/api/v1/auth/register", h.Register)
		r.Post("/api/v1/auth/login", h.Login)
		r.Get("/api/v1/products/{id}", h.GetProduct)
		r.Get("/api/v1/users/{id}/products", h.GetUserProducts)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(authmware.AuthMiddleware(cfg.JWTSecret))
		r.Get("/api/v1/auth/me", h.Me)
		r.Put("/api/v1/users/{id}", h.UpdateUser)

		// Products
		r.Post("/api/v1/products", h.CreateProduct)
		r.Put("/api/v1/products/{id}", h.UpdateProduct)
		r.Delete("/api/v1/products/{id}", h.DeleteProduct)

		// Bounties
		r.Get("/api/v1/bounties", h.ListBounties)
		r.Post("/api/v1/bounties", h.CreateBounty)
		r.Get("/api/v1/bounties/{id}", h.GetBounty)
		r.Post("/api/v1/bounties/{id}/claim", h.ClaimBounty)
		r.Delete("/api/v1/bounties/{id}", h.CancelBounty)

		// Reviews
		r.Post("/api/v1/reviews", h.CreateReview)
		r.Get("/api/v1/reviews/{id}", h.GetReview)
		r.Put("/api/v1/reviews/{id}", h.UpdateReview)
		r.Post("/api/v1/reviews/{id}/submit", h.SubmitReview)

		// Points
		r.Get("/api/v1/points/balance", h.GetBalance)
		r.Get("/api/v1/points/transactions", h.GetTransactions)
		r.Post("/api/v1/points/transfer", h.TransferPoints)

		// Payments
		r.Post("/api/v1/payments/checkout", h.CreateCheckoutSession)
		r.Get("/api/v1/payments/history", h.GetPaymentHistory)

		// Social
		r.Post("/api/v1/social/follow/{id}", h.FollowUser)
		r.Delete("/api/v1/social/follow/{id}", h.UnfollowUser)
		r.Get("/api/v1/social/followers/{id}", h.GetFollowers)
		r.Get("/api/v1/social/following/{id}", h.GetFollowing)
		r.Get("/api/v1/social/feed", h.GetActivityFeed)

		// Comments
		r.Post("/api/v1/comments", h.AddComment)
		r.Get("/api/v1/comments", h.GetComments)
		r.Delete("/api/v1/comments/{id}", h.DeleteComment)

		// Badges
		r.Get("/api/v1/badges", h.GetUserBadges)
		r.Post("/api/v1/badges/check", h.CheckAndAwardBadges)

		// Gamification
		r.Get("/api/v1/gamification/leaderboard", h.GetLeaderboard)
		r.Get("/api/v1/gamification/streak", h.GetStreak)
		r.Post("/api/v1/gamification/streak/update", h.UpdateStreak)

		// Analytics
		r.Get("/api/v1/analytics/overview", h.GetAnalyticsOverview)
		r.Get("/api/v1/analytics/bounties", h.GetBountyMetrics)
		r.Get("/api/v1/analytics/reviews", h.GetReviewMetrics)
		r.Get("/api/v1/analytics/revenue", h.GetRevenueStats)
		r.Get("/api/v1/analytics/activity", h.GetUserActivity)
		r.Get("/api/v1/analytics/user", h.GetUserAnalytics)
	})

	// Webhook (not protected by auth)
	r.Post("/api/v1/payments/webhook", h.HandlePaymentWebhook)

	// Health check
	r.Get("/health", h.HealthCheck)

	return r
}
