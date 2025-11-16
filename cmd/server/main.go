package main

import (
	"net/http"

	"github.com/Artymka/avito-entrance-test/internal/config"
	pullRequestCreate "github.com/Artymka/avito-entrance-test/internal/handlers/pullRequest/create"
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
	_, err = users.New(p.DB)
	if err != nil {
		panic(err)
	}
	_, err = teams.New(p.DB)
	if err != nil {
		panic(err)
	}
	pullRequests, err := pullrequests.New(p.DB)
	if err != nil {
		panic(err)
	}
	reviewers, err := reviewers.New(p.DB)
	if err != nil {
		panic(err)
	}

	// router
	r := chi.NewRouter()
	r.Post("/pullRequest/create", pullRequestCreate.New(pullRequests, reviewers))

	// listen
	logging.Info("main", "server is listenning...")
	err = http.ListenAndServe(cfg.Server.Address, r)
	if err != nil {
		logging.Err("server down", err)
	}
}
