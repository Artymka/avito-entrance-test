package repos

import "github.com/Artymka/avito-entrance-test/internal/storage/models"

type ReviewersRepo interface {
	Create(reviewer models.Reviewer) error
	GetCandidates(authorID string, limit int) ([]models.User, error)
	Get(pullRequestID string) ([]models.User, error)
}
