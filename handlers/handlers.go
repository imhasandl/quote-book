package handlers

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/imhasandl/quote-book/database"
	"github.com/imhasandl/quote-book/database/models"
	"github.com/imhasandl/quote-book/helper"
)

type apiConfig struct {
	db database.DBInterface
}

func NewConfig(db database.DBInterface) *apiConfig {
	return &apiConfig{
		db: db,
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
	author := req.URL.Query().Get("author")

	var quotes []models.Quote
	var err error

	if author != "" {
		quotes, err = cfg.db.GetQuotesByFilter(context.Background(), author)
	} else {
		quotes, err = cfg.db.GetAllQuotes(context.Background())
	}

	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "can't get quotes from database", err)
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, quotes)
}

func (cfg *apiConfig) RandomQuote(w http.ResponseWriter, req *http.Request) {
	quotes, err := cfg.db.GetAllQuotes(context.Background())
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "can't get all quotes from database", err)
		return
	}

	randomQuote := rand.Intn(len(quotes))

	helper.RespondWithJSON(w, http.StatusOK, quotes[randomQuote])
}

func (cfg *apiConfig) DeleteQuote(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	idStr, ok := vars["id"]
	if !ok || idStr == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "id parameter is missing in URL", nil)
		return
	}

	quoteID, err := strconv.Atoi(idStr)
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "can't convert string id to int", err)
		return
	}

	err = cfg.db.DeleteQuote(context.Background(), quoteID)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "can't delete quote by id", err)
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, map[string]bool{
		"status": true,
	})
}
