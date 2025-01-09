// The web controller layer is responsible for handling incoming http requests.
package api

import (
	"io"
	"net/http"
)

// Storager defines an interface for interacting with various storage mechanisms,
// such as PathStorage. It describes operations to store, retrieve,
// check for existence, and delete a "path" entity.
type Storager interface {
	GetPath(name string) string
	AddPath(fullPath string) string
}

type Handlers struct {
	repo Storager
}

func New(s Storager) *Handlers {
	return &Handlers{
		repo: s,
	}
}

func (h *Handlers) PostShortPath(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error in body", http.StatusBadRequest) // StatusInternalServerError
		return
	}

	shortPath := h.repo.AddPath(string(body))

	shortenPath := "http://localhost:8080/" + shortPath

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortenPath))

}

func (h *Handlers) GetID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	id := r.PathValue("id")

	w.Header().Set("Location", h.repo.GetPath(id))
	w.WriteHeader(http.StatusTemporaryRedirect)
}
