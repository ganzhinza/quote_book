package service

import (
	"errors"
	"quote_book/pkg/db"
	"quote_book/pkg/entities"
)

type quoteServiceImpl struct {
	db db.DB
}

func NewQuoteService(db db.DB) *quoteServiceImpl {
	return &quoteServiceImpl{db: db}
}

func (qs *quoteServiceImpl) AddQuote(quote entities.Quote) error {
	err := qs.db.AddQuote(quote)
	if err != nil {
		return errors.Join(errors.New("service AddQuote: "), err)
	}
	return nil
}

func (qs *quoteServiceImpl) GetQuotes(author string) ([]entities.Quote, error) {
	var quotes []entities.Quote
	var err error

	if author == "" {
		quotes, err = qs.db.GetAllQuotes()
	} else {
		quotes, err = qs.db.GetAuthorQuotes(author)
	}
	if err != nil {
		return []entities.Quote{}, errors.Join(errors.New("service GetQuotes: "), err)
	}

	return quotes, nil
}

func (qs *quoteServiceImpl) GetRandomQuote() (entities.Quote, error) {
	quotes, err := qs.db.GetRandomQuote()
	if err != nil {
		return entities.Quote{}, errors.Join(errors.New("service GetRandomQuote: "), err)
	}

	return quotes, nil
}

func (qs *quoteServiceImpl) DeleteQuote(id int) error {
	err := qs.db.DeleteQuote(id)
	if err != nil {
		return errors.Join(errors.New("service DeleteQuote: "), qs.db.DeleteQuote(id))
	}
	return err
}
