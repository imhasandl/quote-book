package database

import (
	"database/sql"
	"fmt"
	"log"
)

var DB *sql.DB

func InitDatabase(db_url string) error {
	var err error
	DB, err = sql.Open("postgres", db_url)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("Database connected")
	return nil
}

func CloseDatabase(){
	if DB != nil {
		DB.Close()
		log.Println("Database connection closed")
	}
}