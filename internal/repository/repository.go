package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"log/slog"
)

var logger slog.Logger

func init() {
	logger = slog.Logger{}
}

type TeamRepository interface {
	GetTeamUsers(teamName string) ([]User, error) // /team/get — список участников команды
}

type UserRepository interface {
	GetUser(userID string) (*User, error)                 // получение  пользователя
	GetActiveUsersByTeam(teamName string) ([]User, error) // выбор активных кандидатов в ревьюверы
	GetUserReviews(userID string) ([]PullRequest, error)  // /users/getReview — список PR, где user ревьювер
}

type PullRequestRepository interface {
	CreatePullRequest(pr PullRequest) error           // создание PR
	GetPullRequest(prID string) (*PullRequest, error) // получение одного PR
	ListPullRequests() ([]PullRequest, error)         // /prs/get — список всех PR
	MergePullRequest(prID string) error               // merge PR (смена статуса + merged_at)
}

type PRReviewersRepository interface {
	BeginTx(ctx context.Context) (pgx.Tx, error)

	AddReviewer(ctx context.Context, tx pgx.Tx, prID, userID string) error    // Добавить минимально нужное число ревьюеров
	RemoveReviewer(ctx context.Context, tx pgx.Tx, prID, userID string) error // Убрать и заменить на другого ревьюера
	GetReviewers(ctx context.Context, tx pgx.Tx, prID string) ([]User, error) // Узнать ревьюеров по PR
	CountReviewers(ctx context.Context, tx pgx.Tx, prID string) (uint, error) // Подсчитать сколько ревьюеров у PR
}
