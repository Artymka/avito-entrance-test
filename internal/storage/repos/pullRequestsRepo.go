package repos

import "github.com/Artymka/avito-entrance-test/internal/storage/models"

type PullRequestsRepo interface {
	Create(pr models.PullRequest) error
}
