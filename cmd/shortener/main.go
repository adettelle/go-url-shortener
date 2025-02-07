package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/adettelle/go-url-shortener/internal/api"
	"github.com/adettelle/go-url-shortener/internal/config"
	"github.com/adettelle/go-url-shortener/internal/db"
	"github.com/adettelle/go-url-shortener/internal/migrator"
	"github.com/adettelle/go-url-shortener/internal/storage/dbstorage"
	"github.com/adettelle/go-url-shortener/internal/storage/urlstorage"
)

func main() {
	err := initializeServer()
	if err != nil {
		log.Fatal(err)
	}
}

// var addressStorage *urlstorage.AddressStorage

// var readFromFile bool = true

func initializeServer() error {
	cfg, err := config.New()
	if err != nil {
		return err
	}
	log.Println("Config:", cfg)

	migrator.MustApplyMigrations(cfg.DBParams)

	storager, err := initStorager(cfg)
	if err != nil {
		return err
	}

	urlAPI := api.New(storager, cfg)
	router := api.NewRouter(storager, urlAPI)
	srv := &http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}
	// addressStorage = urlstorage.New(cfg.FileStoragePath)

	// // если файл по адресу "/tmp/metrics-db.json" есть, но пустой (fi.Size() == 0),
	// // далее проверяем, что именно он записан у нас в config.StoragePath
	// // и то стираем всё из хранилища addressStorage, т. к. загружать нечего, ведь файл пустой
	// fi, err := os.Stat(cfg.FileStoragePath)
	// if err != nil && !os.IsNotExist(err) { // если какая-то другая ошибка
	// 	log.Fatal(err)
	// }

	// if fi != nil && fi.Size() > 0 { // !os.IsNotExist(err) &&
	// 	err := urlstorage.ReadJSONFromFile(cfg.FileStoragePath, storager) // addressStorage
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }

	go startServer(cfg, storager)

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

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatal(err) // failure/timeout shutting down the server gracefully
	}

	urlAPI.Finalizing = true

	err = storager.Finalize()
	if err != nil {
		log.Println(err)
		log.Println("unable to write to file")
	}

	// err = urlstorage.WriteAddressStorageToJSONFile(cfg.FileStoragePath, storager) // addressStorage
	// if err != nil {
	// 	log.Println("unable to write to file")
	// }

	// handlers := api.New(addressStorage, cfg)

	// r := chi.NewRouter()
	// r.Post("/", mware.WithLogging(mware.GzipMiddleware(handlers.CreateShortAddressPlainText)))
	// r.Get("/{id}", mware.WithLogging(mware.GzipMiddleware(handlers.GetFullAddress)))
	// r.Post("/api/shorten", mware.WithLogging(mware.GzipMiddleware(handlers.CreateShortAddressJSON)))

	// fmt.Printf("Starting server on port %s\n", cfg.Address)
	// return http.ListenAndServe(cfg.Address, r)
	return nil
}

func startServer(cfg *config.Config, storager api.Storager) { // addrStorage *urlstorage.AddressStorage
	fmt.Printf("Starting server on port %s\n", cfg.Address)

	handlers := api.New(storager, cfg) // addrStorage

	r := api.NewRouter(storager, handlers) // addrStorage
	// r := chi.NewRouter()
	// r.Post("/", mware.WithLogging(mware.GzipMiddleware(handlers.CreateShortAddressPlainText)))
	// r.Get("/{id}", mware.WithLogging(mware.GzipMiddleware(handlers.GetFullAddress)))
	// r.Post("/api/shorten", mware.WithLogging(mware.GzipMiddleware(handlers.CreateShortAddressJSON)))

	err := http.ListenAndServe(cfg.Address, r)
	if err != nil {
		log.Fatal(err)
	}
}

func initStorager(cfg *config.Config) (api.Storager, error) {
	var storager api.Storager

	if cfg.DBParams != "" {
		db, err := db.NewDBConnection(cfg.DBParams).Connect()
		if err != nil {
			log.Fatal(err)
		}

		storager = &dbstorage.DBStorage{
			Ctx: context.Background(),
			DB:  db,
		}

		migrator.MustApplyMigrations(cfg.DBParams)
	} else {
		var urlSt *urlstorage.AddressStorage
		log.Println("cfg.StoragePath in initStorager:", cfg.FileStoragePath)
		urlSt, err := urlstorage.New(cfg.ShouldRestore(), cfg.FileStoragePath)
		if err != nil {
			return nil, err
		}

		storager = urlSt
	}
	return storager, nil
}
