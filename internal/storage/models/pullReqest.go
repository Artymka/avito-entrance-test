package models

import "time"

type PullRequest struct {
	ID       string     `db:"id"`
	Name     string     `db:"name"`
	AuthorID string     `db:"author_id"`
	Status   string     `db:"status"`
	MergedAt *time.Time `db:"merged_at"`
}

const (
	PRStatusOpen   = "OPEN"
	PRStatusMerged = "MERGED"
)
