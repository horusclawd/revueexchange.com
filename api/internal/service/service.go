package service

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
	"github.com/revueexchange/api/internal/config"
	"github.com/revueexchange/api/internal/model"
	"github.com/revueexchange/api/internal/repository"
)

// Services holds all services
type Services struct {
	UserService       *UserService
	AuthService      *AuthService
	ProductService   *ProductService
	BountyService    *BountyService
	ReviewService    *ReviewService
	PointsService    *PointsService
	PaymentService   *PaymentService
	SocialService    *SocialService
	BadgeService     *BadgeService
	GamificationService *GamificationService
	AnalyticsService *AnalyticsService
}

// NewServices creates all services
func NewServices(repo *repository.Repository, dynamoDB *dynamodb.Client, cfg *config.Config) *Services {
	services := &Services{
		UserService:    NewUserService(repo),
		AuthService:   NewAuthService(repo, cfg),
		ProductService: NewProductService(repo),
		BountyService: NewBountyService(repo),
		ReviewService: NewReviewService(repo),
		PointsService: NewPointsService(repo),
		PaymentService: NewPaymentService(repo, cfg),
		SocialService:  NewSocialService(repo),
	}

	// Initialize badge service if DynamoDB is available
	if dynamoDB != nil {
		badgeRepo := repository.NewBadgeRepository(dynamoDB)
		services.BadgeService = NewBadgeService(repo, badgeRepo)

		// Initialize gamification service
		gamificationRepo := repository.NewGamificationRepository(dynamoDB)
		services.GamificationService = NewGamificationService(repo, gamificationRepo)
	}

	// Always initialize analytics service (uses PostgreSQL)
	services.AnalyticsService = NewAnalyticsService(repo)

	return services
}

// UserService handles user operations
type UserService struct {
	repo *repository.Repository
}

func NewUserService(repo *repository.Repository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}

func (s *UserService) UpdateUser(ctx context.Context, user *model.User) error {
	return s.repo.UpdateUser(ctx, user)
}

// AuthService handles authentication
type AuthService struct {
	repo *repository.Repository
	cfg  *config.Config
}

func NewAuthService(repo *repository.Repository, cfg *config.Config) *AuthService {
	return &AuthService{repo: repo, cfg: cfg}
}

