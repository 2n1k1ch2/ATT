package repository

import (
	"AvitoTestTask/internal/dto"

	"context"
	"github.com/jackc/pgx/v5"
)

type TeamRepository interface {
	GetTeamUsers(ctx context.Context, teamName string) ([]User, error)             // получение юзеров по команде
	CreateTeam(ctx context.Context, teamName string) error                         //создание команды
	AddUsers(ctx context.Context, teamName string, members []dto.TeamMember) error // добавляет юзеров одним батчем
}

type UserRepository interface {
	ChangeActivityFlag(ctx context.Context, userID string, flag bool) (User, error) //изменить флаг активности
	GetUserReviews(ctx context.Context, userID string) ([]PullRequestShort, error)  // список PR, где user ревьювер

	GetUser(ctx context.Context, userID string) (*User, error)                 // получение  пользователя
	GetActiveUsersByTeam(ctx context.Context, teamName string) ([]User, error) // выбор активных кандидатов в ревьюверы

}

type PullRequestRepository interface {
	CreatePullRequest(ctx context.Context, prID, author_id, pr_name string) (PullRequest, error) // создание PR
	GetPullRequest(ctx context.Context, prID string) (PullRequest, error)                        // получение одного PR
	MergePullRequest(ctx context.Context, prID string) (PullRequest, error)
	ListPullRequests(ctx context.Context) ([]PullRequest, error) // список всех PR
	// merge PR (смена статуса + merged_at)
}

type PRReviewersRepository interface {
	BeginTx(ctx context.Context) (pgx.Tx, error)
	AddReviewer(ctx context.Context, tx pgx.Tx, prID, userID string) error    // Добавить минимально нужное число ревьюеров
	RemoveReviewer(ctx context.Context, tx pgx.Tx, prID, userID string) error // Убрать и заменить на другого ревьюера
	GetReviewers(ctx context.Context, tx pgx.Tx, prID string) ([]User, error) // Узнать ревьюеров по PR
	CountReviewers(ctx context.Context, tx pgx.Tx, prID string) (uint, error) // Подсчитать сколько ревьюеров у PR
}
type StatisticRepository interface {
	GetUsersStats(ctx context.Context) ([]UserStatsResponse, error) //статистика юзеров
}

type Tx interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
