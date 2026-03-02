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

func (r *Repository) Close() {
	r.db.Close()
}
