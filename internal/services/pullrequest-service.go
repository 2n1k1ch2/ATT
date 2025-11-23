package service

import (
	"AvitoTestTask/internal/dto"
	"AvitoTestTask/internal/repository"
	"context"
	"errors"
	"math/rand"
)

const MAX_REVIEWS = 2

var (
	ErrEnoughReviewers               = errors.New("pr already has 2 reviewers")
	ErrAlreadyZeroErrEnoughReviewers = errors.New("no reviewers to remove")
	ErrAlreadyMerged                 = errors.New("pr already merged")
)

type PullRequestServ struct {
	pr      repository.PullRequestRepository
	reviews repository.PRReviewersRepository
	users   repository.UserRepository
}

func NewPullRequestServ(pr repository.PullRequestRepository, reviews repository.PRReviewersRepository, users repository.UserRepository) PullRequestServ {
	return PullRequestServ{pr: pr,
		reviews: reviews,
		users:   users,
	}
}

func (s *PullRequestServ) CreatePR(ctx context.Context, prID, authorID, prName string) (dto.PullRequest, error) {
	// создаем PR
	pr, err := s.pr.CreatePullRequest(ctx, prID, authorID, prName)
	if err != nil {
		return dto.PullRequest{}, err
	}
	author, err := s.users.GetUser(ctx, pr.AuthorID)
	if err != nil {
		return dto.PullRequest{}, err
	}

	// старт транзакции
	tx, err := s.reviews.BeginTx(ctx)
	if err != nil {
		return dto.PullRequest{}, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	// доступные пользователи
	users, err := s.users.GetActiveUsersByTeam(ctx, author.TeamName)
	if err != nil {
		return dto.PullRequest{}, err
	}
	users = removeUser(users, pr.AuthorID)

	// выбираем
	newReviewers := selectRandomUnique(users, MAX_REVIEWS)

	// назначаем
	for _, reviewer := range newReviewers {
		if err = s.reviews.AddReviewer(ctx, tx, prID, reviewer.UserID); err != nil {
			return dto.PullRequest{}, err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return dto.PullRequest{}, err
	}

	return s.getPullRequest(ctx, prID)
}

func (s *PullRequestServ) MergePR(ctx context.Context, prID string) (dto.PullRequest, error) {
	pr, err := s.pr.GetPullRequest(ctx, prID)
	if err != nil {
		return dto.PullRequest{}, err
	}
	if pr.Status == "MERGED" {
		return dto.PullRequest{}, ErrAlreadyMerged
	}
	_, err = s.pr.MergePullRequest(ctx, prID)
	if err != nil {
		return dto.PullRequest{}, err
	}
	resp, err := s.getPullRequest(ctx, prID)
	if err != nil {
		return dto.PullRequest{}, err
	}

	return resp, nil

}

func (s *PullRequestServ) ReassignReviewer(ctx context.Context, prID string, oldReviewerID string) (dto.PullRequest, string, error) {

	pr, err := s.pr.GetPullRequest(ctx, prID)
	if err != nil {
		return dto.PullRequest{}, "", err
	}
	if pr.Status == "MERGED" {
		return dto.PullRequest{}, "", ErrAlreadyMerged
	}
	tx, err := s.reviews.BeginTx(ctx)
	if err != nil {
		return dto.PullRequest{}, "", err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()
	reviewers, err := s.reviews.GetReviewers(ctx, tx, prID)
	if err != nil {
		return dto.PullRequest{}, "", err
	}

	if len(reviewers) == 0 {
		return dto.PullRequest{}, "", ErrAlreadyZeroErrEnoughReviewers
	}

	if err = s.reviews.RemoveReviewer(ctx, tx, prID, oldReviewerID); err != nil {
		return dto.PullRequest{}, "", err
	}

	users, err := s.users.GetActiveUsersByTeam(ctx, reviewers[0].TeamName)
	if err != nil {
		return dto.PullRequest{}, "", err
	}

	pr, err = s.pr.GetPullRequest(ctx, prID)
	if err != nil {
		return dto.PullRequest{}, "", err
	}
	users = removeUser(users, pr.AuthorID)

	for _, r := range reviewers {
		users = removeUser(users, r.UserID)
	}

	cand := selectRandomUnique(users, 1)
	if len(cand) > 0 {
		if err = s.reviews.AddReviewer(ctx, tx, prID, cand[0].UserID); err != nil {
			return dto.PullRequest{}, "", err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return dto.PullRequest{}, "", err
	}

	resp, err := s.getPullRequest(ctx, prID)
	if err != nil {
		return dto.PullRequest{}, "", err
	}

	return resp, cand[0].UserID, nil
}

// UTILS

func (s *PullRequestServ) getPullRequest(ctx context.Context, prID string) (dto.PullRequest, error) {

	pr, err := s.pr.GetPullRequest(ctx, prID)
	if err != nil {
		return dto.PullRequest{}, err

	}
	tx, err := s.reviews.BeginTx(ctx)
	if err != nil {
		return dto.PullRequest{}, err
	}
	reviewers, err := s.reviews.GetReviewers(ctx, tx, prID)
	if err != nil {
		return dto.PullRequest{}, err
	}
	if len(reviewers) == 0 {
		return dto.PullRequest{}, ErrAlreadyZeroErrEnoughReviewers
	}
	resp := dto.PullRequest{
		AuthorID:          pr.AuthorID,
		PullRequestID:     prID,
		PullRequestName:   pr.PrName,
		CreatedAt:         pr.CreatedAt,
		MergedAt:          pr.MergedAt,
		AssignedReviewers: make([]string, len(reviewers)),
	}

	for i, reviewer := range reviewers {
		resp.AssignedReviewers[i] = reviewer.UserID
	}
	return resp, tx.Commit(ctx)
}

func removeUser(items []repository.User, userID string) []repository.User {
	for i, u := range items {
		if u.UserID == userID {
			items[i] = items[len(items)-1]
			return items[:len(items)-1]
		}
	}
	return items
}

func selectRandomUnique(items []repository.User, count int) []repository.User {
	if count <= 0 {
		return []repository.User{}
	}
	if count >= len(items) {
		return items
	}
	rand.Shuffle(len(items), func(i, j int) {
		items[i], items[j] = items[j], items[i]
	})
	return items[:count]
}
