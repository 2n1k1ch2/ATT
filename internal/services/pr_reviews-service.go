package service

import (
	"AvitoTestTask/internal/repository"
	"context"
	"errors"
	"math/rand"
)

type PrReviewsServ struct {
	reviews repository.PRReviewersRepository
	users   repository.UserRepository
	pr      repository.PullRequestRepository
}

const MAX_REVIEWS = 2

var (
	ErrEnoughReviewers               = errors.New("pr already has 2 reviewers")
	ErrAlreadyZeroErrEnoughReviewers = errors.New("no reviewers to remove")
	ErrAlreadyMerged                 = errors.New("pr already merged")
)

func NewPrReviewsService(
	reviews repository.PRReviewersRepository,
	users repository.UserRepository,
	pr repository.PullRequestRepository,
) PrReviewsServ {
	return PrReviewsServ{reviews: reviews, users: users, pr: pr}
}

func (s *PrReviewsServ) AssignReviewers(ctx context.Context, prID string) ([]repository.User, error) {
	tx, err := s.reviews.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			logger.Error("warning: transaction rollback failed: %v", rollbackErr)
		}
	}()

	// Получаем PR
	pr, err := s.pr.GetPullRequest(prID)
	if err != nil {
		return nil, err
	}
	if pr.Status == "Merged" {
		return nil, ErrAlreadyMerged
	}

	// Получаем доступных пользователей
	users, err := s.users.GetActiveUsersByTeam(pr.AuthorID)
	if err != nil {
		return nil, err
	}
	users = removeUser(users, pr.AuthorID)

	// Уже назначенные ревьюверы
	reviewers, err := s.reviews.GetReviewers(ctx, tx, prID)
	if err != nil {
		return nil, err
	}
	if len(reviewers) >= MAX_REVIEWS {
		return nil, ErrEnoughReviewers
	}

	// Убираем тех, кто уже ревьюверы
	for _, r := range reviewers {
		users = removeUser(users, r.UserID)
	}

	newReviewers := selectRandomUnique(users, MAX_REVIEWS-len(reviewers))

	// Сохраняем новых ревьюверов по транзакциям
	for _, reviewer := range newReviewers {
		if err = s.reviews.AddReviewer(ctx, tx, prID, reviewer.UserID); err != nil {
			return nil, err
		}
	}

	return newReviewers, tx.Commit(ctx)
}

func (s *PrReviewsServ) RemoveReviewer(ctx context.Context, prID, userID string) error {
	tx, err := s.reviews.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			logger.Error("warning: transaction rollback failed: %v", rollbackErr)
		}
	}()
	//Получаем ревьюеров
	reviewers, err := s.reviews.GetReviewers(ctx, tx, prID)
	if err != nil {
		return err
	}

	if len(reviewers) == 0 {
		return ErrAlreadyZeroErrEnoughReviewers
	}
	//Пытаемся удалить ревьюера
	if err = s.reviews.RemoveReviewer(ctx, tx, prID, userID); err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	// Пытаемся назначить нового ревьюера после транзакции
	_, err = s.AssignReviewers(ctx, prID)
	return err
}

// UTILS
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
