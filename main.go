package main

import (
	"log"
	"net/http"

	"lince/entities"
	"lince/setup"

	"github.com/gorilla/mux"
)

func main() {
	cfg := entities.Config{}
	r := mux.NewRouter()

	setup.SetupModules(r, cfg)

	srv := &http.Server{Addr: ":8080", Handler: r}
	log.Printf("HTTP server listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
