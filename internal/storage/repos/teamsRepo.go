package repos

import "github.com/Artymka/avito-entrance-test/internal/storage/models"

type TeamsRepo interface {
	Create(team models.Team) error
	Exists(team models.Team) (bool, error)
	GetMembers(team models.Team) ([]models.User, error)
}
