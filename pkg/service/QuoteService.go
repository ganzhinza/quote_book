package service

import "quote_book/pkg/entities"

type QuoteService interface {
	AddQuote(quote entities.Quote) error
	GetQuotes(author string) ([]entities.Quote, error)
	GetRandomQuote() (entities.Quote, error)
	DeleteQuote(id int) error
}
