package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/revueexchange/api/internal/model"
	"github.com/revueexchange/api/internal/repository"
)

// GamificationService handles leaderboard and streak operations
type GamificationService struct {
	repo           *repository.Repository
	gamificationRepo *repository.GamificationRepository
}

// NewGamificationService creates a new gamification service
func NewGamificationService(repo *repository.Repository, gamificationRepo *repository.GamificationRepository) *GamificationService {
	return &GamificationService{
		repo:             repo,
		gamificationRepo: gamificationRepo,
	}
}

// GetLeaderboard gets the top users by points
func (s *GamificationService) GetLeaderboard(ctx context.Context, limit int) ([]model.LeaderboardEntry, error) {
	return s.gamificationRepo.GetLeaderboard(ctx, limit)
}

// UpdateLeaderboard updates a user's position on the leaderboard
func (s *GamificationService) UpdateLeaderboard(ctx context.Context, userID uuid.UUID) error {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	reviews, err := s.repo.GetReviewsByReviewer(ctx, userID)
	if err != nil {
		return err
	}

	reviewCount := 0
	for _, r := range reviews {
		if r.Status == "published" {
			reviewCount++
		}
	}

	return s.gamificationRepo.UpdateLeaderboard(ctx, userID.String(), user.Username, user.DisplayName, user.Points, reviewCount)
}

// GetStreak gets a user's current streak
func (s *GamificationService) GetStreak(ctx context.Context, userID uuid.UUID) (*model.Streak, error) {
	return s.gamificationRepo.GetStreak(ctx, userID.String())
}

// UpdateStreak updates a user's streak based on activity
func (s *GamificationService) UpdateStreak(ctx context.Context, userID uuid.UUID) (*model.Streak, error) {
	return s.gamificationRepo.UpdateStreak(ctx, userID.String())
}
