package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/adettelle/go-url-shortner/internal/storage"
)

var pathStorage = storage.PathStorage{
	Paths: map[string]string{},
}

func PostShortPath(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error in body", http.StatusBadRequest) // StatusInternalServerError
		return
	}

	shortPath := pathStorage.AddPath(string(body))

	shortenPath := "http://localhost:8080/" + shortPath

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortenPath))

}

func GetID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	id := r.PathValue("id")

	w.Header().Set("Location", pathStorage.GetPath(id))
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	r := http.NewServeMux() // создаем сервер
	r.HandleFunc("POST /", PostShortPath)
	r.HandleFunc("GET /{id}", GetID)
	port := ":8080"
	fmt.Printf("Starting server on port %s\n", port)
	return (http.ListenAndServe(port, r))
}
