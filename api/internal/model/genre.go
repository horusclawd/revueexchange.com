package model

// Genre and expertise models

type UserGenre struct {
	UserID     string   `json:"user_id"`
	Genres     []string `json:"genres"` // fiction, non-fiction, tech, business, etc.
	Expertise  []string `json:"expertise"` // topics user is expert in
	Interests  []string `json:"interests"` // topics user interested in
}

type GenrePreference struct {
	ID        string   `json:"id"`
	UserID    string   `json:"user_id"`
	Genre     string   `json:"genre"`
	Preference string  `json:"preference"` // expert, interested, none
	Weight    int      `json:"weight"` // 1-10
}
