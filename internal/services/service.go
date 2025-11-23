package service

import (
	"AvitoTestTask/internal/dto"
	"context"
	"log/slog"
)

var logger slog.Logger

func init() {
	logger = slog.Logger{}
}

type TeamService interface {
	CreateTeam(ctx context.Context, teamName string, members []dto.TeamMember) (*dto.CreateTeamResponse, error)
	ListTeamUsers(ctx context.Context, teamName string) (dto.GetTeamResponse, error)
}

type UserService interface {
	ChangeActivityFlag(ctx context.Context, user_id string, flag bool) (*dto.SetUserIsActiveResponse, error)
	GetUserAssignedPR(ctx context.Context, user_id string) (dto.GetUserReviewsResponse, error)
}

type PullRequestService interface {
	CreatePR(ctx context.Context, prID, authorID, prName string) (dto.PullRequest, error)
	MergePR(ctx context.Context, prID string) (dto.PullRequest, error)
	ReassignReviewer(ctx context.Context, prID string, oldReviewerID string) (dto.PullRequest, string, error)
}
type StatisticService interface {
	GetUsersStats(ctx context.Context) (*dto.GetUsersStatsResponse, error)
}
