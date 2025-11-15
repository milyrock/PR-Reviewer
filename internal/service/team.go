package service

import (
	"database/sql"
	"errors"

	"github.com/milyrock/PR-Reviewer/internal/models"
	"github.com/milyrock/PR-Reviewer/internal/repository"
)

type TeamService struct {
	repo *repository.Repository
}

func NewTeamService(repo *repository.Repository) *TeamService {
	return &TeamService{repo: repo}
}

func (s *TeamService) AddTeam(req models.CreateTeamRequest) (*models.Team, error) {
	exists, err := s.repo.TeamExists(req.TeamName)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrTeamExists
	}

	if err := s.repo.CreateTeam(req.TeamName, req.Members); err != nil {
		return nil, err
	}

	team, err := s.repo.GetTeam(req.TeamName)
	if err != nil {
		return nil, err
	}

	return team, nil
}

func (s *TeamService) GetTeam(teamName string) (*models.Team, error) {
	team, err := s.repo.GetTeam(teamName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTeamNotFound
		}
		return nil, err
	}

	return team, nil
}
