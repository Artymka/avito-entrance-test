package models

type Reviewer struct {
	UserID        int64 `db:"user_id"`
	PullRequestID int64 `db:"pull_request_id"`
}
