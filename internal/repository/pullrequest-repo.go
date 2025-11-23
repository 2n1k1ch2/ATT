package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

type PullRequestRepo struct {
	db *pgxpool.Pool
}

func NewPullRequestRepo(db *pgxpool.Pool) *PullRequestRepo {
	return &PullRequestRepo{db: db}
}

func (r *PullRequestRepo) CreatePullRequest(ctx context.Context, prID, authorID, prName string) (PullRequest, error) {
	t := time.Now()
	pr := PullRequest{
		PrID:      prID,
		PrName:    prName,
		AuthorID:  authorID,
		Status:    "OPEN",
		CreatedAt: &t,
		MergedAt:  nil,
	}
	_, err := r.db.Exec(context.TODO(), `
        INSERT INTO pull_requests 
        (pull_request_id, pull_request_name, author_id, status, created_at, merged_at) 
        VALUES ($1, $2, $3, $4, $5, $6)`,
		pr.PrID, pr.PrName, pr.AuthorID, pr.Status, pr.CreatedAt, pr.MergedAt)
	if err != nil {
		log.Print(err.Error())
	}
	return pr, err
}

func (r *PullRequestRepo) GetPullRequest(ctx context.Context, prId string) (PullRequest, error) {
	row, err := r.db.Query(context.TODO(), `SELECT pull_request_id,pull_request_name,author_id,status,
       											created_at,merged_at FROM pull_requests WHERE pull_request_id=$1`, prId)
	pr := PullRequest{}
	if err != nil {
		return pr, err
	}

	defer row.Close()
	if row.Next() {

		err = row.Scan(
			&pr.PrID,
			&pr.PrName,
			&pr.AuthorID,
			&pr.Status,
			&pr.CreatedAt,
			&pr.MergedAt)
		if err != nil {

			return PullRequest{}, err
		}
		return pr, nil
	}
	return pr, errors.New("PullRepo: request not found")

}
func (r *PullRequestRepo) ListPullRequests(ctx context.Context) ([]PullRequest, error) {
	rows, err := r.db.Query(context.TODO(), `SELECT pull_request_id,pull_request_name,author_id,status,
       											created_at,merged_at FROM pull_requests`)
	if err != nil {

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
			return nil, err
		}
		prs = append(prs, pr)

	}
	return prs, nil
}
func (r *PullRequestRepo) MergePullRequest(ctx context.Context, prID string) (PullRequest, error) {
	_ = r.db.QueryRow(ctx, `
        UPDATE pull_requests 
        SET status = 'MERGED', merged_at = COALESCE(merged_at, CURRENT_TIMESTAMP)
        WHERE pull_request_id = $1 AND status != 'MERGED'`,
		prID)

	pr, err := r.GetPullRequest(ctx, prID)

	if err != nil {
		log.Print(err.Error())
	}
	return pr, err
}
func (r *PullRequestRepo) ListPRsForReviewer(ctx context.Context, userID string) ([]PullRequest, error) {
	rows, err := r.db.Query(ctx, `
		SELECT p.pr_id, p.title, p.author_id, p.created_at
		FROM pull_requests p
		JOIN pr_reviewers r ON p.pr_id = r.pr_id
		WHERE r.user_id = $1
	`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []PullRequest
	for rows.Next() {
		var pr PullRequest
		if err = rows.Scan(&pr.AuthorID, &pr.PrName, &pr.Status, &pr.CreatedAt); err != nil {
			return nil, err
		}
		prs = append(prs, pr)
	}
	return prs, nil
}
