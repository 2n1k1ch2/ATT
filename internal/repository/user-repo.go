package repository

import (
	"context"
	"errors"
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

func (r *UserRepo) GetUser(userID string) (*User, error) {
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

func (r *UserRepo) GetActiveUsersByTeam(teamName string) ([]User, error) {
	rows, err := r.db.Query(context.TODO(), `SELECT user_id,username,team_name,is_active 
												  FROM users WHERE team_name = $1`, teamName)
	if err != nil {
		log.Fatalf("User-Repository: %s", err.Error())
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
			logger.Error("UserRepo GetActiveUsersByTeam Error", err.Error())
			return nil, err

		}
		users = append(users, user)
	}
	return users, nil
}
func (r *UserRepo) GetUserReviews(userID string) ([]PullRequest, error) {
	rows, err := r.db.Query(context.TODO(), `SELECT pull_request_id, pull_request_name,author_id,status,created_at,
       											 merged_at FROM pull_requests WHERE author_id = $1`, userID)
	if err != nil {
		logger.Error("UserRepo GetUserReviews Error ", err.Error())
		return nil, err

	}
	defer rows.Close()
	var pullRequests []PullRequest
	for rows.Next() {
		var pullRequest PullRequest
		err = rows.Scan(&pullRequest.PrID,
			&pullRequest.PrName,
			&pullRequest.AuthorID,
			&pullRequest.Status,
			&pullRequest.CreatedAt,
			&pullRequest.MergedAt)
		if err != nil {
			logger.Error("UserRepo GetUserReviews Error ", err)
			return nil, err
		}
		pullRequests = append(pullRequests, pullRequest)
	}

	return pullRequests, nil
}
