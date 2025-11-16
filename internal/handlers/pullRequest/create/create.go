package create

import (
	"errors"
	"net/http"

	"github.com/Artymka/avito-entrance-test/internal/lib/logging"
	"github.com/Artymka/avito-entrance-test/internal/lib/response"
	"github.com/Artymka/avito-entrance-test/internal/schemas"
	"github.com/Artymka/avito-entrance-test/internal/storage"
	"github.com/Artymka/avito-entrance-test/internal/storage/models"
	"github.com/Artymka/avito-entrance-test/internal/storage/repos"
	"github.com/go-chi/render"
)

// Создать PR и автоматически назначить до 2 ревьюверов из команды автора
func New(pr repos.PullRequestsRepo, re repos.ReviewersRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pull_request.create"
		// json check
		var req schemas.PRCreateRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			logging.Err(op, err)
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Err(response.CodeBadRequest, ""))
			return
		}

		// create pull request
		err := pr.Create(models.PullRequest{
			ID:       req.ID,
			Name:     req.Name,
			AuthorID: req.AuthorID,
		})
		if err != nil {
			if errors.Is(err, storage.ErrUnique) {
				logging.Err(op, err)
				w.WriteHeader(http.StatusConflict)
				render.JSON(w, r, response.Err(response.CodePRExists, "PR id already exists"))
				return
			}
			if errors.Is(err, storage.ErrForeignKey) {
				logging.Err(op, err)
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, response.Err(response.CodeNotFound, ""))
				return
			}
			logging.Err(op, err)
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Err(response.CodeBadRequest, ""))
			return
		}

		// assign reviewers
		candidates, err := re.GetCandidates(req.AuthorID, 2)
		if err != nil {
			logging.Err(op, err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Err(response.CodeBadRequest, ""))
			return
		}

		res := make([]string, 0, 2)
		for _, candidate := range candidates {
			if candidate.ID != "" {
				err = re.Create(models.Reviewer{
					UserID:        candidate.ID,
					PullRequestID: req.ID,
				})

				if err != nil {
					logging.Err(op, err)
					w.WriteHeader(http.StatusInternalServerError)
					render.JSON(w, r, response.Err(response.CodeBadRequest, ""))
					return
				}

				res = append(res, candidate.ID)
			}
		}

		// success
		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, schemas.PRCreateResponse{
			PR: schemas.PRResponse{
				ID:               req.ID,
				Name:             req.Name,
				AuthorID:         req.AuthorID,
				Status:           "OPEN",
				AssignedReviwers: res,
			},
		})
	}
}
