package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/revueexchange/api/internal/model"
	"github.com/revueexchange/api/internal/repository"
)

// FraudService handles fraud detection
type FraudService struct {
	repo *repository.Repository
}

// NewFraudService creates a new fraud service
func NewFraudService(repo *repository.Repository) *FraudService {
	return &FraudService{repo: repo}
}

// CheckReviewForFraud analyzes a review for suspicious patterns
func (s *FraudService) CheckReviewForFraud(ctx context.Context, userID uuid.UUID, bountyID uuid.UUID) *model.FraudAlert {
	// Check for rapid reviews (more than 3 reviews in 1 hour)
	recentReviews, err := s.repo.GetRecentReviewCount(ctx, userID, 1)
	if err == nil && recentReviews > 3 {
		return &model.FraudAlert{
			ID:          uuid.New().String(),
			UserID:      userID.String(),
			Type:        "rapid_reviews",
			Severity:    "medium",
			Description: "User submitted more than 3 reviews in 1 hour",
			CreatedAt:   time.Now().Format(time.RFC3339),
		}
	}

	// Check for very short reviews from new accounts
	reviews, err := s.repo.GetReviewsByReviewer(ctx, userID)
	if err == nil && len(reviews) < 5 {
		for _, r := range reviews {
			if r.WordCount != nil && *r.WordCount < 20 {
				return &model.FraudAlert{
					ID:          uuid.New().String(),
					UserID:      userID.String(),
					Type:        "suspicious_pattern",
					Severity:    "low",
					Description: "New user submitting very short reviews",
					CreatedAt:   time.Now().Format(time.RFC3339),
				}
			}
		}
	}

	return nil
}

// GetUserTrustScore calculates a user's trust score
func (s *FraudService) GetUserTrustScore(ctx context.Context, userID uuid.UUID) (*model.UserReputation, error) {
	reputation := &model.UserReputation{
		UserID:        userID.String(),
		TrustScore:    50.0, // base score
		ActivityScore: 50.0,
	}

	// Get review count
	reviewCount, err := s.repo.GetReviewCount(ctx, userID)
	if err == nil {
		// More reviews = higher trust
		if reviewCount > 10 {
			reputation.TrustScore += 20
		} else if reviewCount > 5 {
			reputation.TrustScore += 10
		}
	}

	// Get user account age
	user, err := s.repo.GetUserByID(ctx, userID)
	if err == nil {
		age := time.Since(user.CreatedAt).Hours()
		if age > 720 { // 30 days
			reputation.TrustScore += 20
		} else if age > 168 { // 7 days
			reputation.TrustScore += 10
		}
	}

	// Cap at 100
	if reputation.TrustScore > 100 {
		reputation.TrustScore = 100
	}

	return reputation, nil
}

// FlagReviewForVerification marks a review for manual verification
func (s *FraudService) FlagReviewForVerification(ctx context.Context, reviewID, amazonLink string) error {
	return s.repo.FlagReviewForVerification(ctx, reviewID, amazonLink)
}

// GetVerificationStatus gets the verification status of a review
func (s *FraudService) GetVerificationStatus(ctx context.Context, reviewID string) (*model.ReviewVerification, error) {
	return s.repo.GetReviewVerification(ctx, reviewID)
}
