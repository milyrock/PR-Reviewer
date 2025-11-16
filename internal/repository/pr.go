package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/milyrock/PR-Reviewer/internal/models"
)

const (
	prExists = `SELECT EXISTS(SELECT 1 FROM pull_requests WHERE pull_request_id = $1)`

	selectPR = `
		SELECT pull_request_id, pull_request_name, author_id, status, created_at, merged_at
		FROM pull_requests
		WHERE pull_request_id = $1
	`

	selectPRReviewers = `
		SELECT user_id
		FROM pr_reviewers
		WHERE pull_request_id = $1
		ORDER BY user_id
	`

	insertPR = `
		INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	insertPRReviewer = `
		INSERT INTO pr_reviewers (pull_request_id, user_id)
		VALUES ($1, $2)
	`

	mergePR = `
		UPDATE pull_requests
		SET status = 'MERGED', merged_at = $1
		WHERE pull_request_id = $2
	`

	reviewerAssigned = `
		SELECT EXISTS(
			SELECT 1 FROM pr_reviewers
			WHERE pull_request_id = $1 AND user_id = $2
		)
	`

	deletePRReviewer = `
		DELETE FROM pr_reviewers
		WHERE pull_request_id = $1 AND user_id = $2
	`

	selectPRReviewStats = `
		SELECT 
			pr.pull_request_id,
			pr.pull_request_name,
			COUNT(prr.user_id) as reviewer_count
		FROM pull_requests pr
		LEFT JOIN pr_reviewers prr ON pr.pull_request_id = prr.pull_request_id
		GROUP BY pr.pull_request_id, pr.pull_request_name
		ORDER BY pr.created_at DESC
	`
)

func (r *Repository) PRExists(pullRequestID string) (bool, error) {
	var exists bool
	err := r.db.Get(&exists, prExists, pullRequestID)
	return exists, err
}

func (r *Repository) IsReviewerAssigned(pullRequestID, userID string) (bool, error) {
	var exists bool
	err := r.db.Get(&exists, reviewerAssigned, pullRequestID, userID)
	return exists, err
}

func (r *Repository) GetPR(pullRequestID string) (*models.PullRequest, error) {
	var pr models.PullRequest
	err := r.db.Get(&pr, selectPR, pullRequestID)
	if err != nil {
		return nil, err
	}

	err = r.db.Select(&pr.AssignedReviewers, selectPRReviewers, pullRequestID)
	if err != nil {
		return nil, err
	}

	return &pr, nil
}

func (r *Repository) CreatePR(pr *models.PullRequest, reviewerIDs []string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	now := time.Now()
	_, err = tx.Exec(insertPR, pr.PullRequestID, pr.PullRequestName, pr.AuthorID, pr.Status, now)
	if err != nil {
		return err
	}

	for _, reviewerID := range reviewerIDs {
		_, err = tx.Exec(insertPRReviewer, pr.PullRequestID, reviewerID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *Repository) MergePR(pullRequestID string) error {
	pr, err := r.GetPR(pullRequestID)
	if err != nil {
		return err
	}

	if pr.Status == "MERGED" {
		return nil
	}

	now := time.Now()
	result, err := r.db.Exec(mergePR, now, pullRequestID)
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

func (r *Repository) ReassignReviewer(pullRequestID, oldUserID, newUserID string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	result, err := tx.Exec(deletePRReviewer, pullRequestID, oldUserID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("reviewer not assigned")
	}

	_, err = tx.Exec(insertPRReviewer, pullRequestID, newUserID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *Repository) GetPRReviewStats() ([]models.PRReviewStats, error) {
	var stats []models.PRReviewStats
	err := r.db.Select(&stats, selectPRReviewStats)
	return stats, err
}
