package service

import (
	"database/sql"
	"errors"

	"github.com/milyrock/PR-Reviewer/internal/models"
	"github.com/milyrock/PR-Reviewer/internal/repository"
)

type UserService struct {
	repo *repository.Repository
}

func NewUserService(repo *repository.Repository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) SetIsActive(req models.SetIsActiveRequest) (*models.User, error) {
	if err := s.repo.SetUserIsActive(req.UserID, req.IsActive); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	user, err := s.repo.GetUser(req.UserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetReview(userID string) ([]models.PullRequestShort, error) {
	_, err := s.repo.GetUser(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	prs, err := s.repo.GetUserReviewPRs(userID)
	if err != nil {
		return nil, err
	}

	return prs, nil
}
