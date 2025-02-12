// The web controller layer is responsible for handling incoming http requests.
package api

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/adettelle/go-url-shortener/internal/config"
	"github.com/adettelle/go-url-shortener/internal/db"
	"github.com/adettelle/go-url-shortener/internal/logger"
	"go.uber.org/zap"
)

var errlog *zap.Logger = logger.Logger

// Storager defines an interface for interacting with various storage mechanisms,
// such as PathStorage. It describes operations to store, retrieve,
// check for existence, and delete a "address" entity.
type Storager interface {
	GetAddress(shortURL string) (string, error)
	AddAddress(originalURL string) (string, error) // full address
	Finalize() error                               // отрабатывает завершение приложения (при штатном завершении работы)
}

type Handlers struct {
	repo       Storager
	config     *config.Config
	DBCon      db.DBConnector
	Finalizing bool // true означает, что надо делать graceful shutdowm
}

func New(s Storager, cfg *config.Config) *Handlers {
	return &Handlers{
		repo:   s,
		config: cfg,
		DBCon:  db.NewDBConnection(cfg.DBParams),
	}
}

func (h *Handlers) CreateShortAddressPlainText(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var err error

	body, err := io.ReadAll(r.Body)
	if err != nil {
		errlog.Error("error in reading body", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	/*
		shortenAddress, err := helper(h, string(body))
		if err != nil {
			errlog.Error("error in adding address", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	*/
	shortAddress, err := h.repo.AddAddress(string(body))
	if err != nil {
		errlog.Error("error in adding address", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	shortenAddress := h.config.URLAddress + "/" + shortAddress

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(shortenAddress))
	if err != nil {
		errlog.Error("error in writing response", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) GetFullAddress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	id := r.PathValue("id")
	fullAddress, err := h.repo.GetAddress(id)
	if err != nil {
		errlog.Error("error in getting address", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if fullAddress == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Location", fullAddress)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

type shortAddrCreateRequestDTO struct {
	URL string `json:"url"`
}

type shortAddrCreateResponseDTO struct {
	Result string `json:"result"`
}

// func helper(h *Handlers, body string) (string, error) {
// 	shortAddress, err := h.repo.AddAddress(string(body))
// 	if err != nil {
// 		errlog.Error("error in adding address", zap.Error(err))
// 		return "", err
// 	}

// 	shortenAddress := h.config.URLAddress + "/" + shortAddress
// 	return shortenAddress, nil
// }

func (h *Handlers) CreateShortAddressJSON(w http.ResponseWriter, r *http.Request) {
	var requestBody shortAddrCreateRequestDTO
	var buf bytes.Buffer

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		errlog.Error("error in reading body", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Deserialize JSON into requestBody
	if err = json.Unmarshal(buf.Bytes(), &requestBody); err != nil {
		errlog.Error("error in unmarshalling json", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	/*
		shortenAddress, err := helper(h, requestBody.URL)
		if err != nil {
			errlog.Error("error in adding address", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	*/
	shortAddress, err := h.repo.AddAddress(requestBody.URL) // shortAddress is: vN
	if err != nil {
		errlog.Error("error in adding address", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	shortenAddress := h.config.URLAddress + "/" + shortAddress // http://localhost:8000/vN

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	respDTO := shortAddrCreateResponseDTO{Result: shortenAddress}
	resp, err := json.Marshal(respDTO)
	if err != nil {
		errlog.Error("error in marshalling json", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = w.Write(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) CheckConnectionToDB(w http.ResponseWriter, r *http.Request) {
	log.Println("Checking DB")

	_, err := h.DBCon.Connect()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type PostBatchRequestDTO struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type PostBatchResponseDTO struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// PostBatch processes an HTTP POST request
// that contains a batch of urls in JSON format.
func (h *Handlers) PostBatch(w http.ResponseWriter, r *http.Request) {
	var requestBody []PostBatchRequestDTO
	var buf bytes.Buffer
	var shortURL string
	var result []PostBatchResponseDTO

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		errlog.Error("error in reading body", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Deserialize JSON into requestBody
	if err = json.Unmarshal(buf.Bytes(), &requestBody); err != nil {
		errlog.Error("error in unmarshalling json", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(requestBody) == 0 {
		// проверить на пустоту requestBody len == 0, если да то я возвращаю пустой массив
		result = []PostBatchResponseDTO{}
	}
	/*
		shortenAddress, err := helper(h, requestBody.URL)
		if err != nil {
			errlog.Error("error in adding address", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	*/

	for _, elem := range requestBody {
		charStr, err := h.repo.AddAddress(elem.OriginalURL) // charStr is: vN
		if err != nil {
			errlog.Error("error in adding address", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		shortURL = h.config.URLAddress + "/" + charStr // http://localhost:8000/vN

		res := PostBatchResponseDTO{
			CorrelationID: elem.CorrelationID,
			ShortURL:      shortURL,
		}
		result = append(result, res)
		// в этом цикле сразу писать в массив, который я потом верну
	}

	// respDTO := PostBatchResponseDTO{ShortURL: shortURL}
	resp, err := json.Marshal(result) // respDTO
	if err != nil {
		errlog.Error("error in marshalling json", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
}
