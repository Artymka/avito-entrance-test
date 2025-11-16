package teams

import (
	"errors"
	"fmt"

	"github.com/Artymka/avito-entrance-test/internal/storage"
	"github.com/Artymka/avito-entrance-test/internal/storage/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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
			name VARCHAR(255) PRIMARY KEY
		)
	`)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *TeamsRepo) Create(team models.Team) error {
	const op = "postgres.teams_repo.create"
	_, err := r.db.Exec(`
		INSERT INTO teams
		(name)
		VALUES ($1)
	`, team.Name)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23505": // unique_violation
				return fmt.Errorf("%s: %w", op, storage.ErrUnique)
			default:
				return fmt.Errorf("%s: %w", op, err)
			}
		}
	}
	return nil
}
