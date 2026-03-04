package model

// Analytics models

type OverviewStats struct {
	TotalUsers         int `json:"total_users"`
	TotalBounties      int `json:"total_bounties"`
	TotalReviews       int `json:"total_reviews"`
	TotalPointsAwarded int `json:"total_points_awarded"`
	TotalPointsSpent   int `json:"total_points_spent"`
}

type BountyMetrics struct {
	Status    string `json:"status"`
	Count     int    `json:"count"`
	TotalBounties int `json:"total_bounties"`
}

type ReviewMetrics struct {
	Status     string `json:"status"`
	Count      int    `json:"count"`
	AvgRating  float64 `json:"avg_rating"`
	TotalWords int    `json:"total_words"`
}

type RevenueStats struct {
	TotalRevenue    int `json:"total_revenue"`
	CompletedPayments int `json:"completed_payments"`
	PendingPayments int `json:"pending_payments"`
}

type UserActivity struct {
	Date       string `json:"date"`
	NewUsers   int    `json:"new_users"`
	NewReviews int    `json:"new_reviews"`
	NewBounties int  `json:"new_bounties"`
}
