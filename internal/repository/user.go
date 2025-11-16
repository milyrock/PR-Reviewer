package repository

import (
	"database/sql"

	"github.com/milyrock/PR-Reviewer/internal/models"
)

const (
	selectUser = `
		SELECT user_id, username, team_name, is_active
		FROM users
		WHERE user_id = $1
	`

	updateUserActive = `
		UPDATE users
		SET is_active = $1
		WHERE user_id = $2
	`

	selectActiveUsersByTeam = `
		SELECT user_id, username, team_name, is_active
		FROM users
		WHERE team_name = $1 AND is_active = true AND user_id != $2
		ORDER BY user_id
	`

	selectUserReviewPRs = `
		SELECT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status
		FROM pull_requests pr
		INNER JOIN pr_reviewers prr ON pr.pull_request_id = prr.pull_request_id
		WHERE prr.user_id = $1
		ORDER BY pr.created_at DESC
	`

	selectUserReviewStats = `
		SELECT 
			u.user_id,
			u.username,
			COUNT(prr.pull_request_id) as review_count
		FROM users u
		LEFT JOIN pr_reviewers prr ON u.user_id = prr.user_id
		GROUP BY u.user_id, u.username
		ORDER BY review_count DESC, u.user_id
	`
)

func (r *Repository) GetUser(userID string) (*models.User, error) {
	var user models.User
	err := r.db.Get(&user, selectUser, userID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) SetUserIsActive(userID string, isActive bool) error {
	result, err := r.db.Exec(updateUserActive, isActive, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *Repository) GetActiveUsersByTeamName(teamName string, excludeUserID string) ([]models.User, error) {
	var users []models.User
	err := r.db.Select(&users, selectActiveUsersByTeam, teamName, excludeUserID)
	return users, err
}

func (r *Repository) GetUserReviewPRs(userID string) ([]models.PullRequestShort, error) {
	var prs []models.PullRequestShort
	err := r.db.Select(&prs, selectUserReviewPRs, userID)
	return prs, err
}

func (r *Repository) GetUserReviewStats() ([]models.UserReviewStats, error) {
	var stats []models.UserReviewStats
	err := r.db.Select(&stats, selectUserReviewStats)
	return stats, err
}
