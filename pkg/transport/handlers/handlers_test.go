package handlers_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"quote_book/pkg/db/memdb"
	"quote_book/pkg/entities"
	"quote_book/pkg/service"
	"quote_book/pkg/transport/handlers"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
)

var (
	svc    service.QuoteService
	logger *slog.Logger
	router *mux.Router
)

func TestMain(m *testing.M) {
	db := memdb.New()
	svc = service.NewQuoteService(db)
	logger = slog.Default()
	router = setupRouter()
	os.Exit(m.Run())
}

func setupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/quotes", handlers.NewAddQuoteHandler(svc, logger)).Methods(http.MethodPost)
	r.HandleFunc("/quotes", handlers.NewGetQuotesHandler(svc, logger)).Methods(http.MethodGet)
	r.HandleFunc("/quotes/random", handlers.NewGetRandomQuotesHandler(svc, logger)).Methods(http.MethodGet)
	r.HandleFunc("/quotes/{id}", handlers.NewDeleteQuoteHandler(svc, logger)).Methods(http.MethodDelete)
	return r
}

func TestAddGetDeleteQuote(t *testing.T) {

	// Добавляем цитату
	quote := entities.Quote{Text: "Test quote", Author: "Tester"}
	body, _ := json.Marshal(quote)

	req := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("AddQuote: expected status %d, got %d", http.StatusCreated, w.Code)
	}

	// Получаем все цитаты
	req = httptest.NewRequest(http.MethodGet, "/quotes", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GetQuotes: expected status %d, got %d", http.StatusOK, w.Code)
	}

	var quotes []entities.Quote
	err := json.NewDecoder(w.Body).Decode(&quotes)
	if err != nil {
		t.Fatalf("GetQuotes: decode error: %v", err)
	}

	if len(quotes) == 0 {
		t.Fatal("GetQuotes: expected at least one quote")
	}

	// Проверяем, что добавленная цитата есть
	found := false
	for _, q := range quotes {
		if q.Text == quote.Text && q.Author == quote.Author {
			found = true
		}
	}
	if !found {
		t.Fatal("GetQuotes: added quote not found")
	}

	// Получаем случайную цитату
	req = httptest.NewRequest(http.MethodGet, "/quotes/random", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GetRandomQuote: expected status %d, got %d", http.StatusOK, w.Code)
	}

	var randomQuote entities.Quote
	err = json.NewDecoder(w.Body).Decode(&randomQuote)
	if err != nil {
		t.Fatalf("GetRandomQuote: decode error: %v", err)
	}

	if randomQuote.Text == "" {
		t.Fatal("GetRandomQuote: empty quote text")
	}

	// Удаляем цитату по ID
	// Для удаления используем ID из случайной цитаты
	req = httptest.NewRequest(http.MethodDelete, "/quotes/"+strconv.Itoa(randomQuote.ID), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("DeleteQuote: expected status %d, got %d", http.StatusNoContent, w.Code)
	}

	// Проверяем, что цитата удалена
	req = httptest.NewRequest(http.MethodGet, "/quotes", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GetQuotes after delete: expected status %d, got %d", http.StatusOK, w.Code)
	}

	err = json.NewDecoder(w.Body).Decode(&quotes)
	if err != nil {
		t.Fatalf("GetQuotes after delete: decode error: %v", err)
	}

	for _, q := range quotes {
		if q.ID == randomQuote.ID {
			t.Fatal("DeleteQuote: quote still exists after deletion")
		}
	}
}

func TestAddQuoteBadJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewReader([]byte("{bad json")))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("AddQuote with bad json: expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestDeleteQuoteInvalidID(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/quotes/not-an-id", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("DeleteQuote with invalid id: expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestDeleteQuoteNoID(t *testing.T) {

	// Путь без id, должен вернуть 404 (mux не найдет маршрут)
	req := httptest.NewRequest(http.MethodDelete, "/quotes/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("DeleteQuote with no id: expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestGetRandomQuote(t *testing.T) {
	// Сначала добавим цитату, чтобы она была в базе
	quote := entities.Quote{Text: "Random test quote", Author: "Random Tester"}
	err := svc.AddQuote(quote)
	if err != nil {
		t.Fatalf("failed to add quote for test: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/quotes/random", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GetRandomQuote: expected status %d, got %d", http.StatusOK, w.Code)
	}

	var gotQuote entities.Quote
	err = json.NewDecoder(w.Body).Decode(&gotQuote)
	if err != nil {
		t.Fatalf("GetRandomQuote: decode error: %v", err)
	}

	if gotQuote.Text == "" {
		t.Fatal("GetRandomQuote: empty quote text")
	}
	if gotQuote.Author == "" {
		t.Fatal("GetRandomQuote: empty quote author")
	}
}

func TestGetRandomQuoteEmptyDB(t *testing.T) {
	// Создаём новый сервис с пустой базой
	db := memdb.New()
	emptySvc := service.NewQuoteService(db)
	logger := slog.Default()

	r := mux.NewRouter()
	r.HandleFunc("/quotes/random", handlers.NewGetRandomQuotesHandler(emptySvc, logger)).Methods(http.MethodGet)

	req := httptest.NewRequest(http.MethodGet, "/quotes/random", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("GetRandomQuote empty DB: expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	var errResp map[string]string
	err := json.NewDecoder(w.Body).Decode(&errResp)
	if err != nil {
		t.Fatalf("GetRandomQuote empty DB: decode error: %v", err)
	}

	if errMsg, ok := errResp["error"]; !ok || errMsg == "" {
		t.Fatal("GetRandomQuote empty DB: expected error message in response")
	}
}
