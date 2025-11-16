package add

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

// Создать команду с участниками (создаёт/обновляет пользователей)
func New(t repos.TeamsRepo, u repos.UsersRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.team.add"
		// json check
		var req schemas.TeamAddRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			logging.Err(op, err)
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Err(response.CodeBadRequest, ""))
			return
		}

		// create team
		err := t.Create(models.Team{
			Name: req.Name,
		})
		if err != nil {
			if errors.Is(err, storage.ErrUnique) {
				logging.Err(op, err)
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Err(response.CodeTeamExists, "team_name already exists"))
				return
			}
			logging.Err(op, err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Err(response.CodeServer, ""))
			return
		}

		// create & update users
		for _, member := range req.Members {
			user := models.User{
				ID:       member.ID,
				Name:     member.Name,
				IsActive: member.IsActive,
				TeamName: req.Name,
			}

			if err = u.Create(user); err != nil {
				if errors.Is(err, storage.ErrUnique) {
					err = u.Update(user)
					if err == nil {
						continue
					}
				}

				logging.Err(op, err)
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Err(response.CodeBadRequest, ""))
				return
			}
		}

		// success
		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, req)
	}
}
