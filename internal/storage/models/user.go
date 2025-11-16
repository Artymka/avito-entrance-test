package models

type User struct {
	ID       int64  `db:"id"`
	Name     string `db:"name"`
	IsActive bool   `db:"is_active"`
	TeamID   int64  `db:"team_id"`
}
