package repos

import "github.com/Artymka/avito-entrance-test/internal/storage/models"

type UsersRepo interface {
	Create(user models.User) error
	Update(user models.User) error
	SetIsActive(user models.User) error
	Exists(userID string) (bool, error)
	Get(userID string) (models.User, error)
}
