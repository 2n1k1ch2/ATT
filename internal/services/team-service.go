package service

import (
	"AvitoTestTask/internal/dto"
	"AvitoTestTask/internal/repository"
	"context"
)

type TeamServ struct {
	team repository.TeamRepository
	user repository.UserRepository
}

func NewTeamServ(team repository.TeamRepository, user repository.UserRepository) TeamServ {
	return TeamServ{team: team, user: user}
}

func (s *TeamServ) CreateTeam(ctx context.Context, teamName string, members []dto.TeamMember) (*dto.CreateTeamResponse, error) {
	err := s.team.CreateTeam(ctx, teamName)

	if err != nil {

		logger.Error("TeamService ", "error", err.Error())
		return nil, err
	}

	err = s.team.AddUsers(ctx, teamName, members)
	if err != nil {
		logger.Error("TeamService ", "error", err.Error())
		return nil, err
	}

	return &dto.CreateTeamResponse{
		Team: dto.Team{
			TeamName: teamName,
			Members:  members,
		},
	}, nil

}
func (s *TeamServ) ListTeamUsers(ctx context.Context, teamName string) (dto.GetTeamResponse, error) {
	users, err := s.team.GetTeamUsers(ctx, teamName)
	if err != nil {
		logger.Error("TeamService ", "error", err.Error())
		return dto.GetTeamResponse{}, err
	}

	members := make([]dto.TeamMember, 0, len(users))
	for _, u := range users {
		members = append(members, dto.TeamMember{
			UserID:   u.UserID,
			Username: u.Username,
			IsActive: u.IsActive,
		})
	}

	return dto.GetTeamResponse{
		TeamName: teamName,
		Members:  members,
	}, nil
}
