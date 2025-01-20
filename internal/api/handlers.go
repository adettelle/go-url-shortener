// The web controller layer is responsible for handling incoming http requests.
package api

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/adettelle/go-url-shortener/internal/config"
)

// Storager defines an interface for interacting with various storage mechanisms,
// such as PathStorage. It describes operations to store, retrieve,
// check for existence, and delete a "path" entity.
type Storager interface {
	GetPath(name string) (string, error)
	AddPath(fullPath string) (string, error)
}

type Handlers struct {
	repo   Storager
	config *config.Config
}

func New(s Storager, cfg *config.Config) *Handlers {
	return &Handlers{
		repo:   s,
		config: cfg,
	}
}

func (h *Handlers) PostShortPath(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("error in writing reading body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortPath, err := h.repo.AddPath(string(body))
	if err != nil {
		log.Println("error in adding path")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	shortenPath := h.config.URLAddress + "/" + shortPath

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(shortenPath))
	if err != nil {
		log.Println("error in writing response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) GetID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	id := r.PathValue("id")
	fullPath, err := h.repo.GetPath(id)
	if err != nil {
		log.Println("error in getting path")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if fullPath == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Location", fullPath)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

type ShortAddrCreateRequestDTO struct {
	URL string `json:"url"`
}

type ShortAddrCreateResponseDTO struct {
	Result string `json:"result"`
}

func (h *Handlers) ShortAddressCreate(w http.ResponseWriter, r *http.Request) {
	var requestBody ShortAddrCreateRequestDTO
	var buf bytes.Buffer

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		log.Println("error in writing reading body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Deserialize JSON into requestBody
	if err = json.Unmarshal(buf.Bytes(), &requestBody); err != nil {
		log.Println("error in unmarshalling json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortPath, err := h.repo.AddPath(string(requestBody.URL)) // shortPath is: vN
	if err != nil {
		log.Println("error in adding path")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	shortenPath := h.config.URLAddress + "/" + shortPath // http://localhost:8000/vN

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	respDTO := ShortAddrCreateResponseDTO{Result: shortenPath}
	resp, err := json.Marshal(respDTO)
	if err != nil {
		log.Println("error in marshalling json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = w.Write(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
