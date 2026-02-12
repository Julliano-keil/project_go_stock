package main

import (
	"log"

	"lince/entities"
	"lince/setup"
)

func main() {
	cfg := entities.Config{}

	run := setup.SetupModules(cfg)
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
