package main

import (
	"net/http"

	"github.com/Artymka/avito-entrance-test/internal/config"
	pullRequestCreate "github.com/Artymka/avito-entrance-test/internal/handlers/pullRequest/create"
	"github.com/Artymka/avito-entrance-test/internal/lib/logging"
	"github.com/Artymka/avito-entrance-test/internal/storage/postgres"
	pullrequests "github.com/Artymka/avito-entrance-test/internal/storage/postgres/pullRequests"
	"github.com/Artymka/avito-entrance-test/internal/storage/postgres/reviewers"
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
	pr, err := pullrequests.New(p.DB)
	if err != nil {
		panic(err)
	}
	re, err := reviewers.New(p.DB)
	if err != nil {
		panic(err)
	}

	// router
	r := chi.NewRouter()
	r.Post("/pullRequest/create", pullRequestCreate.New(pr, re))

	// listen
	err = http.ListenAndServe(cfg.Server.Address, r)
	if err != nil {
		logging.Err("server down", err)
	}
}
