package service

import (
	"AvitoTestTask/internal/dto"
	"AvitoTestTask/internal/repository"
	"context"
)

type UserServ struct {
	users repository.UserRepository
	team  repository.TeamRepository
}

func NewUserServ(users repository.UserRepository, team repository.TeamRepository) UserServ {
	return UserServ{users: users, team: team}
}

func (s *UserServ) ChangeActivityFlag(ctx context.Context, user_id string, flag bool) (*dto.SetUserIsActiveResponse, error) {
	user, err := s.users.ChangeActivityFlag(ctx, user_id, flag)
	if err != nil {
		return nil, err
	}

	return &dto.SetUserIsActiveResponse{
		User: dto.User{
			UserID:   user.UserID,
			Username: user.Username,
			TeamName: user.TeamName,
			IsActive: flag,
		},
	}, nil

}

func (s *UserServ) GetUserAssignedPR(ctx context.Context, user_id string) (dto.GetUserReviewsResponse, error) {
	prs, err := s.users.GetUserReviews(ctx, user_id)
	if err != nil {
		return dto.GetUserReviewsResponse{}, err
	}
	res := dto.GetUserReviewsResponse{
		UserID:       user_id,
		PullRequests: make([]dto.PullRequestShort, len(prs)),
	}

	for i, pr := range prs {
		res.PullRequests[i] = dto.PullRequestShort{
			PullRequestID:   pr.PullRequestID,
			PullRequestName: pr.PullRequestName,
			AuthorID:        pr.AuthorID,
			Status:          pr.Status,
		}
	}
	return res, nil
}
