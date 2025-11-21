package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TeamRepo struct {
	db *pgxpool.Pool
}

func NewTeamRepo(db *pgxpool.Pool) *TeamRepo {
	return &TeamRepo{db: db}
}
func (r *TeamRepo) GetTeamUsers(teamName string) ([]User, error) {
	rows, err := r.db.Query(context.TODO(), `SELECT user_id,username,team_name,is_active 
													FROM users WHERE team_name = $1`, teamName)
	if err != nil {
		logger.Error("GetTeamUsers err:", err)
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
			logger.Error("GetTeamUsers err:", err)
			return nil, err
		}
		users = append(users, user)

	}
	return users, nil
}
