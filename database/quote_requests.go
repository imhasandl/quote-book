package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/imhasandl/quote-book/database/models"
)

type DBInterface interface {
	InsertQuote(ctx context.Context, params CreateQuoteParams) (models.Quote, error)
	GetAllQuotes(ctx context.Context) ([]models.Quote, error)
	GetQuotesByFilter(ctx context.Context, filter string) ([]models.Quote, error)
	DeleteQuote(ctx context.Context, id int) error
}

type DBQueries struct {
	DB *sql.DB
}

type CreateQuoteParams struct {
	Author string `json:"author"`
	Quote  string `json:"quote"`
}

const InsertQuoteParams = `
INSERT INTO quotes (author, quote) VALUES ($1, $2) RETURNING id, author, quote; 
`

func (r *DBQueries) InsertQuote(ctx context.Context, params CreateQuoteParams) (models.Quote, error) {
	var q models.Quote
	err := r.DB.QueryRowContext(ctx, InsertQuoteParams, params.Author, params.Quote).Scan(&q.ID, &q.Author, &q.Quote)
	if err != nil {
		log.Printf("can't insert quote in db: %v", err)
		return models.Quote{}, err
	}

	return q, nil
}

const GetAllQuotesParams = `
SELECT id, author, quote FROM quotes ORDER BY id ASC;
`

func (r *DBQueries) GetAllQuotes(ctx context.Context) ([]models.Quote, error) {
	var quotes []models.Quote

	rows, err := r.DB.QueryContext(ctx, GetAllQuotesParams)
	if err != nil {
		log.Printf("error scanning quote row: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var q models.Quote
		if err := rows.Scan(&q.ID, &q.Author, &q.Quote); err != nil {
			log.Printf("error scanning rows: %v", err)
			return nil, err
		}

		quotes = append(quotes, q)
	}

	if err = rows.Err(); err != nil {
		fmt.Printf("error iterating quote rows: %v", err)
		return nil, err
	}

	return quotes, nil
}

const GetQuotesByFilterParams = `
SELECT id, author, quote FROM quotes
WHERE author = $1;
`

func (r *DBQueries) GetQuotesByFilter(ctx context.Context, author string) ([]models.Quote, error) {
	var quotes []models.Quote

	rows, err := r.DB.QueryContext(ctx, GetQuotesByFilterParams, author)
	if err != nil {
		fmt.Printf("Can't get author quotes from db: %v", err)
		return nil, err
	}

	for rows.Next() {
		var q models.Quote
		err := rows.Scan(&q.ID, &q.Author, &q.Quote)
		if err != nil {
			fmt.Printf("Can't scan author quotes: %v", err)
			return nil, err
		}
		quotes = append(quotes, q)
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("error iterating over author quote: %v", err)
		return nil, err
	}

	return quotes, nil
}

const DeleteQuoteParams = `
DELETE FROM quotes WHERE id = $1
`

func (r *DBQueries) DeleteQuote(ctx context.Context, id int) error {
	_, err := r.DB.ExecContext(ctx, DeleteQuoteParams, id)
	if err != nil {
		fmt.Printf("can't delete quote from database: %v", err)
		return err
	}

	return nil
}
