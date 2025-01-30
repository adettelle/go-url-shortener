package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/adettelle/go-url-shortener/internal/api"
	"github.com/adettelle/go-url-shortener/internal/config"
	"github.com/adettelle/go-url-shortener/internal/helpers"
	"github.com/adettelle/go-url-shortener/internal/storage"
	"github.com/adettelle/go-url-shortener/pkg/mware"
	"github.com/go-chi/chi/v5"
)

func main() {
	err := initialize()
	if err != nil {
		log.Fatal(err)
	}
}

var addressStorage *storage.AddressStorage

// var readFromFile bool = true

func initialize() error {
	cfg, err := config.New()
	if err != nil {
		return err
	}
	log.Println("Config:", cfg)

	addressStorage = storage.New(cfg.FileStoragePath)

	// если файл по адресу "/tmp/metrics-db.json" есть, но пустой (fi.Size() == 0),
	// далее проверяем, что именно он записан у нас в config.StoragePath
	// и то стираем всё из хранилища addressStorage, т. к. загружать нечего, ведь файл пустой
	fi, err := os.Stat(cfg.FileStoragePath)
	if err != nil && !os.IsNotExist(err) { // если какая-то другая ошибка
		fmt.Println("--------")
		log.Fatal(err)
	}

	if fi != nil && fi.Size() > 0 { // !os.IsNotExist(err) &&
		fmt.Println("!!!!!!!!!!")
		err := helpers.ReadJSONFromFile(cfg.FileStoragePath, addressStorage)
		if err != nil {
			log.Fatal(err)
		}
	}

	go startServer(cfg, addressStorage)

	// создаю канал, который принимает объект типа os.Signal.
	c := make(chan os.Signal, 1)
	// говорим программе, что когда произойдет os.Interrupt,
	// то есть сигнал на штатное завершение, написать в этот канал
	signal.Notify(c, os.Interrupt)
	// вычитка из кaнала, в который еще ничего не записано,
	// блокирует выполнение текущей горутины до тех пор, пока в канале что-то не появится.
	// Когда в канале что-то повится, мы пойдем дальше по коду
	s := <-c
	log.Printf("got termination signal: %s. Grathful shutdown", s)
	err = storage.WriteAddressStorageToJSONFile(cfg.FileStoragePath, addressStorage)
	if err != nil {
		log.Println("unable to write to file")
	}

	// handlers := api.New(addressStorage, cfg)

	// r := chi.NewRouter()
	// r.Post("/", mware.WithLogging(mware.GzipMiddleware(handlers.CreateShortAddressPlainText)))
	// r.Get("/{id}", mware.WithLogging(mware.GzipMiddleware(handlers.GetFullAddress)))
	// r.Post("/api/shorten", mware.WithLogging(mware.GzipMiddleware(handlers.CreateShortAddressJSON)))

	// fmt.Printf("Starting server on port %s\n", cfg.Address)
	// return http.ListenAndServe(cfg.Address, r)
	return nil
}

func startServer(cfg *config.Config, addrStorage *storage.AddressStorage) {
	fmt.Printf("Starting server on port %s\n", cfg.Address)

	handlers := api.New(addrStorage, cfg)

	r := chi.NewRouter()
	r.Post("/", mware.WithLogging(mware.GzipMiddleware(handlers.CreateShortAddressPlainText)))
	r.Get("/{id}", mware.WithLogging(mware.GzipMiddleware(handlers.GetFullAddress)))
	r.Post("/api/shorten", mware.WithLogging(mware.GzipMiddleware(handlers.CreateShortAddressJSON)))

	err := http.ListenAndServe(cfg.Address, r)
	if err != nil {
		log.Fatal(err)
	}
}
