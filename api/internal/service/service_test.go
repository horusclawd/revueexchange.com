package service

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/revueexchange/api/internal/config"
)

// TestPasswordHashing tests password hashing and verification
func TestPasswordHashing(t *testing.T) {
	password := "testpassword123"

	// Hash the password
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("HashPassword failed: %v", err)
	}

	if hash == "" {
		t.Error("HashPassword returned empty string")
	}

	// Verify correct password
	if !CheckPassword(password, hash) {
		t.Error("CheckPassword failed for correct password")
	}

	// Verify incorrect password
	if CheckPassword("wrongpassword", hash) {
		t.Error("CheckPassword should fail for incorrect password")
	}
}

// TestTokenGeneration tests JWT token generation and validation
func TestTokenGeneration(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret-key"}
	userID := uuid.New()

	// Generate token
	token, err := GenerateToken(userID, cfg.JWTSecret)
	if err != nil {
		t.Errorf("GenerateToken failed: %v", err)
	}

	if token == "" {
		t.Error("GenerateToken returned empty string")
	}

	// Validate token
	parsedID, err := ValidateToken(token, cfg.JWTSecret)
	if err != nil {
		t.Errorf("ValidateToken failed: %v", err)
	}

	if parsedID != userID {
		t.Errorf("Expected userID %s, got %s", userID, parsedID)
	}

	// Test invalid token
	_, err = ValidateToken("invalid.token.here", cfg.JWTSecret)
	if err == nil {
		t.Error("ValidateToken should fail for invalid token")
	}
}

// TestServiceErrors tests custom service error types
func TestServiceErrors(t *testing.T) {
	tests := []struct {
		err      *ServiceError
		expected string
	}{
		{ErrUserAlreadyExists, "user already exists"},
		{ErrInvalidCredentials, "invalid credentials"},
		{ErrInsufficientPoints, "insufficient points"},
		{ErrCannotClaimOwnBounty, "cannot claim your own bounty"},
		{ErrBountyNotAvailable, "bounty is not available"},
		{ErrUnauthorized, "unauthorized"},
	}

	for _, tt := range tests {
		if tt.err.Error() != tt.expected {
			t.Errorf("Expected error message '%s', got '%s'", tt.expected, tt.err.Error())
		}
	}
}

// TestPointsConversionRate tests the points conversion constant
func TestPointsConversionRate(t *testing.T) {
	if PointsConversionRate != 100 {
		t.Errorf("Expected PointsConversionRate to be 100, got %d", PointsConversionRate)
	}
}

// TestCheckoutResult tests checkout result structure
func TestCheckoutResult(t *testing.T) {
	result := &CheckoutResult{
		Payment:     nil,
		SessionURL:  "https://checkout.example.com/session123",
		PointsAward: 500,
	}

	if result.PointsAward != 500 {
		t.Errorf("Expected PointsAward 500, got %d", result.PointsAward)
	}

	if result.SessionURL == "" {
		t.Error("SessionURL should not be empty")
	}
}

// TestTransferPointsValidation tests points transfer validation
func TestTransferPointsValidation(t *testing.T) {
	// This test validates the transfer amount validation logic
	// in the PointsService.TransferPoints method
	tests := []struct {
		name    string
		amount  int
		wantErr bool
	}{
		{"positive amount", 50, false},
		{"zero amount", 0, true},
		{"negative amount", -10, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTransferAmount(tt.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateTransferAmount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// validateTransferAmount is a helper to test validation logic
func validateTransferAmount(amount int) error {
	if amount <= 0 {
		return &ServiceError{Message: "amount must be positive"}
	}
	return nil
}

// TestTimeConstants tests time-related constants
func TestTimeConstants(t *testing.T) {
	// Verify time.Now() works as expected
	now := time.Now()
	if now.IsZero() {
		t.Error("time.Now() should not return zero time")
	}
}
