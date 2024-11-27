package main

import (
	"log"

	"csr/internal/config"
	"csr/internal/handlers"
	"csr/internal/repository"
	"csr/internal/service"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

type Repository interface {
	GetFoobar(req *repository.GetFoobarRequest) (*repository.GetFoobarResponse, error)
	SetFoobar(req *repository.SetFoobarRequest)
}

func run() error {
	cfg := config.GetConfig()

	store := repository.NewStore()
	foobarService := service.NewFoobar(store)

	return handlers.Serve(cfg.Handlers, foobarService)
}
