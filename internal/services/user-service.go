package service

import (
	"AvitoTestTask/internal/repository"
	"context"
)

type UserServ struct {
	users repository.UserRepository
	team  repository.TeamRepository
}

func NewUserServ(users repository.UserRepository, team repository.TeamRepository) UserServ {
	return UserServ{users: users, team: team}
}

func (s *UserServ) GetUser(ctx context.Context, userID string) (*repository.User, error) {
	user, err := s.users.GetUser(userID)
	return user, err
}

func (s *UserServ) ListTeamUsers(ctx context.Context, teamName string) ([]repository.User, error) {
	users, err := s.team.GetTeamUsers(teamName)
	if err != nil {
		return nil, err
	}
	return users, nil
}
