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

func initializeServer() error {
	cfg, err := config.New()
	if err != nil {
		return err
	}
	log.Println("Config:", cfg)

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

	return nil
}

func startServer(cfg *config.Config, storager api.Storager) {
	fmt.Printf("Starting server on port %s\n", cfg.Address)

	handlers := api.New(storager, cfg)

	r := api.NewRouter(storager, handlers)

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
