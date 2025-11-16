package pullrequests

import (
	"errors"
	"fmt"

	"github.com/Artymka/avito-entrance-test/internal/storage"
	"github.com/Artymka/avito-entrance-test/internal/storage/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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
			id VARCHAR(255) PRIMARY KEY NOT NULL,
			name VARCHAR(255) NOT NULL,
			author_id VARCHAR(255) NOT NULL,
			status pull_request_status NOT NULL DEFAULT 'OPEN',
			FOREIGN KEY (authour_id) REFERENCES users(id)
		)
	`)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *PullRequestsRepo) Create(pr models.PullRequest) error {
	const op = "postgres.pull_requests_repo.create"
	_, err := r.db.Exec(`
		INSERT INTO pull_requests
		(id, name, author_id)
		VALUES ($1, $2, $3)
	`, pr.ID, pr.Name, pr.AuthorID)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23505": // unique_violation
				return fmt.Errorf("%s: %w", op, storage.ErrUnique)
			case "23503": // foreign_key_violation
				return fmt.Errorf("referenced entity does not exist")
			default:
				return fmt.Errorf("%s: %w", op, err)
			}
		}
	}

	return nil
}
