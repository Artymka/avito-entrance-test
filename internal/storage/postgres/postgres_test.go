package postgres

import (
	"testing"

	"github.com/Artymka/avito-entrance-test/internal/config"
	"github.com/Artymka/avito-entrance-test/internal/storage/models"
	pullrequests "github.com/Artymka/avito-entrance-test/internal/storage/postgres/pullRequests"
	"github.com/Artymka/avito-entrance-test/internal/storage/postgres/reviewers"
	"github.com/Artymka/avito-entrance-test/internal/storage/postgres/teams"
	"github.com/Artymka/avito-entrance-test/internal/storage/postgres/users"
	"github.com/stretchr/testify/require"
)

var cfg *config.Config
var p *Postgres

var usersRepo *users.UsersRepo
var teamsRepo *teams.TeamsRepo
var prRepo *pullrequests.PullRequestsRepo
var reviewersRepo *reviewers.ReviewersRepo
var err error

func TestTables(t *testing.T) {
	t.Run("Connecting to db test", func(t *testing.T) {
		cfg, err = config.New("./../../../config/local.yaml")
		require.Nil(t, err)

		p, err = New(cfg)
		require.Nil(t, err)
	})

	t.Run("Creating tables test", func(t *testing.T) {
		usersRepo, err = users.New(p.DB)
		require.Nil(t, err)

		teamsRepo, err = teams.New(p.DB)
		require.Nil(t, err)

		prRepo, err = pullrequests.New(p.DB)
		require.Nil(t, err)

		reviewersRepo, err = reviewers.New(p.DB)
		require.Nil(t, err)
	})

	t.Run("Creating records test", func(t *testing.T) {
		user := models.User{
			ID:       "123",
			Name:     "john",
			IsActive: true,
			TeamID:   "123",
		}
		err = usersRepo.Create(&user)
		require.Nil(t, err)
		t.Logf("user id: %s", user.ID)

		pr := models.PullRequest{
			ID:       "321",
			Name:     "smth",
			AuthorID: "123",
			Status:   "OPEN",
		}
		err = prRepo.Create(pr)
		require.Nil(t, err)
		t.Logf("pr id: %s", pr.ID)
	})
}
