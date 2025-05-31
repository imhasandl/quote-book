package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	"github.com/imhasandl/quote-book/database"
	"github.com/imhasandl/quote-book/database/models"
)

type mockDB struct {
	insertQuoteFunc       func(ctx context.Context, params database.CreateQuoteParams) (models.Quote, error)
	getAllQuotesFunc      func(ctx context.Context) ([]models.Quote, error)
	getQuotesByFilterFunc func(ctx context.Context, filter string) ([]models.Quote, error)
	deleteQuoteFunc       func(ctx context.Context, id int) error
}

func (m *mockDB) InsertQuote(ctx context.Context, params database.CreateQuoteParams) (models.Quote, error) {
	return m.insertQuoteFunc(ctx, params)
}
func (m *mockDB) GetAllQuotes(ctx context.Context) ([]models.Quote, error) {
	return m.getAllQuotesFunc(ctx)
}
func (m *mockDB) GetQuotesByFilter(ctx context.Context, filter string) ([]models.Quote, error) {
	return m.getQuotesByFilterFunc(ctx, filter)
}
func (m *mockDB) DeleteQuote(ctx context.Context, id int) error {
	return m.deleteQuoteFunc(ctx, id)
}

func TestCreateQuote(t *testing.T) {
	tests := []struct {
		name       string
		body       interface{}
		mockInsert func(ctx context.Context, params database.CreateQuoteParams) (models.Quote, error)
		wantStatus int
	}{
		{
			name: "success",
			body: database.CreateQuoteParams{Author: "A", Quote: "Q"},
			mockInsert: func(ctx context.Context, params database.CreateQuoteParams) (models.Quote, error) {
				return models.Quote{ID: 1, Author: params.Author, Quote: params.Quote}, nil
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "bad request",
			body: "not a struct",
			mockInsert: func(ctx context.Context, params database.CreateQuoteParams) (models.Quote, error) {
				return models.Quote{}, nil
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "db error",
			body: database.CreateQuoteParams{Author: "A", Quote: "Q"},
			mockInsert: func(ctx context.Context, params database.CreateQuoteParams) (models.Quote, error) {
				return models.Quote{}, errors.New("db error")
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockDB{
				insertQuoteFunc: tt.mockInsert,
			}
			cfg := &apiConfig{db: mock}

			var bodyBytes []byte
			if s, ok := tt.body.(string); ok {
				bodyBytes = []byte(s)
			} else {
				bodyBytes, _ = json.Marshal(tt.body)
			}

			req := httptest.NewRequest("POST", "/quotes", bytes.NewReader(bodyBytes))
			w := httptest.NewRecorder()
			cfg.CreateQuote(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestGetQuotes(t *testing.T) {
	tests := []struct {
		name         string
		query        string
		mockAll      func(ctx context.Context) ([]models.Quote, error)
		mockByFilter func(ctx context.Context, filter string) ([]models.Quote, error)
		wantStatus   int
	}{
		{
			name:  "get all quotes success",
			query: "",
			mockAll: func(ctx context.Context) ([]models.Quote, error) {
				return []models.Quote{{ID: 1, Author: "A", Quote: "Q"}}, nil
			},
			mockByFilter: func(ctx context.Context, filter string) ([]models.Quote, error) {
				return nil, nil
			},
			wantStatus: http.StatusOK,
		},
		{
			name:  "get quotes by filter success",
			query: "?author=A",
			mockAll: func(ctx context.Context) ([]models.Quote, error) {
				return nil, nil
			},
			mockByFilter: func(ctx context.Context, filter string) ([]models.Quote, error) {
				return []models.Quote{{ID: 2, Author: "A", Quote: "Q2"}}, nil
			},
			wantStatus: http.StatusOK,
		},
		{
			name:  "db error",
			query: "",
			mockAll: func(ctx context.Context) ([]models.Quote, error) {
				return nil, errors.New("db error")
			},
			mockByFilter: func(ctx context.Context, filter string) ([]models.Quote, error) {
				return nil, nil
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockDB{
				getAllQuotesFunc:      tt.mockAll,
				getQuotesByFilterFunc: tt.mockByFilter,
			}
			cfg := &apiConfig{db: mock}

			req := httptest.NewRequest("GET", "/quotes"+tt.query, nil)
			w := httptest.NewRecorder()
			cfg.GetQuotes(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestRandomQuote(t *testing.T) {
	tests := []struct {
		name       string
		mockAll    func(ctx context.Context) ([]models.Quote, error)
		wantStatus int
	}{
		{
			name: "success",
			mockAll: func(ctx context.Context) ([]models.Quote, error) {
				return []models.Quote{{ID: 1, Author: "A", Quote: "Q"}}, nil
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "db error",
			mockAll: func(ctx context.Context) ([]models.Quote, error) {
				return nil, errors.New("db error")
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockDB{
				getAllQuotesFunc: tt.mockAll,
			}
			cfg := &apiConfig{db: mock}

			req := httptest.NewRequest("GET", "/quotes/random", nil)
			w := httptest.NewRecorder()
			cfg.RandomQuote(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestDeleteQuote(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		mockDelete func(ctx context.Context, id int) error
		wantStatus int
	}{
		{
			name: "success",
			id:   "1",
			mockDelete: func(ctx context.Context, id int) error {
				return nil
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "missing id",
			id:   "",
			mockDelete: func(ctx context.Context, id int) error {
				return nil
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid id",
			id:   "abc",
			mockDelete: func(ctx context.Context, id int) error {
				return nil
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "db error",
			id:   "2",
			mockDelete: func(ctx context.Context, id int) error {
				return errors.New("db error")
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockDB{
				deleteQuoteFunc: tt.mockDelete,
			}
			cfg := &apiConfig{db: mock}

			url := "/quotes"
			if tt.id != "" {
				url += "/" + tt.id
			}
			req := httptest.NewRequest("DELETE", url, nil)
			req = setMuxVars(req, map[string]string{"id": tt.id})

			w := httptest.NewRecorder()
			cfg.DeleteQuote(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func setMuxVars(r *http.Request, vars map[string]string) *http.Request {
	return mux.SetURLVars(r, vars)
}
