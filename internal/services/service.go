package service

import (
	"AvitoTestTask/internal/repository"
	"context"
	"log/slog"
)

var logger slog.Logger

func init() {
	logger = slog.Logger{}
}

type UserService interface {
	GetUser(ctx context.Context, userID string) (*repository.User, error)
	ListTeamUsers(ctx context.Context, teamName string) ([]repository.User, error)
}

type PullRequestService interface {
	CreatePR(ctx context.Context, pr repository.PullRequest) error
	GetPR(ctx context.Context, prID string) (*repository.PullRequest, error)
	ListPRs(ctx context.Context) ([]repository.PullRequest, error)
	MergePR(ctx context.Context, prID string) error
}

type PRReviewerService interface {
	AssignReviewers(ctx, prID string) ([]repository.User, error)
	ReassignReviewer(ctx, prID string, oldReviewerID string) (repository.User, error)
	GetUserAssignedPR(ctx, userID string) ([]repository.PullRequest, error)
}
