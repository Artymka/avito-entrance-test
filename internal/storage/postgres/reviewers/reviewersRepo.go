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
			user_id VARCHAR(255) NOT NULL,
			pull_request_id VARCHAR(255) NOT NULL,
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

func (r *ReviewersRepo) GetCandidates(authorID string, limit int) ([]models.User, error) {
	const op = "postgres.reviewers_repo.get_candidates"

	res := make([]models.User, 0, 2)
	err := r.db.Select(res, `
		SELECT users.id, users.name, users.team_id, users.is_active
		FROM users
		CROSS JOIN (SELECT id, team_id FROM users WHERE id = $1) AS auth
		WHERE users.team_id = auth.team_id
			AND users.is_active = true
			AND users.id <> auth.id
		LIMIT $2
	`, authorID, limit)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return res, nil
}
