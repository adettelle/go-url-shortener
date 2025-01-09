package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/adettelle/go-url-shortener/internal/api"
	"github.com/adettelle/go-url-shortener/internal/storage"
)

// var pathStorage = storage.New()

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	pathStorage := storage.New()

	handlers := api.New(pathStorage)

	r := http.NewServeMux() // создаем сервер
	r.HandleFunc("POST /", handlers.PostShortPath)
	r.HandleFunc("GET /{id}", handlers.GetID)
	port := ":8080"
	fmt.Printf("Starting server on port %s\n", port)
	return (http.ListenAndServe(port, r))
}
