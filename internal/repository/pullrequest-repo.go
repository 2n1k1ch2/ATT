package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PullRequestRepo struct {
	db *pgxpool.Pool
}

func NewPullRequestRepo(db *pgxpool.Pool) *PullRequestRepo {
	return &PullRequestRepo{db: db}
}

func (r *PullRequestRepo) CreatePullRequest(pr PullRequest) error {
	_, err := r.db.Exec(context.TODO(), `
        INSERT INTO pull_requests 
        (pull_request_id, pull_request_name, author_id, status, created_at, merged_at) 
        VALUES ($1, $2, $3, $4, $5, $6)`,
		pr.PrID, pr.PrName, pr.AuthorID, pr.Status, pr.CreatedAt, pr.MergedAt)
	if err != nil {
		logger.Error("PullRepo", err.Error())
	}
	return err
}

func (r *PullRequestRepo) GetPullRequest(prId string) (*PullRequest, error) {
	row, err := r.db.Query(context.TODO(), `SELECT pull_request_id,pull_request_name,author_id,status,
       											created_at,merged_at FROM pull_requests WHERE pull_request_id=$1`, prId)
	if err != nil {
		logger.Error("PullRepo", err.Error())
		return nil, err
	}
	defer row.Close()
	if row.Next() {
		pr := PullRequest{}
		err = row.Scan(
			&pr.PrID,
			&pr.PrName,
			&pr.AuthorID,
			&pr.Status,
			&pr.CreatedAt,
			&pr.MergedAt)
		if err != nil {
			logger.Error("PullRepo", err.Error())
			return nil, err
		}
		return &pr, nil
	}
	return nil, errors.New("PullRepo: request not found")

}
func (r *PullRequestRepo) ListPullRequests() ([]PullRequest, error) {
	rows, err := r.db.Query(context.TODO(), `SELECT pull_request_id,pull_request_name,author_id,status,
       											created_at,merged_at FROM pull_requests`)
	if err != nil {
		logger.Error("PullRepo", err.Error())
		return nil, err
	}
	defer rows.Close()
	var prs []PullRequest
	for rows.Next() {
		pr := PullRequest{}
		err = rows.Scan(
			&pr.PrID,
			&pr.PrName,
			&pr.AuthorID,
			&pr.Status,
			&pr.CreatedAt,
			&pr.MergedAt,
		)
		if err != nil {
			logger.Error("PullRepo", err.Error())
			return nil, err
		}
		prs = append(prs, pr)

	}
	return prs, nil
}
func (r *PullRequestRepo) MergePullRequest(prID string) error {
	_, err := r.db.Exec(context.TODO(), `
        UPDATE pull_requests 
        SET status = 'MERGED', merged_at = COALESCE(merged_at, CURRENT_TIMESTAMP)
        WHERE pull_request_id = $1 AND status != 'MERGED'`,
		prID)
	return err
}
