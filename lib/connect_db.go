package lib

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

func ConnectDB() *sql.DB {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var db *sql.DB
	var err error

	// Retry connection with backoff
	for i := 0; i < 5; i++ {
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("Database connection failed (attempt %d): %v", i+1, err)
			time.Sleep(time.Duration(i) * 2 * time.Second)
			continue
		}

		if err = db.Ping(); err == nil {
			break
		}
	}

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Successfully connected to database")
	return db
}
