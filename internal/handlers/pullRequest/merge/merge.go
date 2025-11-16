package merge

import (
	"errors"
	"net/http"
	"time"

	"github.com/Artymka/avito-entrance-test/internal/lib/logging"
	"github.com/Artymka/avito-entrance-test/internal/lib/response"
	"github.com/Artymka/avito-entrance-test/internal/schemas"
	"github.com/Artymka/avito-entrance-test/internal/storage"
	"github.com/Artymka/avito-entrance-test/internal/storage/repos"
	"github.com/go-chi/render"
)

// Пометить PR как MERGED (идемпотентная операция)
func New(pr repos.PullRequestsRepo, re repos.ReviewersRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pull_request.merge"

		// json check
		var req schemas.PRMergeRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			logging.Err(op, err)
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Err(response.CodeBadRequest, ""))
			return
		}

		// check pull request
		pullRequest, err := pr.Get(req.ID)
		if err != nil {
			if errors.Is(err, storage.ErrNoRows) {
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, response.Err(response.CodeNotFound, ""))
				return
			}
		}

		// merge
		if pullRequest.Status == "OPEN" {
			now := time.Now()
			pullRequest.MergedAt = &now
			if err = pr.Merge(pullRequest.ID, now); err != nil {
				logging.Err(op, err)
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Err(response.CodeServer, ""))
				return
			}
		}

		// send
		users, err := re.Get(req.ID)
		if err != nil {
			logging.Err(op, err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Err(response.CodeServer, ""))
			return
		}

		reviewers := make([]string, 0, 2)
		for _, user := range users {
			reviewers = append(reviewers, user.ID)
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, schemas.PRMergeResponse{
			PR: schemas.PRMerge{
				ID:               pullRequest.ID,
				Name:             pullRequest.Name,
				AuthorID:         pullRequest.AuthorID,
				Status:           "MERGED",
				AssignedReviwers: reviewers,
				MergedAt:         pullRequest.MergedAt.Format("2006-01-02T15:04:05Z07:00"),
			},
		})
	}
}
