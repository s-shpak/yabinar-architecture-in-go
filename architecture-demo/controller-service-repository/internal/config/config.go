package config

import (
	handlersConfig "csr/internal/handlers/config"
	"flag"
)

type Config struct {
	Handlers handlersConfig.Config
}

func GetConfig() Config {
	cfg := Config{}
	flag.StringVar(&cfg.Handlers.ServerAddr, "addr", "localhost:8080", "address of HTTP server")

	flag.Parse()
	return cfg
}
