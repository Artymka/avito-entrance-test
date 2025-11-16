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

func (r *TeamsRepo) Exists(team models.Team) (bool, error) {
	const op = "postgres.teams_repo.exists"

	var res bool
	err := r.db.Get(&res, `
		SELECT EXISTS (SELECT 1 FROM teams WHERE name = $1)
	`, team.Name)

	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return res, nil
}

func (r *TeamsRepo) GetMembers(team models.Team) ([]models.User, error) {
	const op = "postgres.teams_repo.get_members"

	res := make([]models.User, 0)
	err := r.db.Select(&res, `
		SELECT id, name, is_active
		FROM users
		WHERE team_name = $1
	`, team.Name)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return res, nil
}
