package repository

import (
	"time"
)

type User struct {
	UserID   string
	Username string
	TeamName string
	IsActive bool
}

type PullRequest struct {
	// pr - pull request
	PrID      string
	PrName    string
	AuthorID  string
	Status    string
	CreatedAt *time.Time
	MergedAt  *time.Time
}

type PullRequestShort struct {
	PullRequestID   string
	PullRequestName string
	AuthorID        string
	Status          string
}

type UserStatsResponse struct {
	UserID    string
	Username  string
	TotalPRs  int // всего назначенных PR
	OpenPRs   int // открытых PR
	MergedPRs int // мержнутых PR

}
