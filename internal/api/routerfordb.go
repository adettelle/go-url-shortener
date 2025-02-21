package api

import (
	"github.com/adettelle/go-url-shortener/pkg/mware"
	"github.com/go-chi/chi/v5"
)

// storager Storager, , urlHandlers *URLsHandlers
func NewRouterForDB(handlers *CustomerHandlers, urlHandlers *Handlers) *chi.Mux {
	r := chi.NewMux()

	// withAuth wraps a given HTTP handler with authentication middleware.
	// withAuth := func(h http.HandlerFunc) http.HandlerFunc {
	// 	return mware.AuthMwr(h, handlers.SignKey, jwtChecker)
	// }

	// User authentication routes
	// if handlers != nil {
	// 	r.Post("/api/user/register", handlers.RegisterCustomer)
	// 	r.Post("/api/user/login", handlers.Login)
	// }

	// URLs management routes
	r.Post("/", mware.WithLogging(mware.GzipMiddleware(urlHandlers.CreateShortAddressPlainText)))
	r.Get("/{id}", mware.WithLogging(mware.GzipMiddleware(urlHandlers.GetFullAddress)))
	r.Post("/api/shorten", mware.WithLogging(mware.GzipMiddleware(urlHandlers.CreateShortAddressJSON)))
	r.Get("/ping", mware.WithLogging(mware.GzipMiddleware(urlHandlers.CheckConnectionToDB)))
	r.Post("/api/shorten/batch", mware.WithLogging(mware.GzipMiddleware(urlHandlers.PostBatch)))
	r.Get("/api/user/urls", mware.WithLogging(mware.GzipMiddleware(urlHandlers.GetAllURLS)))

	return r
}
