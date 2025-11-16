package service

import (
	"github.com/milyrock/PR-Reviewer/internal/models"
	"github.com/milyrock/PR-Reviewer/internal/repository"
)

type StatisticsService struct {
	repo *repository.Repository
}

func NewStatisticsService(repo *repository.Repository) *StatisticsService {
	return &StatisticsService{repo: repo}
}

func (s *StatisticsService) GetStatistics() (*models.StatisticsResponse, error) {
	userStats, err := s.repo.GetUserReviewStats()
	if err != nil {
		return nil, err
	}

	prStats, err := s.repo.GetPRReviewStats()
	if err != nil {
		return nil, err
	}

	return &models.StatisticsResponse{
		UserStats: userStats,
		PRStats:   prStats,
	}, nil
}
