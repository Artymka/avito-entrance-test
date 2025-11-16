package models

type PullRequest struct {
	ID       string `db:"id"`
	Name     string `db:"name"`
	AuthorID string `db:"author_id"`
	Status   string `db:"status"`
}

const (
	PRStatusOpen   = "OPEN"
	PRStatusMerged = "MERGED"
)
