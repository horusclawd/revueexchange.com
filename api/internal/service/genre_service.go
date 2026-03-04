package service

import (
	"context"

	"github.com/revueexchange/api/internal/model"
	"github.com/revueexchange/api/internal/repository"
)

// GenreService handles genre matching
type GenreService struct {
	repo *repository.Repository
}

// NewGenreService creates a new genre service
func NewGenreService(repo *repository.Repository) *GenreService {
	return &GenreService{repo: repo}
}

// GetMatchingBounties returns bounties matching user's genre preferences
func (s *GenreService) GetMatchingBounties(ctx context.Context, userID string, limit int) ([]model.Bounty, error) {
	// Get user's genre preferences
	genres, err := s.repo.GetUserGenres(ctx, userID)
	if err != nil || len(genres) == 0 {
		// Return popular bounties if no preferences
		return s.repo.GetPopularBounties(ctx, limit)
	}

	// Get bounties matching genres
	return s.repo.GetBountiesByGenres(ctx, genres, limit)
}

// GetRecommendedReviewers returns reviewers matching bounty genre
func (s *GenreService) GetRecommendedReviewers(ctx context.Context, genre string, limit int) ([]string, error) {
	return s.repo.GetReviewersByGenre(ctx, genre, limit)
}

// UpdateUserGenres updates user's genre preferences
func (s *GenreService) UpdateUserGenres(ctx context.Context, userID string, genres, expertise, interests []string) error {
	return s.repo.UpdateUserGenres(ctx, userID, genres, expertise, interests)
}
