package pullrequests

import (
	"fmt"

	"github.com/Artymka/avito-entrance-test/internal/storage/models"
	"github.com/jmoiron/sqlx"
)

type PullRequestsRepo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) (*PullRequestsRepo, error) {
	const op = "postgres.pull_requests_repo.new"
	repo := PullRequestsRepo{db: db}
	err := repo.createTable()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &repo, nil
}

func (r *PullRequestsRepo) createTable() error {
	const op = "postgres.pull_requests_repo.create_table"

	_, err := r.db.Exec(`
	 	DO $$ 
        BEGIN
            IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'pull_request_status') THEN
                CREATE TYPE pull_request_status AS ENUM('OPEN', 'MERGED');
            END IF;
        END
        $$;

		CREATE TABLE IF NOT EXISTS pull_requests (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			author_id INTEGER NOT NULL,
			status pull_request_status NOT NULL DEFAULT 'OPEN'
		)
	`)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *PullRequestsRepo) Create(pr *models.PullRequest) error {
	const op = "postgres.pull_requests_repo.create"
	err := r.db.Get(&pr.ID, `
		INSERT INTO pull_requests
		(name, author_id, status)
		VALUES ($1, $2, $3)
		RETURNING id
	`, pr.Name, pr.AuthorID, pr.Status)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
