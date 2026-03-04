package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/revueexchange/api/internal/model"
	"github.com/revueexchange/api/internal/repository"
)

// BadgeService handles badge operations
type BadgeService struct {
	repo         *repository.Repository
	badgeRepo    *repository.BadgeRepository
}

// NewBadgeService creates a new badge service
func NewBadgeService(repo *repository.Repository, badgeRepo *repository.BadgeRepository) *BadgeService {
	return &BadgeService{
		repo:      repo,
		badgeRepo: badgeRepo,
	}
}

// CheckAndAwardBadges checks user stats and awards applicable badges
func (s *BadgeService) CheckAndAwardBadges(ctx context.Context, userID uuid.UUID) ([]model.Badge, error) {
	// Get user stats
	stats, err := s.getUserStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	awardedBadges := []model.Badge{}

	// Check each badge type
	for _, def := range model.BadgeDefinitions {
		// Check if user already has this badge type
		hasBadge, err := s.badgeRepo.HasBadge(ctx, userID.String(), def.Type)
		if err != nil {
			continue
		}
		if hasBadge {
			continue
		}

		// Check if user meets condition
		earned, tier := s.checkCondition(def, stats)
		if earned {
			badge := &model.Badge{
				UserID:      userID.String(),
				BadgeType:   def.Type,
				BadgeName:   def.Name,
				Description: def.Description,
				Tier:        tier,
				AwardedAt:   time.Now(),
			}

			if err := s.badgeRepo.AwardBadge(ctx, badge); err != nil {
				continue
			}

			awardedBadges = append(awardedBadges, *badge)
		}
	}

	return awardedBadges, nil
}

// getUserStats gets user statistics for badge checking
func (s *BadgeService) getUserStats(ctx context.Context, userID uuid.UUID) (*model.UserStats, error) {
	// Get review count
	reviewCount, err := s.repo.GetReviewCount(ctx, userID)
	if err != nil {
		reviewCount = 0
	}

	// Get helpful votes (placeholder - would need to add to repo)
	helpfulVotes := 0

	// Get member since
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return &model.UserStats{
			UserID:        userID.String(),
			ReviewCount:   reviewCount,
			HelpfulVotes:  helpfulVotes,
			MemberSince:    time.Now(),
		}, nil
	}

	return &model.UserStats{
		UserID:         userID.String(),
		ReviewCount:    reviewCount,
		HelpfulVotes:   helpfulVotes,
		IsTopReviewer:  false, // Would need leaderboard logic
		MemberSince:    user.CreatedAt,
	}, nil
}

// checkCondition checks if user meets badge condition
func (s *BadgeService) checkCondition(def model.BadgeDefinition, stats *model.UserStats) (bool, string) {
	switch def.Condition {
	case "first_review":
		if stats.ReviewCount >= 1 {
			return true, "bronze"
		}
	case "review_10":
		if stats.ReviewCount >= 10 {
			if stats.ReviewCount >= 50 {
				return true, "gold"
			}
			if stats.ReviewCount >= 25 {
				return true, "silver"
			}
			return true, "bronze"
		}
	case "review_50":
		if stats.ReviewCount >= 50 {
			if stats.ReviewCount >= 100 {
				return true, "platinum"
			}
			if stats.ReviewCount >= 75 {
				return true, "gold"
			}
			return true, "silver"
		}
	case "streak_7":
		if stats.StreakDays >= 7 {
			if stats.StreakDays >= 14 {
				return true, "silver"
			}
			return true, "bronze"
		}
	case "streak_30":
		if stats.StreakDays >= 30 {
			if stats.StreakDays >= 60 {
				return true, "platinum"
			}
			return true, "gold"
		}
	case "top_reviewer":
		if stats.IsTopReviewer {
			return true, "platinum"
		}
	case "helpful_100":
		if stats.HelpfulVotes >= 100 {
			if stats.HelpfulVotes >= 250 {
				return true, "gold"
			}
			return true, "silver"
		}
	case "early_adopter":
		// Beta cutoff - anyone who joined before a certain date
		if stats.MemberSince.Before(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)) {
			return true, "bronze"
		}
	}
	return false, ""
}

// GetUserBadges gets all badges for a user
func (s *BadgeService) GetUserBadges(ctx context.Context, userID uuid.UUID) ([]model.Badge, error) {
	return s.badgeRepo.GetUserBadges(ctx, userID.String())
}
