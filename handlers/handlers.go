package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/imhasandl/quote-book/database"
	"github.com/imhasandl/quote-book/helper"
)

type apiConfig struct {
	db *database.DBQueries
}

func NewConfig(db *sql.DB) *apiConfig {
	return &apiConfig{
		db: &database.DBQueries{DB: db},
	}
}

func (cfg *apiConfig) CreateQuote(w http.ResponseWriter, req *http.Request) {
	var CreateQuoteParams database.CreateQuoteParams
	err := json.NewDecoder(req.Body).Decode(&CreateQuoteParams)
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	quote, err := cfg.db.InsertQuote(context.Background(), CreateQuoteParams)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "can't insert quote into database", err)
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, quote)
}

func (cfg *apiConfig) GetQuotes(w http.ResponseWriter, req *http.Request) {

}

func (cfg *apiConfig) RandomQuote(w http.ResponseWriter, req *http.Request) {

}

func (cfg *apiConfig) FilterByAuthor(w http.ResponseWriter, req *http.Request) {
	author := req.URL.Query().Get("author")
	if author == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "author parameter is required", nil)
		return
	}
}

func (cfg *apiConfig) DeleteQuote(w http.ResponseWriter, req *http.Request) {

}
