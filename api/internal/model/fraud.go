package model

// Fraud detection models

type FraudAlert struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Type        string    `json:"type"` // rapid_reviews, suspicious_pattern, fake_account
	Severity    string    `json:"severity"` // low, medium, high
	Description string    `json:"description"`
	ReferenceID string    `json:"reference_id,omitempty"`
	Resolved    bool      `json:"resolved"`
	CreatedAt   string    `json:"created_at"`
}

type ReviewVerification struct {
	ReviewID        string `json:"review_id"`
	VerifiedPurchase bool  `json:"verified_purchase"`
	AmazonLink      string `json:"amazon_link,omitempty"`
	ScreenshotURL   string `json:"screenshot_url,omitempty"`
	Status          string `json:"status"` // pending, verified, rejected
}

type UserReputation struct {
	UserID         string  `json:"user_id"`
	TrustScore     float64 `json:"trust_score"` // 0-100
	ReviewAuthenticity float64 `json:"review_authenticity"` // based on verification
	ActivityScore  float64 `json:"activity_score"` // based on review patterns
}
