package models

type User struct {
	ID       string `db:"id"`
	Name     string `db:"name"`
	IsActive bool   `db:"is_active"`
	TeamID   string `db:"team_id"`
}
