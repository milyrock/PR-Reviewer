package models

import "time"

type TeamMember struct {
	UserID   string `json:"user_id" db:"user_id"`
	Username string `json:"username" validate:"required" db:"username"`
	IsActive bool   `json:"is_active" db:"is_active"`
}

type Team struct {
	TeamName string       `json:"team_name" db:"team_name"`
	Members  []TeamMember `json:"members"`
}

type User struct {
	UserID   string `json:"user_id" db:"user_id"`
	Username string `json:"username" validate:"required" db:"username"`
	TeamName string `json:"team_name" db:"team_name"`
	IsActive bool   `json:"is_active" db:"is_active"`
}

type PullRequest struct {
	PullRequestID     string     `json:"pull_request_id" db:"pull_request_id"`
	PullRequestName   string     `json:"pull_request_name" db:"pull_request_name"`
	AuthorID          string     `json:"author_id" db:"author_id"`
	Status            string     `json:"status" db:"status"`
	AssignedReviewers []string   `json:"assigned_reviewers"`
	CreatedAt         *time.Time `json:"createdAt,omitempty" db:"created_at"`
	MergedAt          *time.Time `json:"mergedAt,omitempty" db:"merged_at"`
}

type PullRequestShort struct {
	PullRequestID   string `json:"pull_request_id" db:"pull_request_id"`
	PullRequestName string `json:"pull_request_name" db:"pull_request_name"`
	AuthorID        string `json:"author_id" db:"author_id"`
	Status          string `json:"status" db:"status"`
}

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type CreateTeamRequest struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}

type SetIsActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type CreatePRRequest struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

type MergePRRequest struct {
	PullRequestID string `json:"pull_request_id"`
}

type ReassignPRRequest struct {
	PullRequestID string `json:"pull_request_id"`
	OldUserID     string `json:"old_reviewer_id"`
}

type UserReviewStats struct {
	UserID      string `json:"user_id" db:"user_id"`
	Username    string `json:"username" db:"username"`
	ReviewCount int    `json:"review_count" db:"review_count"`
}

type PRReviewStats struct {
	PullRequestID   string `json:"pull_request_id" db:"pull_request_id"`
	PullRequestName string `json:"pull_request_name" db:"pull_request_name"`
	ReviewerCount   int    `json:"reviewer_count" db:"reviewer_count"`
}

type StatisticsResponse struct {
	UserStats []UserReviewStats `json:"user_stats"`
	PRStats   []PRReviewStats   `json:"pr_stats"`
}
