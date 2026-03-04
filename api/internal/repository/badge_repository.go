package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/revueexchange/api/internal/model"
)

const BadgesTableName = "revueexchange-badges"

// BadgeRepository handles badge operations in DynamoDB
type BadgeRepository struct {
	client *dynamodb.Client
}

// NewBadgeRepository creates a new badge repository
func NewBadgeRepository(client *dynamodb.Client) *BadgeRepository {
	return &BadgeRepository{client: client}
}

// CreateBadgesTable creates the badges table if it doesn't exist
func (r *BadgeRepository) CreateBadgesTable(ctx context.Context) error {
	_, err := r.client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(BadgesTableName),
	})
	if err == nil {
		// Table already exists
		return nil
	}

	// Create table
	_, err = r.client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(BadgesTableName),
		KeySchema: []types.KeySchemaElement{
			{AttributeName: aws.String("user_id"), KeyType: types.KeyTypeHash},
			{AttributeName: aws.String("id"), KeyType: types.KeyTypeRange},
		},
		AttributeDefinitions: []types.AttributeDefinition{
			{AttributeName: aws.String("user_id"), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String("badge_type"), AttributeType: types.ScalarAttributeTypeS},
		},
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String("badge_type-index"),
				KeySchema: []types.KeySchemaElement{
					{AttributeName: aws.String("badge_type"), KeyType: types.KeyTypeHash},
					{AttributeName: aws.String("user_id"), KeyType: types.KeyTypeRange},
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
	return err
}

// BadgeItem represents a badge item in DynamoDB
type BadgeItem struct {
	ID          string `dynamodbav:"id"`
	UserID      string `dynamodbav:"user_id"`
	BadgeType   string `dynamodbav:"badge_type"`
	BadgeName   string `dynamodbav:"badge_name"`
	Description string `dynamodbav:"description"`
	Tier        string `dynamodbav:"tier"`
	IconURL     string `dynamodbav:"icon_url,omitempty"`
	AwardedAt   string `dynamodbav:"awarded_at"`
}

// AwardBadge awards a badge to a user
func (r *BadgeRepository) AwardBadge(ctx context.Context, badge *model.Badge) error {
	item := BadgeItem{
		ID:          uuid.New().String(),
		UserID:      badge.UserID,
		BadgeType:   badge.BadgeType,
		BadgeName:   badge.BadgeName,
		Description: badge.Description,
		Tier:        badge.Tier,
		IconURL:     badge.IconURL,
		AwardedAt:   badge.AwardedAt.Format(time.RFC3339),
	}

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal badge: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(BadgesTableName),
		Item:      av,
	})
	return err
}

// GetUserBadges gets all badges for a user
func (r *BadgeRepository) GetUserBadges(ctx context.Context, userID string) ([]model.Badge, error) {
	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(BadgesTableName),
		KeyConditionExpression: aws.String("user_id = :user_id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":user_id": &types.AttributeValueMemberS{Value: userID},
		},
	})
	if err != nil {
		return nil, err
	}

	var items []BadgeItem
	err = attributevalue.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		return nil, err
	}

	badges := make([]model.Badge, len(items))
	for i, item := range items {
		badges[i] = model.Badge{
			ID:          item.ID,
			UserID:      item.UserID,
			BadgeType:   item.BadgeType,
			BadgeName:   item.BadgeName,
			Description: item.Description,
			Tier:        item.Tier,
			IconURL:     item.IconURL,
		}
		if t, err := time.Parse(time.RFC3339, item.AwardedAt); err == nil {
			badges[i].AwardedAt = t
		}
	}

	return badges, nil
}

// HasBadge checks if user already has a badge of a specific type
func (r *BadgeRepository) HasBadge(ctx context.Context, userID, badgeType string) (bool, error) {
	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(BadgesTableName),
		KeyConditionExpression: aws.String("user_id = :user_id"),
		FilterExpression:       aws.String("badge_type = :badge_type"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":user_id":    &types.AttributeValueMemberS{Value: userID},
			":badge_type": &types.AttributeValueMemberS{Value: badgeType},
		},
		Limit: aws.Int32(1),
	})
	if err != nil {
		return false, err
	}

	return len(result.Items) > 0, nil
}
