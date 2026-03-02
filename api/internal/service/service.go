package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/revueexchange/api/internal/config"
	"github.com/revueexchange/api/internal/model"
	"github.com/revueexchange/api/internal/repository"
)

// Services holds all services
type Services struct {
	UserService    *UserService
	AuthService   *AuthService
	BountyService *BountyService
	PointsService *PointsService
}

// NewServices creates all services
func NewServices(repo *repository.Repository, cfg *config.Config) *Services {
	return &Services{
		UserService:    NewUserService(repo),
		AuthService:   NewAuthService(repo, cfg),
		BountyService: NewBountyService(repo),
		PointsService: NewPointsService(repo),
	}
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

// BountyService handles bounty operations
type BountyService struct {
	repo *repository.Repository
}

func NewBountyService(repo *repository.Repository) *BountyService {
	return &BountyService{repo: repo}
}

func (s *BountyService) GetBounties(ctx context.Context, limit, offset int) ([]model.Bounty, error) {
	return s.repo.GetBounties(ctx, limit, offset)
}

func (s *BountyService) GetBountyByID(ctx context.Context, id uuid.UUID) (*model.Bounty, error) {
	return s.repo.GetBountyByID(ctx, id)
}

func (s *BountyService) CreateBounty(ctx context.Context, bounty *model.Bounty) error {
	return s.repo.CreateBounty(ctx, bounty)
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

// Service errors
var (
	ErrUserAlreadyExists = &ServiceError{Message: "user already exists"}
	ErrInvalidCredentials = &ServiceError{Message: "invalid credentials"}
)

type ServiceError struct {
	Message string
}

func (e *ServiceError) Error() string {
	return e.Message
}
