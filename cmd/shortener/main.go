package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/adettelle/go-url-shortener/internal/api"
	"github.com/adettelle/go-url-shortener/internal/config"
	"github.com/adettelle/go-url-shortener/internal/storage"
	"github.com/go-chi/chi/v5"
)

// var pathStorage = storage.New()

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.New()
	if err != nil {
		return err
	}

	pathStorage := storage.New()

	handlers := api.New(pathStorage, cfg)

	r := chi.NewRouter()
	r.Post("/", handlers.PostShortPath)
	r.Get("/{id}", handlers.GetID)

	//port := ":8080"
	fmt.Printf("Starting server on port %s\n", cfg.Address) // port
	return (http.ListenAndServe(cfg.Address, r))            // port
}
