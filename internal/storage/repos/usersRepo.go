package repos

import "github.com/Artymka/avito-entrance-test/internal/storage/models"

type UsersRepo interface {
	Create(user models.User) error
	Update(user models.User) error
}
