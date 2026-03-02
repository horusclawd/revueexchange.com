package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/revueexchange/api/internal/model"
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
	bounties, err := h.BountyService.GetBounties(r.Context(), 20, 0)
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
