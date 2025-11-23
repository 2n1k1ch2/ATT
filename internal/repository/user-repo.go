package repository

import (
	"AvitoTestTask/internal/dto"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) ChangeActivityFlag(ctx context.Context, userID string, flag bool) (User, error) {
	var u User

	err := r.db.QueryRow(ctx, `
        UPDATE users
        SET is_active = $1
        WHERE user_id = $2
        RETURNING user_id, username, team_name, is_active
    `, flag, userID).Scan(&u.UserID, &u.Username, &u.TeamName, &u.IsActive)

	if err != nil {
		return User{}, fmt.Errorf("update is_active for user %s: %w", userID, err)
	}

	return u, nil
}
func (r *UserRepo) GetUserReviews(ctx context.Context, userID string) ([]PullRequestShort, error) {
	rows, err := r.db.Query(ctx, `
        SELECT 
            pr.pull_request_id,
            pr.pull_request_name,
            pr.author_id,
            pr.status
        FROM pull_requests pr
        INNER JOIN pr_reviewers prr ON pr.pull_request_id = prr.pull_request_id
        WHERE prr.user_id = $1
    `, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []PullRequestShort
	for rows.Next() {
		var pr PullRequestShort
		if err = rows.Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status); err != nil {
			return nil, err
		}
		prs = append(prs, pr)
	}

	return prs, nil
}

func (r *UserRepo) GetActiveUsersByTeam(ctx context.Context, teamName string) ([]User, error) {
	rows, err := r.db.Query(context.TODO(), `SELECT user_id,username,team_name,is_active 
												  FROM users WHERE team_name = $1`, teamName)
	if err != nil {

		return nil, err
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		var user User
		err = rows.Scan(&user.UserID,
			&user.Username,
			&user.TeamName,
			&user.IsActive)
		if err != nil {

			return nil, err

		}
		users = append(users, user)
	}
	return users, nil
}
func (r *UserRepo) GetUser(ctx context.Context, userID string) (*User, error) {
	row := r.db.QueryRow(context.TODO(),
		`SELECT user_id, username, team_name, is_active 
         FROM users WHERE user_id = $1`, userID)

	var user User
	err := row.Scan(&user.UserID,
		&user.Username,
		&user.TeamName,
		&user.IsActive)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.New("User not found")
	}
	return &user, err
}
func (r *TeamRepo) AddUsers(ctx context.Context, teamName string, members []dto.TeamMember) error {
	batch := &pgx.Batch{}

	for _, m := range members {
		batch.Queue(`
			INSERT INTO users (user_id, username, team_name, is_active)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (user_id) DO UPDATE SET
				username = EXCLUDED.username,
				team_name = EXCLUDED.team_name,
				is_active = EXCLUDED.is_active
		`, m.UserID, m.Username, teamName, m.IsActive)
	}

	br := r.db.SendBatch(ctx, batch)
	defer func() {
		err := br.Close()
		if err != nil {
			log.Print("Error closing batch")
		}
	}()
	for range members {
		if _, err := br.Exec(); err != nil {
			return err
		}
	}

	return nil
}
