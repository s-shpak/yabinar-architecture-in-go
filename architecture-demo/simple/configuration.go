package main

import "flag"

type Configuration struct {
	Limit int
}

func GetConfiguration() *Configuration {
	cfg := &Configuration{}
	flag.IntVar(&cfg.Limit, "n", 20, "run FizzBuzz up to this number")
	flag.Parse()
	return cfg
}
