package model

import (
	"time"

	"github.com/google/uuid"
)

// DB interface for database
type DB interface {
	Query(query string, args ...interface{}) (Rows, error)
	QueryRow(query string, args ...interface{}) Row
	Exec(query string, args ...interface{}) (Result, error)
}

type Rows interface {
	Close()
	Next() bool
	Scan(dest ...interface{}) error
}

type Row interface {
	Scan(dest ...interface{}) error
}

type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}

// User represents a user in the system
type User struct {
	ID              uuid.UUID  `json:"id"`
	Email           string     `json:"email"`
	PasswordHash    string     `json:"-"`
	Username        string     `json:"username"`
	DisplayName     string     `json:"display_name"`
	AvatarURL       *string    `json:"avatar_url"`
	Bio             *string    `json:"bio"`
	Points          int        `json:"points"`
	ReputationScore float64    `json:"reputation_score"`
	StripeCustomerID *string   `json:"stripe_customer_id,omitempty"`
	SubscriptionTier string    `json:"subscription_tier"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// Product represents a product (book, course, podcast, newsletter)
type Product struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Type        string    `json:"type"` // book, course, podcast, newsletter
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	URL         *string   `json:"url"`
	CoverImageURL *string `json:"cover_image_url"`
	Genre       *string   `json:"genre"`
	WordCount   *int      `json:"word_count"`
	CreatedAt   time.Time `json:"created_at"`
}

// Bounty represents a review bounty
type Bounty struct {
	ID           uuid.UUID  `json:"id"`
	UserID       uuid.UUID  `json:"user_id"`
	ProductID    uuid.UUID  `json:"product_id"`
	BountyPoints int        `json:"bounty_points"`
	BountyCash   *float64  `json:"bounty_cash"`
	Status       string     `json:"status"` // open, claimed, under_review, completed, cancelled
	Requirements *string    `json:"requirements"`
	ClaimedBy    *uuid.UUID `json:"claimed_by"`
	ClaimedAt    *time.Time `json:"claimed_at"`
	CompletedAt   *time.Time `json:"completed_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// Review represents a review
type Review struct {
	ID               uuid.UUID  `json:"id"`
	BountyID         uuid.UUID  `json:"bounty_id"`
	ReviewerID       uuid.UUID  `json:"reviewer_id"`
	Rating           int        `json:"rating"`
	Title            *string    `json:"title"`
	Content          *string    `json:"content"`
	WordCount        *int       `json:"word_count"`
	VerifiedPurchase bool       `json:"verified_purchase"`
	Status           string     `json:"status"` // draft, submitted, published
	AmazonReviewURL  *string    `json:"amazon_review_url"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// PointTransaction represents a point transaction
type PointTransaction struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"user_id"`
	Amount        int       `json:"amount"`
	Type          string    `json:"type"` // earned, spent, bonus, penalty, refund
	ReferenceType *string   `json:"reference_type"`
	ReferenceID   *uuid.UUID `json:"reference_id"`
	Description   *string   `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
}

// Follow represents a follow relationship
type Follow struct {
	FollowerID  uuid.UUID `json:"follower_id"`
	FollowingID uuid.UUID `json:"following_id"`
	CreatedAt   time.Time `json:"created_at"`
}

// Payment represents a payment
type Payment struct {
	ID                  uuid.UUID  `json:"id"`
	UserID              uuid.UUID  `json:"user_id"`
	StripeSessionID    *string    `json:"stripe_session_id"`
	StripePaymentIntent *string    `json:"stripe_payment_intent"`
	AmountCents        int        `json:"amount_cents"`
	Currency           string     `json:"currency"`
	Type               *string    `json:"type"`
	Status             string     `json:"status"`
	CreatedAt          time.Time  `json:"created_at"`
}

// Comment represents a comment on a review
type Comment struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	ReviewID  uuid.UUID  `json:"review_id"`
	ParentID  *uuid.UUID `json:"parent_id"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"created_at"`
}

// Activity represents an activity feed item
type Activity struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Type      string    `json:"type"` // review_submitted, bounty_created, follow, etc.
	Reference *string   `json:"reference"`
	CreatedAt time.Time `json:"created_at"`
}
