package service

import (
	"AvitoTestTask/internal/repository"
	"context"
)

type PullRequestServ struct {
	pr repository.PullRequestRepository
}

func NewUserService(pr repository.PullRequestRepository) PullRequestServ {
	return PullRequestServ{pr: pr}
}

func (s *PullRequestServ) GetPR(ctx context.Context, prID string) (*repository.PullRequest, error) {
	pr, err := s.pr.GetPullRequest(prID)
	if err != nil {
		return nil, err
	}
	return pr, nil

}
func (s *PullRequestServ) CreatePR(ctx context.Context, pr repository.PullRequest) error {
	err := s.pr.CreatePullRequest(pr)
	if err != nil {
		return err
	}
	return nil
}
func (s *PullRequestServ) ListPRs(ctx context.Context) ([]repository.PullRequest, error) {
	prs, err := s.pr.ListPullRequests()
	if err != nil {
		return nil, err

	}
	return prs, nil
}
func (s *PullRequestServ) MergePR(ctx context.Context, prID string) error {
	err := s.pr.MergePullRequest(prID)
	if err != nil {
		return err
	}
	return nil
}
