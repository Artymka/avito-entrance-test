package reviewers

import (
	"fmt"

	"github.com/Artymka/avito-entrance-test/internal/storage/models"
	"github.com/jmoiron/sqlx"
)

type ReviewersRepo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) (*ReviewersRepo, error) {
	const op = "postgres.reviewers_repo.new"
	repo := ReviewersRepo{db: db}
	err := repo.createTable()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &repo, nil
}

func (r *ReviewersRepo) createTable() error {
	const op = "postgres.reviewers_repo.create_table"
	_, err := r.db.Exec(`
		CREATE TABLE IF NOT EXISTS reviewers (
			user_id INTEGER NOT NULL,
			pull_request_id INTEGER NOT NULL,
			PRIMARY KEY (user_id, pull_request_id)
		)
	`)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *ReviewersRepo) Create(reviewer models.Reviewer) error {
	const op = "postgres.reviewers_repo.create"
	_, err := r.db.Exec(`
		INSERT INTO reviewers
		(user_id, pull_request_id)
		VALUES ($1, $2)
	`, reviewer.UserID, reviewer.PullRequestID)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
