package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/revueexchange/api/internal/model"
	"github.com/revueexchange/api/internal/repository"
	"github.com/rs/zerolog/log"
)

// Response is a standard API response
type Response struct {
	Data    interface{} `json:"data,omitempty"`
	Error   string     `json:"error,omitempty"`
	Message string     `json:"message,omitempty"`
}

// TokenResponse represents an auth response with token
type TokenResponse struct {
	User  *model.User `json:"user"`
	Token string      `json:"token"`
}

// HealthCheck handles GET /health
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Message: "OK"})
}

// Register handles POST /api/v1/auth/register
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Username == "" || req.Password == "" {
		http.Error(w, "email, username, and password are required", http.StatusBadRequest)
		return
	}

	user, err := h.AuthService.Register(r.Context(), req.Email, req.Username, req.Password)
	if err != nil {
		log.Error().Err(err).Msg("register failed")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate token
	token, err := h.AuthService.GenerateToken(user.ID)
	if err != nil {
		log.Error().Err(err).Msg("token generation failed")
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Data: TokenResponse{User: user, Token: token}})
}

// Login handles POST /api/v1/auth/login
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.AuthService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate token
	token, err := h.AuthService.GenerateToken(user.ID)
	if err != nil {
		log.Error().Err(err).Msg("token generation failed")
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Data: TokenResponse{User: user, Token: token}})
}

// Me handles GET /api/v1/auth/me
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(uuid.UUID)

	user, err := h.UserService.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Data: user})
}

// UpdateUser handles PUT /api/v1/users/{id}
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	// Verify user is updating themselves
	requestingUserID := r.Context().Value("user_id").(uuid.UUID)
	if userID != requestingUserID {
		http.Error(w, "unauthorized", http.StatusForbidden)
		return
	}

	var req struct {
		DisplayName *string `json:"display_name"`
		Bio         *string `json:"bio"`
		AvatarURL   *string `json:"avatar_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.UserService.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	if req.DisplayName != nil {
		user.DisplayName = *req.DisplayName
	}
	if req.Bio != nil {
		user.Bio = req.Bio
	}
	if req.AvatarURL != nil {
		user.AvatarURL = req.AvatarURL
	}

	if err := h.UserService.UpdateUser(r.Context(), user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Data: user})
}

// ListBounties handles GET /api/v1/bounties
func (h *Handler) ListBounties(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	limit := 20
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		fmt.Sscanf(o, "%d", &offset)
	}

	filters := repository.BountyFilters{
		Status:    r.URL.Query().Get("status"),
		Genre:     r.URL.Query().Get("genre"),
		Type:      r.URL.Query().Get("type"),
	}

	if minPoints := r.URL.Query().Get("min_points"); minPoints != "" {
		fmt.Sscanf(minPoints, "%d", &filters.MinPoints)
	}
	if maxPoints := r.URL.Query().Get("max_points"); maxPoints != "" {
		fmt.Sscanf(maxPoints, "%d", &filters.MaxPoints)
	}

	bounties, err := h.BountyService.GetBounties(r.Context(), filters, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Data: bounties})
}

// GetBounty handles GET /api/v1/bounties/{id}
func (h *Handler) GetBounty(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid bounty id", http.StatusBadRequest)
		return
	}

	bounty, err := h.BountyService.GetBountyByID(r.Context(), id)
	if err != nil {
		http.Error(w, "bounty not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Data: bounty})
}

// CreateBounty handles POST /api/v1/bounties
func (h *Handler) CreateBounty(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(uuid.UUID)

	var req struct {
		ProductID    uuid.UUID `json:"product_id"`
		BountyPoints int       `json:"bounty_points"`
		BountyCash   *float64  `json:"bounty_cash"`
		Requirements *string   `json:"requirements"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bounty := &model.Bounty{
		ID:            uuid.New(),
		UserID:        userID,
		ProductID:     req.ProductID,
		BountyPoints:  req.BountyPoints,
		BountyCash:    req.BountyCash,
		Status:        "open",
		Requirements:  req.Requirements,
	}

	if err := h.BountyService.CreateBounty(r.Context(), bounty); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Data: bounty})
}

// ClaimBounty handles POST /api/v1/bounties/{id}/claim
func (h *Handler) ClaimBounty(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid bounty id", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(uuid.UUID)

	bounty, err := h.BountyService.ClaimBounty(r.Context(), id, userID)
	if err != nil {
		if err.Error() == "cannot claim your own bounty" {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Data: bounty})
}

// CancelBounty handles DELETE /api/v1/bounties/{id}
func (h *Handler) CancelBounty(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid bounty id", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(uuid.UUID)

	err = h.BountyService.CancelBounty(r.Context(), id, userID)
	if err != nil {
		if err.Error() == "bounty is not available" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Message: "Bounty cancelled"})
}

