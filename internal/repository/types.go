package repository

import "time"

type User struct {
	UserID   string `json:"id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type PullRequest struct {
	// pr - pull request
	PrID      string    `json:"pull_request_id"`
	PrName    string    `json:"pull_request_name"`
	AuthorID  string    `json:"author_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	MergedAt  time.Time `json:"merged_at"`
}
