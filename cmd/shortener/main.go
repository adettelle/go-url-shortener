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
	w.Write([]byte("http://localhost:8080/EwHXdJfB"))

}

var data = map[string]string{"EwHXdJfB": "https://practicum.yandex.ru/"}

func GetID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	id := r.PathValue("id")

	w.Header().Set("Location", data[id])
	w.WriteHeader(http.StatusTemporaryRedirect)
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