// CreateReview handles POST /api/v1/reviews
func (h *Handler) CreateReview(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(uuid.UUID)

	var req struct {
		BountyID uuid.UUID `json:"bounty_id"`
		Rating   int       `json:"rating"`
		Title    *string   `json:"title"`
		Content  *string   `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Calculate word count from content
	var wordCount *int
	if req.Content != nil {
		count := len(*req.Content)
		wordCount = &count
	}

	review := &model.Review{
		ID:         uuid.New(),
		BountyID:   req.BountyID,
		ReviewerID: userID,
		Rating:     req.Rating,
		Title:      req.Title,
		Content:    req.Content,
		WordCount:  wordCount,
		Status:     "draft",
	}

	if err := h.ReviewService.CreateReview(r.Context(), review); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Response{Data: review})
}

// GetReview handles GET /api/v1/reviews/{id}
func (h *Handler) GetReview(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid review id", http.StatusBadRequest)
		return
	}

	review, err := h.ReviewService.GetReviewByID(r.Context(), id)
	if err != nil {
		http.Error(w, "review not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Data: review})
}

// UpdateReview handles PUT /api/v1/reviews/{id}
func (h *Handler) UpdateReview(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid review id", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(uuid.UUID)

	// Get existing review
	review, err := h.ReviewService.GetReviewByID(r.Context(), id)
	if err != nil {
		http.Error(w, "review not found", http.StatusNotFound)
		return
	}

	// Verify ownership
	if review.ReviewerID != userID {
		http.Error(w, "unauthorized", http.StatusForbidden)
		return
	}

	var req struct {
		Rating  *int    `json:"rating"`
		Title   *string `json:"title"`
		Content *string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Rating != nil {
		review.Rating = *req.Rating
	}
	if req.Title != nil {
		review.Title = req.Title
	}
	if req.Content != nil {
		review.Content = req.Content
		count := len(*req.Content)
		review.WordCount = &count
	}

	review.UpdatedAt = time.Now()

	if err := h.ReviewService.UpdateReview(r.Context(), review); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Data: review})
}

// SubmitReview handles POST /api/v1/reviews/{id}/submit
func (h *Handler) SubmitReview(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid review id", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(uuid.UUID)

	// Get bounty to find points
	bounty, err := h.BountyService.GetBountyByID(r.Context(), id)
	if err != nil {
		http.Error(w, "bounty not found", http.StatusNotFound)
		return
	}

	review, err := h.ReviewService.SubmitReview(r.Context(), id, userID, bounty.BountyPoints)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Data: review})
}

// GetBalance handles GET /api/v1/points/balance
func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(uuid.UUID)

	balance, err := h.PointsService.GetBalance(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Data: map[string]int{"balance": balance}})
}

// GetTransactions handles GET /api/v1/points/transactions
func (h *Handler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(uuid.UUID)

	txs, err := h.PointsService.GetTransactions(r.Context(), userID, 20, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Data: txs})
}

// TransferPoints handles POST /api/v1/points/transfer
func (h *Handler) TransferPoints(w http.ResponseWriter, r *http.Request) {
	fromUserID := r.Context().Value("user_id").(uuid.UUID)

	var req struct {
		ToUserID uuid.UUID `json:"to_user_id"`
		Amount   int       `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 {
		http.Error(w, "amount must be positive", http.StatusBadRequest)
		return
	}

	err := h.PointsService.TransferPoints(r.Context(), fromUserID, req.ToUserID, req.Amount)
	if err != nil {
		if err.Error() == "insufficient points" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Message: "Points transferred successfully"})
}

// CreateProduct handles POST /api/v1/products
func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(uuid.UUID)

	var req struct {
		Type          string  `json:"type"`
		Title         string  `json:"title"`
		Description   *string `json:"description"`
		URL           *string `json:"url"`
		CoverImageURL *string `json:"cover_image_url"`
		Genre         *string `json:"genre"`
		WordCount     *int    `json:"word_count"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Type == "" || req.Title == "" {
		http.Error(w, "type and title are required", http.StatusBadRequest)
		return
	}

	product := &model.Product{
		ID:            uuid.New(),
		UserID:        userID,
		Type:          req.Type,
		Title:         req.Title,
		Description:   req.Description,
		URL:           req.URL,
		CoverImageURL: req.CoverImageURL,
		Genre:         req.Genre,
		WordCount:     req.WordCount,
	}

	if err := h.ProductService.CreateProduct(r.Context(), product); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Response{Data: product})
}

// GetProduct handles GET /api/v1/products/{id}
func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	product, err := h.ProductService.GetProductByID(r.Context(), id)
	if err != nil {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Data: product})
}

// GetUserProducts handles GET /api/v1/users/{id}/products
func (h *Handler) GetUserProducts(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	products, err := h.ProductService.GetProductsByUserID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Data: products})
}

// UpdateProduct handles PUT /api/v1/products/{id}
func (h *Handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(uuid.UUID)

	// Get existing product
	product, err := h.ProductService.GetProductByID(r.Context(), id)
	if err != nil {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}

	// Verify ownership
	if product.UserID != userID {
		http.Error(w, "unauthorized", http.StatusForbidden)
		return
	}

	var req struct {
		Title         *string `json:"title"`
		Description   *string `json:"description"`
		URL           *string `json:"url"`
		CoverImageURL *string `json:"cover_image_url"`
		Genre         *string `json:"genre"`
		WordCount     *int    `json:"word_count"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Title != nil {
		product.Title = *req.Title
	}
	if req.Description != nil {
		product.Description = req.Description
	}
	if req.URL != nil {
		product.URL = req.URL
	}
	if req.CoverImageURL != nil {
		product.CoverImageURL = req.CoverImageURL
	}
	if req.Genre != nil {
		product.Genre = req.Genre
	}
	if req.WordCount != nil {
		product.WordCount = req.WordCount
	}

	if err := h.ProductService.UpdateProduct(r.Context(), product); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Data: product})
}

