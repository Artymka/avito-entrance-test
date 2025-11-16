package setisactive

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

// Установить флаг активности пользователя
func New(u repos.UsersRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.users.set_is_active"

		// json check
		var req schemas.UserSetIsActiveRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			logging.Err(op, err)
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Err(response.CodeBadRequest, ""))
			return
		}

		// check user
		user, err := u.Get(req.ID)

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

		// update user
		err = u.SetIsActive(models.User{
			ID:       req.ID,
			IsActive: req.IsActive,
		})

		if err != nil {
			logging.Err(op, err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Err(response.CodeServer, ""))
			return
		}

		// success
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, schemas.UserRsponse{
			User: schemas.User{
				ID:       req.ID,
				Name:     user.Name,
				TeamName: user.TeamName,
				IsActive: req.IsActive,
			},
		})
	}
}
