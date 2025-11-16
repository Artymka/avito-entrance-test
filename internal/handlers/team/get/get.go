package get

import (
	"net/http"

	"github.com/Artymka/avito-entrance-test/internal/lib/logging"
	"github.com/Artymka/avito-entrance-test/internal/lib/response"
	"github.com/Artymka/avito-entrance-test/internal/schemas"
	"github.com/Artymka/avito-entrance-test/internal/storage/models"
	"github.com/Artymka/avito-entrance-test/internal/storage/repos"
	"github.com/go-chi/render"
)

// Получить команду с участниками
func New(t repos.TeamsRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.team.get"
		// query check
		teamName := r.URL.Query().Get("team_name")
		if teamName == "" {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Err(response.CodeBadRequest, ""))
			return
		}
		team := models.Team{Name: teamName}

		// check team
		exists, err := t.Exists(team)
		if err != nil {
			logging.Err(op, err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Err(response.CodeServer, ""))
			return
		}

		if !exists {
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, response.Err(response.CodeNotFound, ""))
			return
		}

		// get users
		members, err := t.GetMembers(team)
		if err != nil {
			logging.Err(op, err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Err(response.CodeServer, ""))
			return
		}

		// success
		resMembers := make([]schemas.Member, 0, len(members))
		for _, member := range members {
			resMembers = append(resMembers, schemas.Member{
				ID:       member.ID,
				Name:     member.Name,
				IsActive: member.IsActive,
			})
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, schemas.TeamAddRequest{
			Name:    teamName,
			Members: resMembers,
		})
	}
}
