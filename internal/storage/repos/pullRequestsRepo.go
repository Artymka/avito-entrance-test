package repos

import (
	"time"

	"github.com/Artymka/avito-entrance-test/internal/storage/models"
)

type PullRequestsRepo interface {
	Create(pr models.PullRequest) error
	Merge(id string, mergedAt time.Time) error
	Get(id string) (models.PullRequest, error)
}