func (s *AuthService) Register(ctx context.Context, email, username, password string) (*model.User, error) {
	// Check if user exists
	existing, _ := s.repo.GetUserByEmail(ctx, email)
	if existing != nil {
		return nil, ErrUserAlreadyExists
	}

	// Hash password
	hash, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &model.User{
		ID:              uuid.New(),
		Email:           email,
		PasswordHash:    hash,
		Username:        username,
		DisplayName:    username,
		Points:          100, // Welcome bonus
		ReputationScore: 0,
		SubscriptionTier: "free",
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*model.User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !CheckPassword(password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

// GenerateToken generates a JWT token for a user
func (s *AuthService) GenerateToken(userID uuid.UUID) (string, error) {
	return GenerateToken(userID, s.cfg.JWTSecret)
}

// ProductService handles product operations
type ProductService struct {
	repo *repository.Repository
}

func NewProductService(repo *repository.Repository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) CreateProduct(ctx context.Context, product *model.Product) error {
	return s.repo.CreateProduct(ctx, product)
}

func (s *ProductService) GetProductByID(ctx context.Context, id uuid.UUID) (*model.Product, error) {
	return s.repo.GetProductByID(ctx, id)
}

func (s *ProductService) GetProductsByUserID(ctx context.Context, userID uuid.UUID) ([]model.Product, error) {
	return s.repo.GetProductsByUserID(ctx, userID)
}

func (s *ProductService) UpdateProduct(ctx context.Context, product *model.Product) error {
	return s.repo.UpdateProduct(ctx, product)
}

func (s *ProductService) DeleteProduct(ctx context.Context, id, userID uuid.UUID) error {
	return s.repo.DeleteProduct(ctx, id, userID)
}

// BountyService handles bounty operations
type BountyService struct {
	repo *repository.Repository
}

func NewBountyService(repo *repository.Repository) *BountyService {
	return &BountyService{repo: repo}
}

func (s *BountyService) GetBounties(ctx context.Context, filters repository.BountyFilters, limit, offset int) ([]model.Bounty, error) {
	return s.repo.GetBounties(ctx, filters, limit, offset)
}

func (s *BountyService) GetBountyByID(ctx context.Context, id uuid.UUID) (*model.Bounty, error) {
	return s.repo.GetBountyByID(ctx, id)
}

func (s *BountyService) CreateBounty(ctx context.Context, bounty *model.Bounty) error {
	return s.repo.CreateBounty(ctx, bounty)
}

// ClaimBounty claims a bounty for review
func (s *BountyService) ClaimBounty(ctx context.Context, bountyID, reviewerID uuid.UUID) (*model.Bounty, error) {
	bounty, err := s.repo.GetBountyByID(ctx, bountyID)
	if err != nil {
		return nil, err
	}

	// Anti-swap: Check if user has already reviewed this author's products
	// (simplified version - in production would check review history)
	if bounty.UserID == reviewerID {
		return nil, ErrCannotClaimOwnBounty
	}

	if bounty.Status != "open" {
		return nil, ErrBountyNotAvailable
	}

	now := time.Now()
	bounty.Status = "claimed"
	bounty.ClaimedBy = &reviewerID
	bounty.ClaimedAt = &now
	bounty.UpdatedAt = now

	if err := s.repo.UpdateBounty(ctx, bounty); err != nil {
		return nil, err
	}

	return bounty, nil
}

// CancelBounty cancels an open bounty
func (s *BountyService) CancelBounty(ctx context.Context, bountyID, userID uuid.UUID) error {
	bounty, err := s.repo.GetBountyByID(ctx, bountyID)
	if err != nil {
		return err
	}

	// Only the creator can cancel
	if bounty.UserID != userID {
		return ErrUnauthorized
	}

	if bounty.Status != "open" {
		return ErrBountyNotAvailable
	}

	bounty.Status = "cancelled"
	bounty.UpdatedAt = time.Now()

	return s.repo.UpdateBounty(ctx, bounty)
}

// Service errors
var (
	ErrCannotClaimOwnBounty = &ServiceError{Message: "cannot claim your own bounty"}
	ErrBountyNotAvailable   = &ServiceError{Message: "bounty is not available"}
	ErrUnauthorized          = &ServiceError{Message: "unauthorized"}
)

// ReviewService handles review operations
type ReviewService struct {
	repo *repository.Repository
}

func NewReviewService(repo *repository.Repository) *ReviewService {
	return &ReviewService{repo: repo}
}

func (s *ReviewService) CreateReview(ctx context.Context, review *model.Review) error {
	return s.repo.CreateReview(ctx, review)
}

func (s *ReviewService) GetReviewByID(ctx context.Context, id uuid.UUID) (*model.Review, error) {
	return s.repo.GetReviewByID(ctx, id)
}

func (s *ReviewService) GetReviewByBountyID(ctx context.Context, bountyID uuid.UUID) (*model.Review, error) {
	return s.repo.GetReviewByBountyID(ctx, bountyID)
}

func (s *ReviewService) GetReviewsByReviewer(ctx context.Context, reviewerID uuid.UUID) ([]model.Review, error) {
	return s.repo.GetReviewsByReviewer(ctx, reviewerID)
}

func (s *ReviewService) UpdateReview(ctx context.Context, review *model.Review) error {
	return s.repo.UpdateReview(ctx, review)
}

// SubmitReview submits a review and awards points
func (s *ReviewService) SubmitReview(ctx context.Context, reviewID, reviewerID uuid.UUID, bountyPoints int) (*model.Review, error) {
	review, err := s.repo.GetReviewByID(ctx, reviewID)
	if err != nil {
		return nil, err
	}

	// Verify reviewer owns the review
	if review.ReviewerID != reviewerID {
		return nil, ErrUnauthorized
	}

	// Can only submit draft reviews
	if review.Status != "draft" {
		return nil, &ServiceError{Message: "review already submitted"}
	}

	// Validate word count (minimum 10 words)
	if review.WordCount != nil && *review.WordCount < 10 {
		return nil, &ServiceError{Message: "review must be at least 10 words"}
	}

	// Validate rating (1-5)
	if review.Rating < 1 || review.Rating > 5 {
		return nil, &ServiceError{Message: "rating must be between 1 and 5"}
	}

	review.Status = "submitted"
	review.UpdatedAt = time.Now()

	if err := s.repo.UpdateReview(ctx, review); err != nil {
		return nil, err
	}

	return review, nil
}

// PointsService handles points operations
type PointsService struct {
	repo *repository.Repository
}

func NewPointsService(repo *repository.Repository) *PointsService {
	return &PointsService{repo: repo}
}

func (s *PointsService) GetBalance(ctx context.Context, userID uuid.UUID) (int, error) {
	return s.repo.GetUserPoints(ctx, userID)
}

func (s *PointsService) GetTransactions(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.PointTransaction, error) {
	return s.repo.GetPointTransactions(ctx, userID, limit, offset)
}

func (s *PointsService) AwardPoints(ctx context.Context, userID uuid.UUID, amount int, description string) error {
	tx := &model.PointTransaction{
		ID:          uuid.New(),
		UserID:      userID,
		Amount:      amount,
		Type:        "earned",
		Description: &description,
	}

	if err := s.repo.CreatePointTransaction(ctx, tx); err != nil {
		return err
	}

	return s.repo.UpdateUserPoints(ctx, userID, amount)
}

func (s *PointsService) DeductPoints(ctx context.Context, userID uuid.UUID, amount int, description string) error {
	tx := &model.PointTransaction{
		ID:          uuid.New(),
		UserID:      userID,
		Amount:      -amount,
		Type:        "spent",
		Description: &description,
	}

	if err := s.repo.CreatePointTransaction(ctx, tx); err != nil {
		return err
	}

	return s.repo.UpdateUserPoints(ctx, userID, -amount)
}

// TransferPoints transfers points from one user to another
func (s *PointsService) TransferPoints(ctx context.Context, fromUserID, toUserID uuid.UUID, amount int) error {
	if amount <= 0 {
		return &ServiceError{Message: "amount must be positive"}
	}

	// Check balance
	balance, err := s.repo.GetUserPoints(ctx, fromUserID)
	if err != nil {
		return err
	}

	if balance < amount {
		return ErrInsufficientPoints
	}

	// Deduct from sender
	tx := &model.PointTransaction{
		ID:            uuid.New(),
		UserID:        fromUserID,
		Amount:        -amount,
		Type:          "transferred",
		ReferenceType: func() *string { s := "transfer"; return &s }(),
		ReferenceID:   func() *uuid.UUID { u := toUserID; return &u }(),
		Description:   func() *string { d := "Transfer to user"; return &d }(),
	}

	if err := s.repo.CreatePointTransaction(ctx, tx); err != nil {
		return err
	}

	if err := s.repo.UpdateUserPoints(ctx, fromUserID, -amount); err != nil {
		return err
	}

	// Add to receiver
	tx.ID = uuid.New()
	tx.UserID = toUserID
	tx.Amount = amount
	tx.Type = "received"
	tx.Description = func() *string { d := "Received transfer"; return &d }()

	if err := s.repo.CreatePointTransaction(ctx, tx); err != nil {
		return err
	}

	return s.repo.UpdateUserPoints(ctx, toUserID, amount)
}

// Service errors
var (
	ErrUserAlreadyExists  = &ServiceError{Message: "user already exists"}
	ErrInvalidCredentials = &ServiceError{Message: "invalid credentials"}
	ErrInsufficientPoints = &ServiceError{Message: "insufficient points"}
)

type ServiceError struct {
	Message string
}

func (e *ServiceError) Error() string {
	return e.Message
}

// PaymentService handles payment operations
type PaymentService struct {
	repo *repository.Repository
	cfg  *config.Config
}

func NewPaymentService(repo *repository.Repository, cfg *config.Config) *PaymentService {
	return &PaymentService{repo: repo, cfg: cfg}
}

// PointsConversionRate is the number of points per dollar (100 points = $1)
const PointsConversionRate = 100

// CheckoutResult contains the checkout session result
type CheckoutResult struct {
	Payment     *model.Payment
	SessionURL  string
	PointsAward int
}

// CreateCheckoutSession creates a Stripe checkout session for purchasing points
func (s *PaymentService) CreateCheckoutSession(ctx context.Context, userID uuid.UUID, amountCents int) (*CheckoutResult, error) {
	if amountCents < 100 {
		return nil, &ServiceError{Message: "minimum purchase is $1.00"}
	}

	// Calculate points to award
	points := (amountCents / 100) * PointsConversionRate

	// Create payment record
	payment := &model.Payment{
		ID:          uuid.New(),
		UserID:      userID,
		AmountCents: amountCents,
		Currency:    "usd",
		Type:        func() *string { t := "points_purchase"; return &t }(),
		Status:      "pending",
		CreatedAt:   time.Now(),
	}

	if err := s.repo.CreatePayment(ctx, payment); err != nil {
		return nil, err
	}

	// In production, this would create a real Stripe checkout session
	// For now, return a mock session URL
	sessionURL := "/checkout/success?session_id=" + payment.ID.String()

	return &CheckoutResult{
		Payment:     payment,
		SessionURL:  sessionURL,
		PointsAward: points,
	}, nil
}

// HandleWebhook processes Stripe webhook events
func (s *PaymentService) HandleWebhook(ctx context.Context, eventType string, sessionID string) error {
	// Get payment by session ID
	payment, err := s.repo.GetPaymentBySessionID(ctx, sessionID)
	if err != nil {
		return err
	}

	switch eventType {
	case "checkout.session.completed":
		// Update payment status
		if err := s.repo.UpdatePaymentStatus(ctx, payment.ID, "completed"); err != nil {
			return err
		}

		// Convert cash to points
		points := (payment.AmountCents / 100) * PointsConversionRate
		if err := s.repo.UpdateUserPoints(ctx, payment.UserID, points); err != nil {
			return err
		}

		// Create point transaction
		tx := &model.PointTransaction{
			ID:          uuid.New(),
			UserID:      payment.UserID,
			Amount:      points,
			Type:        "earned",
			ReferenceType: func() *string { t := "payment"; return &t }(),
			ReferenceID: func() *uuid.UUID { u := payment.ID; return &u }(),
			Description: func() *string { d := "Points purchase"; return &d }(),
			CreatedAt:   time.Now(),
		}
		if err := s.repo.CreatePointTransaction(ctx, tx); err != nil {
			return err
		}

	case "checkout.session.expired":
		if err := s.repo.UpdatePaymentStatus(ctx, payment.ID, "expired"); err != nil {
			return err
		}

	case "payment_intent.payment_failed":
		if err := s.repo.UpdatePaymentStatus(ctx, payment.ID, "failed"); err != nil {
			return err
		}
	}

	return nil
}

// GetPaymentHistory gets payment history for a user
func (s *PaymentService) GetPaymentHistory(ctx context.Context, userID uuid.UUID) ([]model.Payment, error) {
	return s.repo.GetPaymentsByUserID(ctx, userID)
}

// SocialService handles social operations
type SocialService struct {
	repo *repository.Repository
}

func NewSocialService(repo *repository.Repository) *SocialService {
	return &SocialService{repo: repo}
}

// FollowUser follows a user
func (s *SocialService) FollowUser(ctx context.Context, followerID, followingID uuid.UUID) error {
	if followerID == followingID {
		return &ServiceError{Message: "cannot follow yourself"}
	}

	// Check if already following
	exists, err := s.repo.IsFollowing(ctx, followerID, followingID)
	if err != nil {
		return err
	}
	if exists {
		return &ServiceError{Message: "already following this user"}
	}

	follow := &model.Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
		CreatedAt:   time.Now(),
	}

	if err := s.repo.CreateFollow(ctx, follow); err != nil {
		return err
	}

	// Create activity
	ref := followingID.String()
	activity := &model.Activity{
		ID:        uuid.New(),
		UserID:    followerID,
		Type:      "follow",
		Reference: &ref,
		CreatedAt: time.Now(),
	}
	return s.repo.CreateActivity(ctx, activity)
}

// UnfollowUser unfollows a user
func (s *SocialService) UnfollowUser(ctx context.Context, followerID, followingID uuid.UUID) error {
	return s.repo.DeleteFollow(ctx, followerID, followingID)
}

// GetFollowers gets followers of a user
func (s *SocialService) GetFollowers(ctx context.Context, userID uuid.UUID) ([]model.User, error) {
	return s.repo.GetFollowers(ctx, userID)
}

// GetFollowing gets users that a user is following
func (s *SocialService) GetFollowing(ctx context.Context, userID uuid.UUID) ([]model.User, error) {
	return s.repo.GetFollowing(ctx, userID)
}

// IsFollowing checks if a user is following another
func (s *SocialService) IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error) {
	return s.repo.IsFollowing(ctx, followerID, followingID)
}

// AddComment adds a comment to a review
func (s *SocialService) AddComment(ctx context.Context, userID, reviewID uuid.UUID, content string, parentID *uuid.UUID) (*model.Comment, error) {
	comment := &model.Comment{
		ID:        uuid.New(),
		UserID:    userID,
		ReviewID:  reviewID,
		ParentID:  parentID,
		Content:   content,
		CreatedAt: time.Now(),
	}

	if err := s.repo.CreateComment(ctx, comment); err != nil {
		return nil, err
	}

	// Create activity
	ref := reviewID.String()
	activity := &model.Activity{
		ID:        uuid.New(),
		UserID:    userID,
		Type:      "comment",
		Reference: &ref,
		CreatedAt: time.Now(),
	}
	if err := s.repo.CreateActivity(ctx, activity); err != nil {
		return nil, err
	}

	return comment, nil
}

// GetComments gets comments for a review
func (s *SocialService) GetComments(ctx context.Context, reviewID uuid.UUID) ([]model.Comment, error) {
	return s.repo.GetCommentsByReviewID(ctx, reviewID)
}

// DeleteComment deletes a comment
func (s *SocialService) DeleteComment(ctx context.Context, commentID, userID uuid.UUID) error {
	return s.repo.DeleteComment(ctx, commentID, userID)
}

// GetActivityFeed gets activity feed for a user
func (s *SocialService) GetActivityFeed(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Activity, error) {
	return s.repo.GetActivityFeed(ctx, userID, limit, offset)
}
