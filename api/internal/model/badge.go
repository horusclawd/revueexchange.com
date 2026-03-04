package model

import "time"

// Badge represents a user badge
type Badge struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	BadgeType   string    `json:"badge_type"`
	BadgeName   string    `json:"badge_name"`
	Description string    `json:"description"`
	Tier        string    `json:"tier"` // bronze, silver, gold, platinum
	IconURL     string    `json:"icon_url,omitempty"`
	AwardedAt   time.Time `json:"awarded_at"`
}

// BadgeDefinition defines badge types
type BadgeDefinition struct {
	Type        string   `json:"type"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tiers       []string `json:"tiers"` // available tiers
	Condition   string   `json:"condition"` // first_review, review_10, review_50, top_reviewer, streak_7, streak_30
}

// Badge definitions
var BadgeDefinitions = []BadgeDefinition{
	{
		Type:        "first_review",
		Name:        "First Review",
		Description: "Submitted your first review",
		Tiers:       []string{"bronze"},
		Condition:   "first_review",
	},
	{
		Type:        "review_10",
		Name:        "Prolific Reviewer",
		Description: "Submitted 10 reviews",
		Tiers:       []string{"bronze", "silver", "gold"},
		Condition:   "review_10",
	},
	{
		Type:        "review_50",
		Name:        "Master Reviewer",
		Description: "Submitted 50 reviews",
		Tiers:       []string{"silver", "gold", "platinum"},
		Condition:   "review_50",
	},
	{
		Type:        "streak_7",
		Name:        "Week Warrior",
		Description: "7-day activity streak",
		Tiers:       []string{"bronze", "silver"},
		Condition:   "streak_7",
	},
	{
		Type:        "streak_30",
		Name:        "Monthly Master",
		Description: "30-day activity streak",
		Tiers:       []string{"gold", "platinum"},
		Condition:   "streak_30",
	},
	{
		Type:        "top_reviewer",
		Name:        "Top Reviewer",
		Description: "Reached top reviewer status",
		Tiers:       []string{"gold", "platinum"},
		Condition:   "top_reviewer",
	},
	{
		Type:        "helpful_100",
		Name:        "Helpful Hand",
		Description: "Received 100 helpful votes",
		Tiers:       []string{"silver", "gold"},
		Condition:   "helpful_100",
	},
	{
		Type:        "early_adopter",
		Name:        "Early Adopter",
		Description: "Joined during beta",
		Tiers:       []string{"bronze"},
		Condition:   "early_adopter",
	},
}

// UserStats represents user statistics for badge checking
type UserStats struct {
	UserID          string
	ReviewCount     int
	StreakDays      int
	HelpfulVotes    int
	IsTopReviewer   bool
	MemberSince     time.Time
}

// LeaderboardEntry represents a user on the leaderboard
type LeaderboardEntry struct {
	UserID         string    `json:"user_id"`
	Username       string    `json:"username"`
	DisplayName    string    `json:"display_name"`
	Points         int       `json:"points"`
	ReviewCount    int       `json:"review_count"`
	Rank           int       `json:"rank"`
	LastUpdated    time.Time `json:"last_updated"`
}

// Streak represents a user's activity streak
type Streak struct {
	UserID         string    `json:"user_id"`
	CurrentStreak  int       `json:"current_streak"`
	LongestStreak  int       `json:"longest_streak"`
	LastActivityAt time.Time `json:"last_activity_at"`
	StreakStartedAt time.Time `json:"streak_started_at"`
}
