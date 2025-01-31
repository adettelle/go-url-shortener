package api

import (
	"github.com/adettelle/go-url-shortener/pkg/mware"
	"github.com/go-chi/chi/v5"
)

func NewRouter(storager Storager, handlers *Handlers) *chi.Mux {
	r := chi.NewMux()

	r.Post("/", mware.WithLogging(mware.GzipMiddleware(handlers.CreateShortAddressPlainText)))
	r.Get("/{id}", mware.WithLogging(mware.GzipMiddleware(handlers.GetFullAddress)))
	r.Post("/api/shorten", mware.WithLogging(mware.GzipMiddleware(handlers.CreateShortAddressJSON)))
	r.Get("/ping", mware.WithLogging(mware.GzipMiddleware(handlers.CheckConnectionToDB)))

	return r
}
