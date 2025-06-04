package memdb_test

import (
	"quote_book/pkg/db/memdb"
	"quote_book/pkg/entities"
	"testing"
)

func TestAddQuote(t *testing.T) {
	db := memdb.New()

	// Добавляем валидную цитату
	q := entities.Quote{Text: "Hello", Author: "Author"}
	err := db.AddQuote(q)
	if err != nil {
		t.Fatalf("AddQuote failed: %v", err)
	}

	// Добавляем пустую цитату — должна быть ошибка
	err = db.AddQuote(entities.Quote{Text: ""})
	if err == nil {
		t.Fatal("AddQuote with empty text should return error")
	}
}

func TestGetAllQuotes(t *testing.T) {
	db := memdb.New()

	// Должно быть пусто изначально
	quotes, err := db.GetAllQuotes()
	if err != nil {
		t.Fatalf("GetAllQuotes failed: %v", err)
	}
	if len(quotes) != 0 {
		t.Fatal("GetAllQuotes expected empty slice")
	}

	// Добавляем цитату и проверяем
	_ = db.AddQuote(entities.Quote{Text: "Q1", Author: "A1"})
	quotes, err = db.GetAllQuotes()
	if err != nil {
		t.Fatalf("GetAllQuotes failed: %v", err)
	}
	if len(quotes) != 1 {
		t.Fatal("GetAllQuotes expected 1 quote")
	}
}

func TestGetRandomQuote(t *testing.T) {
	db := memdb.New()

	// Пустая база — ожидаем ошибку
	_, err := db.GetRandomQuote()
	if err == nil {
		t.Fatal("GetRandomQuote on empty DB should return error")
	}

	// Добавляем цитату
	_ = db.AddQuote(entities.Quote{Text: "Q1", Author: "A1"})

	q, err := db.GetRandomQuote()
	if err != nil {
		t.Fatalf("GetRandomQuote failed: %v", err)
	}
	if q.Text != "Q1" || q.Author != "A1" {
		t.Fatal("GetRandomQuote returned wrong quote")
	}
}

func TestDeleteQuote(t *testing.T) {
	db := memdb.New()

	_ = db.AddQuote(entities.Quote{Text: "Q1", Author: "A1"})

	quotes, _ := db.GetAllQuotes()
	if len(quotes) == 0 {
		t.Fatal("no quotes to delete")
	}
	id := quotes[0].ID

	err := db.DeleteQuote(id)
	if err != nil {
		t.Fatalf("DeleteQuote failed: %v", err)
	}

	// Удаление несуществующего id — не ошибка
	err = db.DeleteQuote(9999)
	if err != nil {
		t.Fatalf("DeleteQuote non-existent id should not error: %v", err)
	}

	// Проверяем, что цитата помечена как удалённая (не возвращается)
	quotes, _ = db.GetAllQuotes()
	for _, q := range quotes {
		if q.ID == id {
			t.Fatal("Deleted quote still present in GetAllQuotes")
		}
	}
}

func TestGetAuthorQuotes(t *testing.T) {
	db := memdb.New()

	_ = db.AddQuote(entities.Quote{Text: "Q1", Author: "Author1"})
	_ = db.AddQuote(entities.Quote{Text: "Q2", Author: "Author1"})
	_ = db.AddQuote(entities.Quote{Text: "Q3", Author: "Author2"})

	quotes, err := db.GetAuthorQuotes("Author1")
	if err != nil {
		t.Fatalf("GetAuthorQuotes failed: %v", err)
	}
	if len(quotes) != 2 {
		t.Fatalf("GetAuthorQuotes expected 2 quotes, got %d", len(quotes))
	}

	quotes, err = db.GetAuthorQuotes("Unknown")
	if err != nil {
		t.Fatalf("GetAuthorQuotes failed: %v", err)
	}
	if len(quotes) != 0 {
		t.Fatal("GetAuthorQuotes for unknown author should return empty slice")
	}
}
