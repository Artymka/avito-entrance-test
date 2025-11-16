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

func (r *ReviewersRepo) GetCandidates(authorID string, PRID string, limit int) ([]models.User, error) {
	const op = "postgres.reviewers_repo.get_candidates"

	res := make([]models.User, 0, 2)
	err := r.db.Select(&res, `
		SELECT users.id, users.name, users.team_name, users.is_active
		FROM users
		CROSS JOIN (SELECT id, team_name FROM users WHERE id = $1) AS author
		WHERE users.team_name = author.team_name
			AND users.is_active = true
			AND users.id <> author.id
			AND NOT EXISTS (
				SELECT 1 FROM reviewers 
				WHERE reviewers.user_id = users.id 
				AND reviewers.pull_request_id = $2
			)
		LIMIT $3
	`, authorID, PRID, limit)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return res, nil
}

func (r *ReviewersRepo) Get(pullRequestID string) ([]models.User, error) {
	const op = "postgres.reviewers_repo.get"

	res := make([]models.User, 0, 2)
	err := r.db.Select(&res, `
		SELECT users.id, users.name, users.team_name, users.is_active
		FROM users
		JOIN reviewers
		ON users.id = reviewers.user_id
		WHERE reviewers.pull_request_id = $1
	`, pullRequestID)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return res, nil
}

func (r *ReviewersRepo) Has(reviewer models.Reviewer) (bool, error) {
	const op = "postgres.reviewers_repo.has"

	var has bool
	err := r.db.Get(&has, `
		SELECT EXISTS (
			SELECT 1 FROM reviewers 
			WHERE user_id = $1 AND pull_request_id = $2
		)
	`, reviewer.UserID, reviewer.PullRequestID)

	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return has, nil
}

func (r *ReviewersRepo) Delete(reviewer models.Reviewer) error {
	const op = "postgres.reviewers_repo.delete"

	_, err := r.db.Exec(`
		DELETE FROM reviewers
		WHERE user_id = $1 AND pull_request_id = $2
	`, reviewer.UserID, reviewer.PullRequestID)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
