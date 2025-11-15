package service

import (
	"database/sql"
	"errors"
	"math/rand"

	"github.com/milyrock/PR-Reviewer/internal/models"
	"github.com/milyrock/PR-Reviewer/internal/repository"
)

type PRService struct {
	repo *repository.Repository
}

func NewPRService(repo *repository.Repository) *PRService {
	return &PRService{repo: repo}
}

func (s *PRService) CreatePR(req models.CreatePRRequest) (*models.PullRequest, error) {
	exists, err := s.repo.PRExists(req.PullRequestID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrPRExists
	}

	author, err := s.repo.GetUser(req.AuthorID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	candidates, err := s.repo.GetActiveUsersByTeamName(author.TeamName, req.AuthorID)
	if err != nil {
		return nil, err
	}

	reviewerIDs := selectReviewers(candidates, 2)

	pr := &models.PullRequest{
		PullRequestID:     req.PullRequestID,
		PullRequestName:   req.PullRequestName,
		AuthorID:          req.AuthorID,
		Status:            "OPEN",
		AssignedReviewers: reviewerIDs,
	}

	if err := s.repo.CreatePR(pr, reviewerIDs); err != nil {
		return nil, err
	}

	createdPR, err := s.repo.GetPR(req.PullRequestID)
	if err != nil {
		return nil, err
	}

	return createdPR, nil
}

func (s *PRService) MergePR(req models.MergePRRequest) (*models.PullRequest, error) {
	_, err := s.repo.GetPR(req.PullRequestID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPRNotFound
		}
		return nil, err
	}

	if err := s.repo.MergePR(req.PullRequestID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPRNotFound
		}
		return nil, err
	}

	pr, err := s.repo.GetPR(req.PullRequestID)
	if err != nil {
		return nil, err
	}

	return pr, nil
}

func (s *PRService) ReassignPR(req models.ReassignPRRequest) (*models.PullRequest, string, error) {
	pr, err := s.repo.GetPR(req.PullRequestID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", ErrPRNotFound
		}
		return nil, "", err
	}

	if pr.Status == "MERGED" {
		return nil, "", ErrPRMerged
	}

	isAssigned, err := s.repo.IsReviewerAssigned(req.PullRequestID, req.OldUserID)
	if err != nil {
		return nil, "", err
	}
	if !isAssigned {
		return nil, "", ErrReviewerNotAssigned
	}

	oldReviewer, err := s.repo.GetUser(req.OldUserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", ErrUserNotFound
		}
		return nil, "", err
	}

	candidates, err := s.repo.GetActiveUsersByTeamName(oldReviewer.TeamName, req.OldUserID)
	if err != nil {
		return nil, "", err
	}

	filteredCandidates := []models.User{}
	assignedMap := make(map[string]bool)
	for _, reviewerID := range pr.AssignedReviewers {
		assignedMap[reviewerID] = true
	}

	for _, candidate := range candidates {
		if candidate.UserID != pr.AuthorID && !assignedMap[candidate.UserID] {
			filteredCandidates = append(filteredCandidates, candidate)
		}
	}

	if len(filteredCandidates) == 0 {
		return nil, "", ErrNoCandidate
	}

	newReviewer := filteredCandidates[rand.Intn(len(filteredCandidates))]

	if err := s.repo.ReassignReviewer(req.PullRequestID, req.OldUserID, newReviewer.UserID); err != nil {
		return nil, "", err
	}

	updatedPR, err := s.repo.GetPR(req.PullRequestID)
	if err != nil {
		return nil, "", err
	}

	return updatedPR, newReviewer.UserID, nil
}

func selectReviewers(candidates []models.User, maxCount int) []string {
	if len(candidates) == 0 {
		return []string{}
	}

	count := maxCount
	if len(candidates) < maxCount {
		count = len(candidates)
	}

	selected := make([]string, 0, count)
	indices := rand.Perm(len(candidates))

	for i := 0; i < count; i++ {
		selected = append(selected, candidates[indices[i]].UserID)
	}

	return selected
}
