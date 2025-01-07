package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func PostMainPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	_, err := io.ReadAll(r.Body) // body
	if err != nil {
		http.Error(w, "Error in body", http.StatusBadRequest) // StatusInternalServerError
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("http://localhost:8080/EwHXdJfB")) // http://localhost:8080/EwHXdJfB
	// w.Write([]byte("Hello"))

}

func GetID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	id := r.PathValue("id")
	if id == "EwHXdJfB" {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusTemporaryRedirect)
		w.Write([]byte("Location: https://practicum.yandex.ru/"))
	} else {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusTemporaryRedirect)
		w.Write([]byte("Location: https://practicum.yandex.ru/")) // ???
	}

}

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	r := http.NewServeMux() // // создаем сервер
	r.HandleFunc("POST /", PostMainPage)
	r.HandleFunc("GET /{id}", GetID)
	port := ":8080"
	fmt.Printf("Starting server on port %s\n", port)
	return (http.ListenAndServe(port, r))
}
