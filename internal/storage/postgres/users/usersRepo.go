package users

import (
	"errors"
	"fmt"

	"github.com/Artymka/avito-entrance-test/internal/storage"
	"github.com/Artymka/avito-entrance-test/internal/storage/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type UsersRepo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) (*UsersRepo, error) {
	const op = "postgres.users_repo.new"
	repo := UsersRepo{db: db}
	err := repo.createTable()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &repo, nil
}

func (r *UsersRepo) createTable() error {
	const op = "postgres.users_repo.create_table"
	_, err := r.db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id VARCHAR(255) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			team_name VARCHAR(255) NOT NULL,
			is_active BOOLEAN NOT NULL DEFAULT false
		)
	`)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *UsersRepo) Create(user models.User) error {
	const op = "postgres.users_repo.create"
	_, err := r.db.Exec(`
		INSERT INTO users
		(id, name, team_name, is_active)
		VALUES ($1, $2, $3, $4)
	`, user.ID, user.Name, user.TeamName, user.IsActive)

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

func (r *UsersRepo) Update(user models.User) error {
	const op = "postgres.users_repo.update"
	res, err := r.db.Exec(`
		UPDATE users
		SET name = $1,
			team_name = $2,
			is_active = $3
		WHERE id = $4
	`, user.Name, user.TeamName, user.IsActive, user.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rows != 1 {
		return fmt.Errorf("%s: %w", op, storage.ErrWrongUpadtes)
	}

	return nil
}
