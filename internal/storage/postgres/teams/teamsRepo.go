package teams

import (
	"fmt"

	"github.com/Artymka/avito-entrance-test/internal/storage/models"
	"github.com/jmoiron/sqlx"
)

type TeamsRepo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) (*TeamsRepo, error) {
	const op = "postgres.teams_repo.new"
	repo := TeamsRepo{db: db}
	err := repo.createTable()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &repo, nil
}

func (r *TeamsRepo) createTable() error {
	const op = "postgres.teams_repo.create_table"
	_, err := r.db.Exec(`
		CREATE TABLE IF NOT EXISTS teams (
			id VARCHAR(255) PRIMARY KEY,
			name VARCHAR(255) NOT NULL
		)
	`)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *TeamsRepo) Create(team *models.Team) error {
	const op = "postgres.teams_repo.create"
	err := r.db.Get(&team.ID, `
		INSERT INTO teams
		(name)
		VALUES ($1)
		RETURNING id
	`, team.Name)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
