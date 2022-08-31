package configs

import (
	"github.com/caarlos0/env/v6"
	"log"
)

type config struct {
	ServerPort string `env:"PORT" envDefault:"9999"`
}

func Get() *config {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}
	return &cfg
}
