package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

var (
	once      sync.Once
	db        *sql.DB
	dbInitErr error
)

func New() *sql.DB {
	once.Do(func() {
		connectionString := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"))
		fmt.Println(connectionString)
		db, dbInitErr = sql.Open("postgres", connectionString)
		if dbInitErr != nil {
			log.Fatal("Error initializing the database", dbInitErr)
		}

		if dbInitErr = db.Ping(); dbInitErr != nil {
			log.Fatal("Error pinging the database", dbInitErr)
		}

		fmt.Println("Connected to database!")
	})
	return db
}
