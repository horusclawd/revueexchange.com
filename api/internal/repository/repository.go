package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/revueexchange/api/internal/model"
)

// Repository holds the database connection
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new repository
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// User methods
func (r *Repository) CreateUser(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, username, display_name, points, reputation_score, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Exec(ctx, query,
		user.ID, user.Email, user.PasswordHash, user.Username, user.DisplayName,
		user.Points, user.ReputationScore, user.CreatedAt, user.UpdatedAt,
	)
	return err
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT id, email, password_hash, username, display_name, avatar_url, bio, points, reputation_score, subscription_tier, created_at, updated_at FROM users WHERE email = $1`
	row := r.db.QueryRow(ctx, query, email)

	var user model.User
	err := row.Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Username, &user.DisplayName,
		&user.AvatarURL, &user.Bio, &user.Points, &user.ReputationScore,
		&user.SubscriptionTier, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	query := `SELECT id, email, password_hash, username, display_name, avatar_url, bio, points, reputation_score, subscription_tier, created_at, updated_at FROM users WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var user model.User
	err := row.Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Username, &user.DisplayName,
		&user.AvatarURL, &user.Bio, &user.Points, &user.ReputationScore,
		&user.SubscriptionTier, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) UpdateUser(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users SET display_name = $1, avatar_url = $2, bio = $3, updated_at = $4
		WHERE id = $5
	`
	_, err := r.db.Exec(ctx, query, user.DisplayName, user.AvatarURL, user.Bio, user.UpdatedAt, user.ID)
	return err
}

// Health check
func (r *Repository) HealthCheck(ctx context.Context) error {
	return r.db.Ping(ctx)
}

// Product methods
func (r *Repository) CreateProduct(ctx context.Context, product *model.Product) error {
	query := `
		INSERT INTO products (id, user_id, type, title, description, url, cover_image_url, genre, word_count, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.Exec(ctx, query,
		product.ID, product.UserID, product.Type, product.Title, product.Description,
		product.URL, product.CoverImageURL, product.Genre, product.WordCount, product.CreatedAt,
	)
	return err
}

func (r *Repository) GetProductByID(ctx context.Context, id uuid.UUID) (*model.Product, error) {
	query := `SELECT id, user_id, type, title, description, url, cover_image_url, genre, word_count, created_at FROM products WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var product model.Product
	err := row.Scan(&product.ID, &product.UserID, &product.Type, &product.Title, &product.Description, &product.URL, &product.CoverImageURL, &product.Genre, &product.WordCount, &product.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *Repository) GetProductsByUserID(ctx context.Context, userID uuid.UUID) ([]model.Product, error) {
	query := `SELECT id, user_id, type, title, description, url, cover_image_url, genre, word_count, created_at FROM products WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.UserID, &p.Type, &p.Title, &p.Description, &p.URL, &p.CoverImageURL, &p.Genre, &p.WordCount, &p.CreatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *Repository) UpdateProduct(ctx context.Context, product *model.Product) error {
	query := `
		UPDATE products SET title = $1, description = $2, url = $3, cover_image_url = $4, genre = $5, word_count = $6
		WHERE id = $7 AND user_id = $8
	`
	_, err := r.db.Exec(ctx, query, product.Title, product.Description, product.URL, product.CoverImageURL, product.Genre, product.WordCount, product.ID, product.UserID)
	return err
}

func (r *Repository) DeleteProduct(ctx context.Context, id, userID uuid.UUID) error {
	query := `DELETE FROM products WHERE id = $1 AND user_id = $2`
	_, err := r.db.Exec(ctx, query, id, userID)
	return err
}

// Review methods
func (r *Repository) CreateReview(ctx context.Context, review *model.Review) error {
	query := `
		INSERT INTO reviews (id, bounty_id, reviewer_id, rating, title, content, word_count, verified_purchase, status, amazon_review_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := r.db.Exec(ctx, query,
		review.ID, review.BountyID, review.ReviewerID, review.Rating, review.Title,
		review.Content, review.WordCount, review.VerifiedPurchase, review.Status,
		review.AmazonReviewURL, review.CreatedAt, review.UpdatedAt,
	)
	return err
}

func (r *Repository) GetReviewByID(ctx context.Context, id uuid.UUID) (*model.Review, error) {
	query := `SELECT id, bounty_id, reviewer_id, rating, title, content, word_count, verified_purchase, status, amazon_review_url, created_at, updated_at FROM reviews WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var review model.Review
	err := row.Scan(&review.ID, &review.BountyID, &review.ReviewerID, &review.Rating,
		&review.Title, &review.Content, &review.WordCount, &review.VerifiedPurchase,
		&review.Status, &review.AmazonReviewURL, &review.CreatedAt, &review.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *Repository) GetReviewByBountyID(ctx context.Context, bountyID uuid.UUID) (*model.Review, error) {
	query := `SELECT id, bounty_id, reviewer_id, rating, title, content, word_count, verified_purchase, status, amazon_review_url, created_at, updated_at FROM reviews WHERE bounty_id = $1`
	row := r.db.QueryRow(ctx, query, bountyID)

	var review model.Review
	err := row.Scan(&review.ID, &review.BountyID, &review.ReviewerID, &review.Rating,
		&review.Title, &review.Content, &review.WordCount, &review.VerifiedPurchase,
		&review.Status, &review.AmazonReviewURL, &review.CreatedAt, &review.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *Repository) GetReviewsByReviewer(ctx context.Context, reviewerID uuid.UUID) ([]model.Review, error) {
	query := `SELECT id, bounty_id, reviewer_id, rating, title, content, word_count, verified_purchase, status, amazon_review_url, created_at, updated_at FROM reviews WHERE reviewer_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query, reviewerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []model.Review
	for rows.Next() {
		var review model.Review
		if err := rows.Scan(&review.ID, &review.BountyID, &review.ReviewerID, &review.Rating,
			&review.Title, &review.Content, &review.WordCount, &review.VerifiedPurchase,
			&review.Status, &review.AmazonReviewURL, &review.CreatedAt, &review.UpdatedAt); err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	return reviews, nil
}

func (r *Repository) GetReviewCount(ctx context.Context, reviewerID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM reviews WHERE reviewer_id = $1 AND status = 'published'`
	var count int
	err := r.db.QueryRow(ctx, query, reviewerID).Scan(&count)
	return count, err
}

func (r *Repository) UpdateReview(ctx context.Context, review *model.Review) error {
	query := `
		UPDATE reviews SET rating = $1, title = $2, content = $3, word_count = $4, status = $5, amazon_review_url = $6, updated_at = $7
		WHERE id = $8
	`
	_, err := r.db.Exec(ctx, query, review.Rating, review.Title, review.Content, review.WordCount, review.Status, review.AmazonReviewURL, review.UpdatedAt, review.ID)
	return err
}

// Placeholder methods for other entities
func (r *Repository) CreateBounty(ctx context.Context, bounty *model.Bounty) error {
	query := `
		INSERT INTO bounties (id, user_id, product_id, bounty_points, bounty_cash, status, requirements, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Exec(ctx, query,
		bounty.ID, bounty.UserID, bounty.ProductID, bounty.BountyPoints, bounty.BountyCash,
		bounty.Status, bounty.Requirements, bounty.CreatedAt, bounty.UpdatedAt,
	)
	return err
}

// BountyFilters holds filter parameters for bounties
type BountyFilters struct {
	Status     string
	Genre      string
	Type       string
	MinPoints  int
	MaxPoints  int
	UserID     uuid.UUID
	ClaimedBy  uuid.UUID
}

func (r *Repository) GetBounties(ctx context.Context, filters BountyFilters, limit, offset int) ([]model.Bounty, error) {
	// Build dynamic query
	query := `
		SELECT b.id, b.user_id, b.product_id, b.bounty_points, b.bounty_cash, b.status, b.requirements,
		       b.claimed_by, b.claimed_at, b.completed_at, b.created_at, b.updated_at
		FROM bounties b
		JOIN products p ON b.product_id = p.id
		WHERE 1=1
	`
	args := []interface{}{}
	argNum := 1

	if filters.Status != "" {
		query += fmt.Sprintf(" AND b.status = $%d", argNum)
		args = append(args, filters.Status)
		argNum++
	}

	if filters.Genre != "" {
		query += fmt.Sprintf(" AND p.genre = $%d", argNum)
		args = append(args, filters.Genre)
		argNum++
	}

	if filters.Type != "" {
		query += fmt.Sprintf(" AND p.type = $%d", argNum)
		args = append(args, filters.Type)
		argNum++
	}

	if filters.MinPoints > 0 {
		query += fmt.Sprintf(" AND b.bounty_points >= $%d", argNum)
		args = append(args, filters.MinPoints)
		argNum++
	}

	if filters.MaxPoints > 0 {
		query += fmt.Sprintf(" AND b.bounty_points <= $%d", argNum)
		args = append(args, filters.MaxPoints)
		argNum++
	}

	if filters.UserID != uuid.Nil {
		query += fmt.Sprintf(" AND b.user_id = $%d", argNum)
		args = append(args, filters.UserID)
		argNum++
	}

	if filters.ClaimedBy != uuid.Nil {
		query += fmt.Sprintf(" AND b.claimed_by = $%d", argNum)
		args = append(args, filters.ClaimedBy)
		argNum++
	}

	query += fmt.Sprintf(" ORDER BY b.created_at DESC LIMIT $%d OFFSET $%d", argNum, argNum+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bounties []model.Bounty
	for rows.Next() {
		var b model.Bounty
		if err := rows.Scan(&b.ID, &b.UserID, &b.ProductID, &b.BountyPoints, &b.BountyCash, &b.Status, &b.Requirements, &b.ClaimedBy, &b.ClaimedAt, &b.CompletedAt, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		bounties = append(bounties, b)
	}
	return bounties, nil
}

func (r *Repository) GetBountyCount(ctx context.Context, filters BountyFilters) (int, error) {
	query := `SELECT COUNT(*) FROM bounties b JOIN products p ON b.product_id = p.id WHERE 1=1`
	var count int

	if filters.Status != "" {
		query += " AND b.status = '" + filters.Status + "'"
	}

	err := r.db.QueryRow(ctx, query).Scan(&count)
	return count, err
}

func (r *Repository) GetBountyByID(ctx context.Context, id uuid.UUID) (*model.Bounty, error) {
	query := `SELECT id, user_id, product_id, bounty_points, bounty_cash, status, requirements, claimed_by, claimed_at, completed_at, created_at, updated_at FROM bounties WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var b model.Bounty
	err := row.Scan(&b.ID, &b.UserID, &b.ProductID, &b.BountyPoints, &b.BountyCash, &b.Status, &b.Requirements, &b.ClaimedBy, &b.ClaimedAt, &b.CompletedAt, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *Repository) UpdateBounty(ctx context.Context, bounty *model.Bounty) error {
	query := `
		UPDATE bounties SET status = $1, claimed_by = $2, claimed_at = $3, completed_at = $4, updated_at = $5
		WHERE id = $6
	`
	_, err := r.db.Exec(ctx, query, bounty.Status, bounty.ClaimedBy, bounty.ClaimedAt, bounty.CompletedAt, bounty.UpdatedAt, bounty.ID)
	return err
}

// CreatePointTransaction creates a point transaction
func (r *Repository) CreatePointTransaction(ctx context.Context, tx *model.PointTransaction) error {
	query := `
		INSERT INTO point_transactions (id, user_id, amount, type, reference_type, reference_id, description, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(ctx, query,
		tx.ID, tx.UserID, tx.Amount, tx.Type, tx.ReferenceType, tx.ReferenceID, tx.Description, tx.CreatedAt,
	)
	return err
}

// UpdateUserPoints updates user's points
func (r *Repository) UpdateUserPoints(ctx context.Context, userID uuid.UUID, points int) error {
	query := `UPDATE users SET points = points + $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(ctx, query, points, userID)
	return err
}

// GetUserPoints gets user's current points
func (r *Repository) GetUserPoints(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `SELECT points FROM users WHERE id = $1`
	var points int
	err := r.db.QueryRow(ctx, query, userID).Scan(&points)
	return points, err
}

func (r *Repository) GetPointTransactions(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.PointTransaction, error) {
	query := `SELECT id, user_id, amount, type, reference_type, reference_id, description, created_at FROM point_transactions WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []model.PointTransaction
	for rows.Next() {
		var tx model.PointTransaction
		if err := rows.Scan(&tx.ID, &tx.UserID, &tx.Amount, &tx.Type, &tx.ReferenceType, &tx.ReferenceID, &tx.Description, &tx.CreatedAt); err != nil {
			return nil, err
		}
		txs = append(txs, tx)
	}
	return txs, nil
}

// Payment methods
func (r *Repository) CreatePayment(ctx context.Context, payment *model.Payment) error {
	query := `
		INSERT INTO payments (id, user_id, stripe_session_id, stripe_payment_intent, amount_cents, currency, type, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Exec(ctx, query,
		payment.ID, payment.UserID, payment.StripeSessionID, payment.StripePaymentIntent,
		payment.AmountCents, payment.Currency, payment.Type, payment.Status, payment.CreatedAt,
	)
	return err
}

func (r *Repository) GetPaymentBySessionID(ctx context.Context, sessionID string) (*model.Payment, error) {
	query := `SELECT id, user_id, stripe_session_id, stripe_payment_intent, amount_cents, currency, type, status, created_at FROM payments WHERE stripe_session_id = $1`
	row := r.db.QueryRow(ctx, query, sessionID)

	var payment model.Payment
	err := row.Scan(&payment.ID, &payment.UserID, &payment.StripeSessionID, &payment.StripePaymentIntent,
		&payment.AmountCents, &payment.Currency, &payment.Type, &payment.Status, &payment.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *Repository) UpdatePaymentStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `UPDATE payments SET status = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, status, id)
	return err
}

func (r *Repository) GetPaymentsByUserID(ctx context.Context, userID uuid.UUID) ([]model.Payment, error) {
	query := `SELECT id, user_id, stripe_session_id, stripe_payment_intent, amount_cents, currency, type, status, created_at FROM payments WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []model.Payment
	for rows.Next() {
		var p model.Payment
		if err := rows.Scan(&p.ID, &p.UserID, &p.StripeSessionID, &p.StripePaymentIntent,
			&p.AmountCents, &p.Currency, &p.Type, &p.Status, &p.CreatedAt); err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}
	return payments, nil
}

// Follow methods
func (r *Repository) CreateFollow(ctx context.Context, follow *model.Follow) error {
	query := `INSERT INTO follows (follower_id, following_id, created_at) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, query, follow.FollowerID, follow.FollowingID, follow.CreatedAt)
	return err
}

func (r *Repository) DeleteFollow(ctx context.Context, followerID, followingID uuid.UUID) error {
	query := `DELETE FROM follows WHERE follower_id = $1 AND following_id = $2`
	_, err := r.db.Exec(ctx, query, followerID, followingID)
	return err
}

func (r *Repository) IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM follows WHERE follower_id = $1 AND following_id = $2)`
	var exists bool
	err := r.db.QueryRow(ctx, query, followerID, followingID).Scan(&exists)
	return exists, err
}

func (r *Repository) GetFollowers(ctx context.Context, userID uuid.UUID) ([]model.User, error) {
	query := `
		SELECT u.id, u.email, u.password_hash, u.username, u.display_name, u.avatar_url, u.bio, u.points, u.reputation_score, u.subscription_tier, u.created_at, u.updated_at
		FROM users u
		JOIN follows f ON u.id = f.follower_id
		WHERE f.following_id = $1
		ORDER BY f.created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Username, &u.DisplayName,
			&u.AvatarURL, &u.Bio, &u.Points, &u.ReputationScore, &u.SubscriptionTier,
			&u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *Repository) GetFollowing(ctx context.Context, userID uuid.UUID) ([]model.User, error) {
	query := `
		SELECT u.id, u.email, u.password_hash, u.username, u.display_name, u.avatar_url, u.bio, u.points, u.reputation_score, u.subscription_tier, u.created_at, u.updated_at
		FROM users u
		JOIN follows f ON u.id = f.following_id
		WHERE f.follower_id = $1
		ORDER BY f.created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Username, &u.DisplayName,
			&u.AvatarURL, &u.Bio, &u.Points, &u.ReputationScore, &u.SubscriptionTier,
			&u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// Comment methods
func (r *Repository) CreateComment(ctx context.Context, comment *model.Comment) error {
	query := `INSERT INTO comments (id, user_id, review_id, parent_id, content, created_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(ctx, query, comment.ID, comment.UserID, comment.ReviewID, comment.ParentID, comment.Content, comment.CreatedAt)
	return err
}

func (r *Repository) GetCommentsByReviewID(ctx context.Context, reviewID uuid.UUID) ([]model.Comment, error) {
	query := `SELECT id, user_id, review_id, parent_id, content, created_at FROM comments WHERE review_id = $1 ORDER BY created_at ASC`
	rows, err := r.db.Query(ctx, query, reviewID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []model.Comment
	for rows.Next() {
		var c model.Comment
		if err := rows.Scan(&c.ID, &c.UserID, &c.ReviewID, &c.ParentID, &c.Content, &c.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

func (r *Repository) DeleteComment(ctx context.Context, commentID, userID uuid.UUID) error {
	query := `DELETE FROM comments WHERE id = $1 AND user_id = $2`
	_, err := r.db.Exec(ctx, query, commentID, userID)
	return err
}

// Activity methods
func (r *Repository) CreateActivity(ctx context.Context, activity *model.Activity) error {
	query := `INSERT INTO activities (id, user_id, type, reference, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(ctx, query, activity.ID, activity.UserID, activity.Type, activity.Reference, activity.CreatedAt)
	return err
}

func (r *Repository) GetActivityFeed(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Activity, error) {
	query := `
		SELECT a.id, a.user_id, a.type, a.reference, a.created_at
		FROM activities a
		JOIN follows f ON a.user_id = f.following_id
		WHERE f.follower_id = $1
		ORDER BY a.created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []model.Activity
	for rows.Next() {
		var a model.Activity
		if err := rows.Scan(&a.ID, &a.UserID, &a.Type, &a.Reference, &a.CreatedAt); err != nil {
			return nil, err
		}
		activities = append(activities, a)
	}
	return activities, nil
}

// Analytics queries

type PointsStats struct {
	Awarded int
	Spent   int
}

func (r *Repository) GetTotalUsers(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
	return count, err
}

func (r *Repository) GetTotalBounties(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM bounties").Scan(&count)
	return count, err
}

func (r *Repository) GetTotalReviews(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM reviews").Scan(&count)
	return count, err
}

func (r *Repository) GetTotalPointsStats(ctx context.Context) (*PointsStats, error) {
	stats := &PointsStats{}
	err := r.db.QueryRow(ctx, "SELECT COALESCE(SUM(CASE WHEN amount > 0 THEN amount ELSE 0 END), 0), COALESCE(SUM(CASE WHEN amount < 0 THEN ABS(amount) ELSE 0 END), 0) FROM point_transactions").Scan(&stats.Awarded, &stats.Spent)
	return stats, err
}

func (r *Repository) GetBountyMetrics(ctx context.Context) ([]model.BountyMetrics, error) {
	query := `
		SELECT status, COUNT(*), COALESCE(SUM(bounty_points), 0)
		FROM bounties
		GROUP BY status
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []model.BountyMetrics
	for rows.Next() {
		var m model.BountyMetrics
		if err := rows.Scan(&m.Status, &m.Count, &m.TotalBounties); err != nil {
			return nil, err
		}
		metrics = append(metrics, m)
	}
	return metrics, nil
}

func (r *Repository) GetReviewMetrics(ctx context.Context) ([]model.ReviewMetrics, error) {
	query := `
		SELECT status, COUNT(*), COALESCE(AVG(rating), 0), COALESCE(SUM(COALESCE(word_count, 0)), 0)
		FROM reviews
		GROUP BY status
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []model.ReviewMetrics
	for rows.Next() {
		var m model.ReviewMetrics
		if err := rows.Scan(&m.Status, &m.Count, &m.AvgRating, &m.TotalWords); err != nil {
			return nil, err
		}
		metrics = append(metrics, m)
	}
	return metrics, nil
}

func (r *Repository) GetRevenueStats(ctx context.Context) (*model.RevenueStats, error) {
	stats := &model.RevenueStats{}
	err := r.db.QueryRow(ctx, `
		SELECT COALESCE(SUM(amount_cents), 0),
			COUNT(CASE WHEN status = 'completed' THEN 1 END),
			COUNT(CASE WHEN status = 'pending' THEN 1 END)
		FROM payments
	`).Scan(&stats.TotalRevenue, &stats.CompletedPayments, &stats.PendingPayments)
	return stats, err
}

func (r *Repository) GetUserActivity(ctx context.Context, days int) ([]model.UserActivity, error) {
	query := `
		WITH dates AS (
			SELECT generate_series(
				CURRENT_DATE - INTERVAL '1 day' * $1,
				CURRENT_DATE,
				'1 day'::interval
			) AS date
		)
		SELECT
			d.date::text,
			(SELECT COUNT(*) FROM users WHERE created_at::date = d.date) AS new_users,
			(SELECT COUNT(*) FROM reviews WHERE created_at::date = d.date) AS new_reviews,
			(SELECT COUNT(*) FROM bounties WHERE created_at::date = d.date) AS new_bounties
		FROM dates d
		ORDER BY d.date
	`
	rows, err := r.db.Query(ctx, query, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []model.UserActivity
	for rows.Next() {
		var a model.UserActivity
		if err := rows.Scan(&a.Date, &a.NewUsers, &a.NewReviews, &a.NewBounties); err != nil {
			return nil, err
		}
		activities = append(activities, a)
	}
	return activities, nil
}

func (r *Repository) GetUserBountyCount(ctx context.Context, userID uuid.UUID) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM bounties WHERE user_id = $1", userID).Scan(&count)
	return count, err
}

// Fraud detection methods

func (r *Repository) GetRecentReviewCount(ctx context.Context, userID uuid.UUID, hours int) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM reviews
		WHERE reviewer_id = $1
		AND created_at > NOW() - INTERVAL '1 hour' * $2
	`, userID, hours).Scan(&count)
	return count, err
}

func (r *Repository) FlagReviewForVerification(ctx context.Context, reviewID, amazonLink string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE reviews SET amazon_review_url = $1
		WHERE id = $2
	`, amazonLink, reviewID)
	return err
}

func (r *Repository) GetReviewVerification(ctx context.Context, reviewID string) (*model.ReviewVerification, error) {
	var v model.ReviewVerification
	err := r.db.QueryRow(ctx, `
		SELECT id, COALESCE(amazon_review_url, ''), verified_purchase
		FROM reviews WHERE id = $1
	`, reviewID).Scan(&v.ReviewID, &v.AmazonLink, &v.VerifiedPurchase)
	if err != nil {
		return nil, err
	}
	if v.AmazonLink != "" {
		v.Status = "pending"
	}
	return &v, nil
}

func (r *Repository) Close() {
	r.db.Close()
}
