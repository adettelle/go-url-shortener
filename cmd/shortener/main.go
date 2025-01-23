package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/adettelle/go-url-shortener/internal/api"
	"github.com/adettelle/go-url-shortener/internal/config"
	"github.com/adettelle/go-url-shortener/internal/mware"
	"github.com/adettelle/go-url-shortener/internal/storage"
	"github.com/go-chi/chi/v5"
)

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
	log.Println("Config:", cfg)
	addressStorage := storage.New()

	handlers := api.New(addressStorage, cfg)

	r := chi.NewRouter()
	r.Post("/", mware.WithLogging(handlers.CreateShortAddressPlainText))
	r.Get("/{id}", mware.WithLogging(handlers.GetFullAddress))
	r.Post("/api/shorten", mware.WithLogging(handlers.CreateShortAddressJson))

	fmt.Printf("Starting server on port %s\n", cfg.Address)
	return http.ListenAndServe(cfg.Address, r)
}
