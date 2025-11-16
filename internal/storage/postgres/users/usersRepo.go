package users

import (
	"fmt"

	"github.com/Artymka/avito-entrance-test/internal/storage/models"
	"github.com/jmoiron/sqlx"
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
			team_id VARCHAR(255) NOT NULL,
			is_active BOOLEAN NOT NULL DEFAULT false
		)
	`)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *UsersRepo) Create(user *models.User) error {
	const op = "postgres.users_repo.create"
	err := r.db.Get(&user.ID, `
		INSERT INTO users
		(id, name, team_id, is_active)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, user.ID, user.Name, user.TeamID, user.IsActive)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
