package main

import (
	"net/http"

	"github.com/Artymka/avito-entrance-test/internal/config"
	pullRequestCreate "github.com/Artymka/avito-entrance-test/internal/handlers/pullRequest/create"
	pullRequestMerge "github.com/Artymka/avito-entrance-test/internal/handlers/pullRequest/merge"
	pullRequestReassign "github.com/Artymka/avito-entrance-test/internal/handlers/pullRequest/reassign"
	teamAdd "github.com/Artymka/avito-entrance-test/internal/handlers/team/add"
	teamGet "github.com/Artymka/avito-entrance-test/internal/handlers/team/get"
	userSetIsActive "github.com/Artymka/avito-entrance-test/internal/handlers/users/setIsActive"
	"github.com/Artymka/avito-entrance-test/internal/lib/logging"
	"github.com/Artymka/avito-entrance-test/internal/storage/postgres"
	pullrequests "github.com/Artymka/avito-entrance-test/internal/storage/postgres/pullRequests"
	"github.com/Artymka/avito-entrance-test/internal/storage/postgres/reviewers"
	"github.com/Artymka/avito-entrance-test/internal/storage/postgres/teams"
	"github.com/Artymka/avito-entrance-test/internal/storage/postgres/users"
	"github.com/go-chi/chi"
)

func main() {
	// config
	cfg, err := config.New("./config/local.yaml")
	if err != nil {
		panic(err)
	}

	// postgres
	p, err := postgres.New(cfg)
	if err != nil {
		panic(err)
	}
	usersRepo, err := users.New(p.DB)
	if err != nil {
		panic(err)
	}
	teamsRepo, err := teams.New(p.DB)
	if err != nil {
		panic(err)
	}
	prRepo, err := pullrequests.New(p.DB)
	if err != nil {
		panic(err)
	}
	reviewersRepo, err := reviewers.New(p.DB)
	if err != nil {
		panic(err)
	}

	// router
	r := chi.NewRouter()
	r.Post("/team/add", teamAdd.New(teamsRepo, usersRepo))
	r.Get("/team/get", teamGet.New(teamsRepo))
	r.Post("/users/setIsActive", userSetIsActive.New(usersRepo))
	r.Post("/pullRequest/create", pullRequestCreate.New(prRepo, reviewersRepo))
	r.Post("/pullRequest/merge", pullRequestMerge.New(prRepo, reviewersRepo))
	r.Post("/pullRequest/reassign", pullRequestReassign.New(prRepo, reviewersRepo, usersRepo))

	// listen
	logging.Info("main", "server is listenning...")
	err = http.ListenAndServe(cfg.Server.Address, r)
	if err != nil {
		logging.Err("server down", err)
	}
}
