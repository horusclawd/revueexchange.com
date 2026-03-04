package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/revueexchange/api/internal/model"
	"github.com/revueexchange/api/internal/repository"
)

// AnalyticsService handles analytics operations
type AnalyticsService struct {
	repo *repository.Repository
}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService(repo *repository.Repository) *AnalyticsService {
	return &AnalyticsService{repo: repo}
}

// GetOverviewStats returns overall platform statistics
func (s *AnalyticsService) GetOverviewStats(ctx context.Context) (*model.OverviewStats, error) {
	stats := &model.OverviewStats{}

	// Get total users
	totalUsers, err := s.repo.GetTotalUsers(ctx)
	if err != nil {
		return nil, err
	}
	stats.TotalUsers = totalUsers

	// Get total bounties
	totalBounties, err := s.repo.GetTotalBounties(ctx)
	if err != nil {
		return nil, err
	}
	stats.TotalBounties = totalBounties

	// Get total reviews
	totalReviews, err := s.repo.GetTotalReviews(ctx)
	if err != nil {
		return nil, err
	}
	stats.TotalReviews = totalReviews

	// Get total points awarded/spent
	pointsStats, err := s.repo.GetTotalPointsStats(ctx)
	if err != nil {
		return nil, err
	}
	stats.TotalPointsAwarded = pointsStats.Awarded
	stats.TotalPointsSpent = pointsStats.Spent

	return stats, nil
}

// GetBountyMetrics returns bounty performance metrics
func (s *AnalyticsService) GetBountyMetrics(ctx context.Context) ([]model.BountyMetrics, error) {
	return s.repo.GetBountyMetrics(ctx)
}

// GetReviewMetrics returns review metrics
func (s *AnalyticsService) GetReviewMetrics(ctx context.Context) ([]model.ReviewMetrics, error) {
	return s.repo.GetReviewMetrics(ctx)
}

// GetRevenueStats returns revenue analytics
func (s *AnalyticsService) GetRevenueStats(ctx context.Context) (*model.RevenueStats, error) {
	return s.repo.GetRevenueStats(ctx)
}

// GetUserActivity returns daily user activity
func (s *AnalyticsService) GetUserActivity(ctx context.Context, days int) ([]model.UserActivity, error) {
	return s.repo.GetUserActivity(ctx, days)
}

// GetUserAnalytics returns analytics for a specific user
func (s *AnalyticsService) GetUserAnalytics(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	analytics := make(map[string]interface{})

	// Get user's review count
	reviewCount, err := s.repo.GetReviewCount(ctx, userID)
	if err == nil {
		analytics["review_count"] = reviewCount
	}

	// Get user's bounties
	bountyCount, err := s.repo.GetUserBountyCount(ctx, userID)
	if err == nil {
		analytics["bounty_count"] = bountyCount
	}

	// Get user's reviews
	reviews, err := s.repo.GetReviewsByReviewer(ctx, userID)
	if err == nil {
		analytics["reviews"] = reviews
		var totalWords int
		var totalRating float64
		for _, r := range reviews {
			if r.WordCount != nil {
				totalWords += *r.WordCount
			}
			totalRating += float64(r.Rating)
		}
		if len(reviews) > 0 {
			analytics["avg_rating"] = totalRating / float64(len(reviews))
			analytics["total_words"] = totalWords
		}
	}

	return analytics, nil
}
