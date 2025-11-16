package getreview

import (
	"net/http"

	"github.com/Artymka/avito-entrance-test/internal/lib/logging"
	"github.com/Artymka/avito-entrance-test/internal/lib/response"
	"github.com/Artymka/avito-entrance-test/internal/schemas"
	"github.com/Artymka/avito-entrance-test/internal/storage/repos"
	"github.com/go-chi/render"
)

// Получить PR'ы, где пользователь назначен ревьювером
func New(u repos.UsersRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.users.get_review"

		// query check
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Err(response.CodeBadRequest, ""))
			return
		}

		// // check user (not described in openapi.yml)
		// exists, err := u.Exists(userID)
		// if err != nil {
		// 	logging.Err(op, err)
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	render.JSON(w, r, response.Err(response.CodeServer, ""))
		// 	return
		// }
		// if !exists {
		// 	w.WriteHeader(http.StatusNotFound)
		// 	render.JSON(w, r, response.Err(response.CodeNotFound, ""))
		// 	return
		// }

		// get review
		pullRequests, err := u.GetReview(userID)
		if err != nil {
			logging.Err(op, err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Err(response.CodeServer, ""))
			return
		}

		respPRS := make([]schemas.PRShort, 0, len(pullRequests))
		for _, pr := range pullRequests {
			respPRS = append(respPRS, schemas.PRShort{
				ID:       pr.ID,
				Name:     pr.Name,
				AuthorID: pr.AuthorID,
				Status:   pr.Status,
			})
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, schemas.UserReviewResponse{
			UserID: userID,
			PRs:    respPRS,
		})
	}
}
