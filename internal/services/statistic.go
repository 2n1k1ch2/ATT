package service

import (
	"AvitoTestTask/internal/dto"
	"AvitoTestTask/internal/repository"
	"context"
)

type StatisticServ struct {
	statRepo repository.StatisticRepository
}

func NewStatisticService(statRepo repository.StatisticRepository) StatisticServ {
	return StatisticServ{statRepo: statRepo}
}

func (s *StatisticServ) GetUsersStats(ctx context.Context) (*dto.GetUsersStatsResponse, error) {
	userStats, err := s.statRepo.GetUsersStats(ctx)
	if err != nil {
		return nil, err
	}

	stats := make([]dto.UserStatResponse, len(userStats))
	for i, us := range userStats {
		stats[i] = dto.UserStatResponse{
			UserID:    us.UserID,
			Username:  us.Username,
			TotalPRs:  us.TotalPRs,
			OpenPRs:   us.OpenPRs,
			MergedPRs: us.MergedPRs,
		}
	}

	return &dto.GetUsersStatsResponse{Stats: stats}, nil
}
