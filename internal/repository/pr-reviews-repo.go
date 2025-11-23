package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PrReviewsRepo struct {
	db *pgxpool.Pool
}

func NewPrReviewsRepo(db *pgxpool.Pool) *PrReviewsRepo {
	return &PrReviewsRepo{db: db}
}

func (r *PrReviewsRepo) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.db.BeginTx(ctx, pgx.TxOptions{})
}

func (r *PrReviewsRepo) AddReviewer(ctx context.Context, tx pgx.Tx, prID, userID string) error {

	_, err := tx.Exec(ctx, `
		INSERT INTO pr_reviewers (pull_request_id, user_id, assigned_at)
		VALUES ($1, $2, NOW())
	`, prID, userID)
	return err
}

func (r *PrReviewsRepo) RemoveReviewer(ctx context.Context, tx pgx.Tx, prID, userID string) error {
	cmdTag, err := tx.Exec(ctx, `
		DELETE FROM pr_reviewers
		WHERE pull_request_id = $1 AND user_id = $2
	`, prID, userID)

	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("reviewer not found")
	}

	return nil
}
func (r *PrReviewsRepo) GetReviewers(ctx context.Context, tx pgx.Tx, prID string) ([]User, error) {
	rows, err := tx.Query(ctx, `
		SELECT u.user_id, u.username, u.team_name, u.is_active
		FROM pr_reviewers pr
		JOIN users u ON pr.user_id = u.user_id
		WHERE pr.pull_request_id = $1
	`, prID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err = rows.Scan(&u.UserID, &u.Username, &u.TeamName, &u.IsActive); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *PrReviewsRepo) CountReviewers(ctx context.Context, tx pgx.Tx, prID string) (uint, error) {
	row := tx.QueryRow(ctx, `
		SELECT COUNT(*) 
		FROM pr_reviewers 
		WHERE pull_request_id = $1
	`, prID)

	var count uint
	return count, row.Scan(&count)
}
