package db

import "quote_book/pkg/entities"

type DB interface {
	AddQuote(quote entities.Quote) error
	GetAllQuotes() ([]entities.Quote, error)
	GetRandomQuote() (entities.Quote, error)
	GetAuthorQuotes(author string) ([]entities.Quote, error)
	DeleteQuote(id int) error
}
