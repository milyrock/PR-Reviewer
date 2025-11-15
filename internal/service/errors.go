package service

import "errors"

var (
	ErrPRExists            = errors.New("PR id already exists")
	ErrPRNotFound          = errors.New("resource not found")
	ErrPRMerged            = errors.New("cannot reassign on merged PR")
	ErrReviewerNotAssigned = errors.New("reviewer is not assigned to this PR")
	ErrNoCandidate         = errors.New("no active replacement candidate in team")
	ErrUserNotFound        = errors.New("resource not found")
	ErrTeamExists          = errors.New("team_name already exists")
	ErrTeamNotFound        = errors.New("resource not found")
)
