package api

import (
	"log/slog"
	"net/http"
	"quote_book/pkg/db"
	"quote_book/pkg/service"
	"quote_book/pkg/transport/handlers"

	"github.com/gorilla/mux"
)

type API struct {
	logger *slog.Logger
	router *mux.Router
}

func (api *API) Router() *mux.Router {
	return api.router
}

func New(db db.DB, logger *slog.Logger) *API {
	quoteService := service.NewQuoteService(db)

	api := API{
		router: mux.NewRouter(),
		logger: logger,
	}

	api.endpoints(quoteService)

	return &api
}

func (api *API) endpoints(qs service.QuoteService) {

	api.router.HandleFunc("/quotes", handlers.NewAddQuoteHandler(qs, api.logger)).Methods(http.MethodPost)
	api.router.HandleFunc("/quotes", handlers.NewGetQuotesHandler(qs, api.logger)).Methods(http.MethodGet)
	api.router.HandleFunc("/quotes/random", handlers.NewGetRandomQuotesHandler(qs, api.logger)).Methods(http.MethodGet)
	api.router.HandleFunc("/quotes/{id}", handlers.NewDeleteQuoteHandler(qs, api.logger)).Methods(http.MethodDelete)
}
