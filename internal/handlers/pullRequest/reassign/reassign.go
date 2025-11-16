package reassign

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

// Переназначить конкретного ревьювера на другого из его команды
func New(pr repos.PullRequestsRepo, re repos.ReviewersRepo, u repos.UsersRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pull_request.reassign"

		// json check
		var req schemas.PRReassignRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			logging.Err(op, err)
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Err(response.CodeBadRequest, ""))
			return
		}

		// pr and user check
		pullRequest, err := pr.Get(req.PRID)
		if err != nil {
			if errors.Is(err, storage.ErrNoRows) {
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, response.Err(response.CodeNotFound, ""))
				return
			}
			logging.Err(op, err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Err(response.CodeServer, ""))
			return
		}

		exists, err := u.Exists(req.UserID)
		if !exists {
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, response.Err(response.CodeNotFound, ""))
			return
		}
		if err != nil {
			logging.Err(op, err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Err(response.CodeServer, ""))
			return
		}

		// check pr status
		if pullRequest.Status == "MERGED" {
			w.WriteHeader(http.StatusConflict)
			render.JSON(w, r, response.Err(response.CodePRMerged, "cannot reassign on merged PR"))
			return
		}

		// check reviewer
		oldReviewer := models.Reviewer{
			PullRequestID: req.PRID,
			UserID:        req.UserID,
		}
		has, err := re.Has(oldReviewer)
		if err != nil {
			logging.Err(op, err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Err(response.CodeServer, ""))
			return
		}
		if !has {
			w.WriteHeader(http.StatusConflict)
			render.JSON(w, r, response.Err(response.CodeNotAssigned, "reviewer is not assigned to this PR"))
			return
		}

		// check candidates
		candidates, err := re.GetCandidates(pullRequest.AuthorID, req.PRID, 1)
		if err != nil {
			logging.Err(op, err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Err(response.CodeServer, ""))
			return
		}
		if len(candidates) == 0 {
			w.WriteHeader(http.StatusConflict)
			render.JSON(w, r, response.Err(response.CodeNoCandidate, ""))
			return
		}

		// success
		err = re.Create(models.Reviewer{
			UserID:        candidates[0].ID,
			PullRequestID: req.PRID,
		})
		if err != nil {
			logging.Err(op, err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Err(response.CodeServer, ""))
			return
		}

		err = re.Delete(oldReviewer)
		if err != nil {
			logging.Err(op, err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Err(response.CodeServer, ""))
			return
		}

		newUsers, err := re.Get(req.PRID)
		if err != nil {
			logging.Err(op, err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Err(response.CodeServer, ""))
			return
		}

		newCandidates := make([]string, 0, len(newUsers))
		for _, user := range newUsers {
			newCandidates = append(newCandidates, user.ID)
		}

		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, schemas.PRReassignResponse{
			PR: schemas.PRResponse{
				ID:               pullRequest.ID,
				Name:             pullRequest.Name,
				AuthorID:         pullRequest.AuthorID,
				Status:           pullRequest.Status,
				AssignedReviwers: newCandidates,
			},
			ReplacedBy: candidates[0].ID,
		})

		// pr:
		// 	pull_request_id: pr-1001
		// 	pull_request_name: Add search
		// 	author_id: u1
		// 	status: OPEN
		// 	assigned_reviewers: [u3, u5]
		// replaced_by: u5
	}
}