// DeleteProduct handles DELETE /api/v1/products/{id}
func (h *Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(uuid.UUID)

	// Get existing product
	product, err := h.ProductService.GetProductByID(r.Context(), id)
	if err != nil {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}

	// Verify ownership
	if product.UserID != userID {
		http.Error(w, "unauthorized", http.StatusForbidden)
		return
	}

	if err := h.ProductService.DeleteProduct(r.Context(), id, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Message: "Product deleted"})
}

// CheckoutRequest represents a checkout request
type CheckoutRequest struct {
	AmountCents int `json:"amount_cents"`
}

// CreateCheckoutSession handles POST /api/v1/payments/checkout
func (h *Handler) CreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
	var req CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(uuid.UUID)

	result, err := h.PaymentService.CreateCheckoutSession(r.Context(), userID, req.AmountCents)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Data: map[string]interface{}{
			"payment_id":   result.Payment.ID,
			"checkout_url": result.SessionURL,
			"amount_cents": result.Payment.AmountCents,
			"points_award": result.PointsAward,
		},
	})
}

// GetPaymentHistory handles GET /api/v1/payments/history
func (h *Handler) GetPaymentHistory(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(uuid.UUID)

	payments, err := h.PaymentService.GetPaymentHistory(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Data: payments})
}

// HandlePaymentWebhook handles POST /api/v1/payments/webhook
func (h *Handler) HandlePaymentWebhook(w http.ResponseWriter, r *http.Request) {
	// In production, verify Stripe webhook signature
	// For now, parse the event from the request body

	var webhookEvent struct {
		Type      string `json:"type"`
		SessionID string `json:"session_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&webhookEvent); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.PaymentService.HandleWebhook(r.Context(), webhookEvent.Type, webhookEvent.SessionID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Message: "Webhook processed"})
}

// FollowUser handles POST /api/v1/social/follow/{id}
func (h *Handler) FollowUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	followingID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(uuid.UUID)

	if err := h.SocialService.FollowUser(r.Context(), userID, followingID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Message: "Followed user"})
}

// UnfollowUser handles DELETE /api/v1/social/follow/{id}
func (h *Handler) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	followingID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(uuid.UUID)

	if err := h.SocialService.UnfollowUser(r.Context(), userID, followingID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Message: "Unfollowed user"})
}

// GetFollowers handles GET /api/v1/social/followers/{id}
func (h *Handler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	users, err := h.SocialService.GetFollowers(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Data: users})
}

// GetFollowing handles GET /api/v1/social/following/{id}
func (h *Handler) GetFollowing(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	users, err := h.SocialService.GetFollowing(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Data: users})
}

// GetActivityFeed handles GET /api/v1/social/feed
func (h *Handler) GetActivityFeed(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(uuid.UUID)

	activities, err := h.SocialService.GetActivityFeed(r.Context(), userID, 50, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Data: activities})
}

// CommentRequest represents a comment request
type CommentRequest struct {
	ReviewID string `json:"review_id"`
	Content   string `json:"content"`
	ParentID  string `json:"parent_id"`
}

// AddComment handles POST /api/v1/comments
func (h *Handler) AddComment(w http.ResponseWriter, r *http.Request) {
	var req CommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	reviewID, err := uuid.Parse(req.ReviewID)
	if err != nil {
		http.Error(w, "invalid review id", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(uuid.UUID)

	var parentID *uuid.UUID
	if req.ParentID != "" {
		pid, err := uuid.Parse(req.ParentID)
		if err != nil {
			http.Error(w, "invalid parent id", http.StatusBadRequest)
			return
		}
		parentID = &pid
	}

	comment, err := h.SocialService.AddComment(r.Context(), userID, reviewID, req.Content, parentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Data: comment})
}

// DeleteComment handles DELETE /api/v1/comments/{id}
func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	commentID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid comment id", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(uuid.UUID)

	if err := h.SocialService.DeleteComment(r.Context(), commentID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Message: "Comment deleted"})
}

// GetComments handles GET /api/v1/comments
func (h *Handler) GetComments(w http.ResponseWriter, r *http.Request) {
	reviewIDStr := r.URL.Query().Get("review_id")
	if reviewIDStr == "" {
		http.Error(w, "review_id is required", http.StatusBadRequest)
		return
	}

	reviewID, err := uuid.Parse(reviewIDStr)
	if err != nil {
		http.Error(w, "invalid review_id", http.StatusBadRequest)
		return
	}

	comments, err := h.SocialService.GetComments(r.Context(), reviewID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Data: comments})
}
