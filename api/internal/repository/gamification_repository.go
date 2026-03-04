package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/revueexchange/api/internal/model"
)

const (
	LeaderboardTableName = "revueexchange-leaderboard"
	StreaksTableName     = "revueexchange-streaks"
)

// GamificationRepository handles leaderboard and streak operations in DynamoDB
type GamificationRepository struct {
	client *dynamodb.Client
}

// NewGamificationRepository creates a new gamification repository
func NewGamificationRepository(client *dynamodb.Client) *GamificationRepository {
	return &GamificationRepository{client: client}
}

// CreateTables creates leaderboard and streaks tables if they don't exist
func (r *GamificationRepository) CreateTables(ctx context.Context) error {
	// Create leaderboard table
	_, err := r.client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(LeaderboardTableName),
	})
	if err != nil {
		_, err = r.client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(LeaderboardTableName),
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("rank"), KeyType: types.KeyTypeRange},
			},
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("rank"), AttributeType: types.ScalarAttributeTypeN},
				{AttributeName: aws.String("user_id"), AttributeType: types.ScalarAttributeTypeS},
			},
			GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
				{
					IndexName: aws.String("user_id-index"),
					KeySchema: []types.KeySchemaElement{
						{AttributeName: aws.String("user_id"), KeyType: types.KeyTypeHash},
					},
					Projection: &types.Projection{
						ProjectionType: types.ProjectionTypeAll,
					},
					ProvisionedThroughput: &types.ProvisionedThroughput{
						ReadCapacityUnits:  aws.Int64(5),
						WriteCapacityUnits: aws.Int64(5),
					},
				},
			},
			ProvisionedThroughput: &types.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(5),
				WriteCapacityUnits: aws.Int64(5),
			},
		})
		if err != nil {
			return fmt.Errorf("failed to create leaderboard table: %w", err)
		}
	}

	// Create streaks table
	_, err = r.client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(StreaksTableName),
	})
	if err != nil {
		_, err = r.client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(StreaksTableName),
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("user_id"), KeyType: types.KeyTypeHash},
			},
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("user_id"), AttributeType: types.ScalarAttributeTypeS},
			},
			ProvisionedThroughput: &types.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(5),
				WriteCapacityUnits: aws.Int64(5),
			},
		})
		if err != nil {
			return fmt.Errorf("failed to create streaks table: %w", err)
		}
	}

	return nil
}

// LeaderboardItem represents a leaderboard entry in DynamoDB
type LeaderboardItem struct {
	UserID      string `dynamodbav:"user_id"`
	Username    string `dynamodbav:"username"`
	DisplayName string `dynamodbav:"display_name"`
	Points      int    `dynamodbav:"points"`
	ReviewCount int    `dynamodbav:"review_count"`
	Rank        int    `dynamodbav:"rank"`
	LastUpdated string `dynamodbav:"last_updated"`
}

// UpdateLeaderboard updates a user's position on the leaderboard
func (r *GamificationRepository) UpdateLeaderboard(ctx context.Context, userID, username, displayName string, points, reviewCount int) error {
	// First, get the current rank to update it
	// For simplicity, we'll rebuild the leaderboard periodically
	// In production, you'd use a more sophisticated approach

	item := LeaderboardItem{
		UserID:      userID,
		Username:    username,
		DisplayName: displayName,
		Points:      points,
		ReviewCount: reviewCount,
		Rank:        0, // Will be calculated
		LastUpdated: time.Now().Format(time.RFC3339),
	}

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal leaderboard item: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(LeaderboardTableName),
		Item:      av,
	})
	return err
}

// GetLeaderboard gets the top users by points
func (r *GamificationRepository) GetLeaderboard(ctx context.Context, limit int) ([]model.LeaderboardEntry, error) {
	result, err := r.client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(LeaderboardTableName),
		Limit:    aws.Int32(int32(limit)),
	})
	if err != nil {
		return nil, err
	}

	var items []LeaderboardItem
	err = attributevalue.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		return nil, err
	}

	// Sort by points and assign ranks
	entries := make([]model.LeaderboardEntry, len(items))
	for i, item := range items {
		entries[i] = model.LeaderboardEntry{
			UserID:         item.UserID,
			Username:       item.Username,
			DisplayName:    item.DisplayName,
			Points:         item.Points,
			ReviewCount:    item.ReviewCount,
			Rank:           i + 1,
		}
		if t, err := time.Parse(time.RFC3339, item.LastUpdated); err == nil {
			entries[i].LastUpdated = t
		}
	}

	return entries, nil
}

// StreakItem represents a streak entry in DynamoDB
type StreakItem struct {
	UserID          string `dynamodbav:"user_id"`
	CurrentStreak   int    `dynamodbav:"current_streak"`
	LongestStreak   int    `dynamodbav:"longest_streak"`
	LastActivityAt  string `dynamodbav:"last_activity_at"`
	StreakStartedAt string `dynamodbav:"streak_started_at"`
}

// GetStreak gets a user's streak
func (r *GamificationRepository) GetStreak(ctx context.Context, userID string) (*model.Streak, error) {
	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(StreaksTableName),
		Key: map[string]types.AttributeValue{
			"user_id": &types.AttributeValueMemberS{Value: userID},
		},
	})
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return &model.Streak{
			UserID:         userID,
			CurrentStreak:  0,
			LongestStreak:  0,
		}, nil
	}

	var item StreakItem
	err = attributevalue.UnmarshalMap(result.Item, &item)
	if err != nil {
		return nil, err
	}

	streak := model.Streak{
		UserID:         item.UserID,
		CurrentStreak:  item.CurrentStreak,
		LongestStreak:  item.LongestStreak,
	}

	if t, err := time.Parse(time.RFC3339, item.LastActivityAt); err == nil {
		streak.LastActivityAt = t
	}
	if t, err := time.Parse(time.RFC3339, item.StreakStartedAt); err == nil {
		streak.StreakStartedAt = t
	}

	return &streak, nil
}

// UpdateStreak updates a user's streak based on activity
func (r *GamificationRepository) UpdateStreak(ctx context.Context, userID string) (*model.Streak, error) {
	streak, err := r.GetStreak(ctx, userID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	if streak.LastActivityAt.IsZero() {
		// First activity
		streak.CurrentStreak = 1
		streak.LongestStreak = 1
		streak.StreakStartedAt = today
	} else {
		lastActivity := time.Date(streak.LastActivityAt.Year(), streak.LastActivityAt.Month(), streak.LastActivityAt.Day(), 0, 0, 0, 0, streak.LastActivityAt.Location())
		daysSince := int(today.Sub(lastActivity).Hours() / 24)

		if daysSince == 0 {
			// Already updated today, no change
		} else if daysSince == 1 {
			// Consecutive day
			streak.CurrentStreak++
			if streak.CurrentStreak > streak.LongestStreak {
				streak.LongestStreak = streak.CurrentStreak
			}
		} else {
			// Streak broken
			streak.CurrentStreak = 1
			streak.StreakStartedAt = today
		}
	}

	streak.LastActivityAt = now
	streak.UserID = userID

	// Save to DynamoDB
	item := StreakItem{
		UserID:          streak.UserID,
		CurrentStreak:  streak.CurrentStreak,
		LongestStreak:  streak.LongestStreak,
		LastActivityAt:  streak.LastActivityAt.Format(time.RFC3339),
		StreakStartedAt: streak.StreakStartedAt.Format(time.RFC3339),
	}

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return nil, err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(StreaksTableName),
		Item:      av,
	})
	if err != nil {
		return nil, err
	}

	return streak, nil
}
