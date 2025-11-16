package models

type PullRequest struct {
	ID       int64  `db:"id"`
	Name     string `db:"name"`
	AuthorID int64  `db:"author_id"`
	Status   string `db:"status"`
}

const (
	PRStatusOpen   = "OPEN"
	PRStatusMerged = "MERGED"
)
