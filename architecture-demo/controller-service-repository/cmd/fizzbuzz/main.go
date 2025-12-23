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
	GetFizzBuzz(req *repository.GetFizzBuzzRequest) (*repository.GetFizzBuzzResponse, error)
	SetFizzBuzz(req *repository.SetFizzBuzzRequest)
}

func run() error {
	cfg := config.GetConfig()

	store := repository.NewStore()
	fizzBuzzService := service.NewFizzBuzz(store)

	return handlers.Serve(cfg.Handlers, fizzBuzzService)
}
