package models

type Reviewer struct {
	UserID        string `db:"user_id"`
	PullRequestID string `db:"pull_request_id"`
}
