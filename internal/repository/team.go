package repository

import (
	"database/sql"

	"github.com/milyrock/PR-Reviewer/internal/models"
)

const (
	insertTeam = `INSERT INTO teams (team_name) VALUES ($1) ON CONFLICT (team_name) DO NOTHING`

	teamExists = `SELECT EXISTS(SELECT 1 FROM teams WHERE team_name = $1)`

	selectTeamMembers = `
		SELECT user_id, username, is_active
		FROM users
		WHERE team_name = $1
		ORDER BY user_id
	`

	insertOrUpdateUser = `
		INSERT INTO users (user_id, username, team_name, is_active)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) DO UPDATE
		SET username = $2, team_name = $3, is_active = $4
	`
)

func (r *Repository) CreateTeam(teamName string, members []models.TeamMember) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	defer tx.Rollback() //nolint:errcheck

	_, err = tx.Exec(insertTeam, teamName)
	if err != nil {
		return err
	}

	for _, member := range members {
		_, err = tx.Exec(insertOrUpdateUser, member.UserID, member.Username, teamName, member.IsActive)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *Repository) TeamExists(teamName string) (bool, error) {
	var exists bool
	err := r.db.Get(&exists, teamExists, teamName)
	return exists, err
}

func (r *Repository) GetTeam(teamName string) (*models.Team, error) {
	exists, err := r.TeamExists(teamName)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, sql.ErrNoRows
	}

	var team models.Team
	team.TeamName = teamName

	err = r.db.Select(&team.Members, selectTeamMembers, teamName)
	if err != nil {
		return nil, err
	}

	return &team, nil
}
