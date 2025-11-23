package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StatisticRepo struct {
	db *pgxpool.Pool
}

func NewStatisticRepo(db *pgxpool.Pool) *StatisticRepo {
	return &StatisticRepo{db: db}
}
func (r *StatisticRepo) GetUsersStats(ctx context.Context) ([]UserStatsResponse, error) {
	rows, err := r.db.Query(ctx, `
        SELECT 
			u.user_id,
			u.username,
			COUNT(prr.pull_request_id) AS total_prs,
			COUNT(*) FILTER (WHERE pr.status = 'OPEN') AS open_prs,
			COUNT(*) FILTER (WHERE pr.status = 'MERGED') AS merged_prs
		FROM users u
		LEFT JOIN pr_reviewers prr ON u.user_id = prr.user_id
		LEFT JOIN pull_requests pr ON prr.pull_request_id = pr.pull_request_id
		GROUP BY u.user_id, u.username
		ORDER BY total_prs DESC;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []UserStatsResponse
	for rows.Next() {
		var s UserStatsResponse
		err = rows.Scan(&s.UserID, &s.Username, &s.TotalPRs, &s.OpenPRs, &s.MergedPRs)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}
