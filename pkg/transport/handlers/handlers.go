package handlers

import (
	"encoding/json"
	"log/slog"
	"math/rand"
	"net/http"
	"quote_book/pkg/entities"
	"quote_book/pkg/service"
	"strconv"

	"github.com/gorilla/mux"
)

func NewAddQuoteHandler(qs service.QuoteService, logger *slog.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := *logger.With("requestID", rand.Int63(), "func", "AddQuoteHandler") //better use uuid but only standart lib

		var quote entities.Quote
		err := json.NewDecoder(r.Body).Decode(&quote)
		if err != nil {
			logger.Error("JSON parsing failed", "error", err.Error())
			jsonError(w, http.StatusBadRequest, "bad json")
			return
		}

		err = qs.AddQuote(quote)
		if err != nil {
			logger.Error("Adding quote failed", "error", err.Error())
			jsonError(w, http.StatusInternalServerError, "quote not added")
			return
		}

		logger.Info("Quote added")
		w.WriteHeader(http.StatusCreated)
	}
}

func NewGetQuotesHandler(qs service.QuoteService, logger *slog.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := *logger.With("requestID", rand.Int63(), "func", "GetQuotesHandler")

		author := r.URL.Query().Get("author")

		quotes, err := qs.GetQuotes(author)

		if err != nil {
			logger.Error("Getting quotes failed", "error", err.Error())
			jsonError(w, http.StatusInternalServerError, "getting quotes error")
			return
		}

		jsonQuotes, err := json.Marshal(quotes)
		if err != nil {
			logger.Error("Quotes marshaling failed", "error", err.Error())
			jsonError(w, http.StatusInternalServerError, "quotes marshaling error")
			return
		}

		logger.Info("Quotes recived")
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonQuotes)
	}
}

func NewGetRandomQuotesHandler(qs service.QuoteService, logger *slog.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := *logger.With("requestID", rand.Int63(), "func", "GetRandomQuotesHandler")

		quote, err := qs.GetRandomQuote()
		if err != nil {
			logger.Error("Quote getting failed", "error", err.Error())
			jsonError(w, http.StatusInternalServerError, "quote getting error")
		}

		jsonQuote, err := json.Marshal(quote)
		if err != nil {
			logger.Error("Quote marshaling failed", "error", err.Error())
			jsonError(w, http.StatusInternalServerError, "quotes marshaling error")
			return
		}

		logger.Info("Quote recived")
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonQuote)
	}
}

func NewDeleteQuoteHandler(qs service.QuoteService, logger *slog.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := *logger.With("requestID", rand.Int63(), "func", "DeleteQuoteHandler")

		vars := mux.Vars(r)
		rawID, ok := vars["id"]
		if !ok {
			logger.Error("No id for deleting", "error", "no id")
			jsonError(w, http.StatusBadRequest, "no id in request")
			return
		}

		id, err := strconv.Atoi(rawID)
		if err != nil {
			logger.Error("Not valid id", "error", err.Error())
			jsonError(w, http.StatusBadRequest, "not valid id")
			return
		}

		err = qs.DeleteQuote(id)
		if err != nil {
			logger.Error("Deleting quote error", "error", err.Error())
			jsonError(w, http.StatusInternalServerError, "deleting quote error")
			return
		}

		logger.Info("Quote deleted")
		w.WriteHeader(http.StatusNoContent)
	}
}

func jsonError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
