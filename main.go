package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/imhasandl/quote-book/database"
	"github.com/imhasandl/quote-book/handlers"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("can not load env file: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("port must be set in .env file")
	}

	db_url := os.Getenv("DB_URL")
	if db_url == "" {
		log.Fatal("can not get db url from .env file")
	}

	if err := database.InitDatabase(db_url); err != nil {
		log.Fatalf("can't init database connection: %v", err)
	}
	defer database.CloseDatabase()

	apiConfig := handlers.NewConfig(database.DB)

	router := mux.NewRouter()
	
	router.HandleFunc("/quotes", apiConfig.CreateQuote).Methods("POST")
	router.HandleFunc("/quotes", apiConfig.GetQuotes).Methods("GET")
	router.HandleFunc("/quotes/random", apiConfig.RandomQuote).Methods("GET")
	router.HandleFunc("/quotes/{id}", apiConfig.DeleteQuote).Methods("DELETE")

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	fmt.Printf("Server running on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
