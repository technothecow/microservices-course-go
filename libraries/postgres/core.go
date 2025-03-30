package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

var (
	db   *sql.DB
	once sync.Once
)

func getPostgresConnection() {
	var err error

	username := os.Getenv("POSTGRES_USER")
	if username == "" {
		username = "postgres"
	}
	password := os.Getenv("POSTGRES_PASSWORD")
	if password == "" {
		password = "postgres"
	}
	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		panic("Please provide POSTGRES_HOST")
	}
	port := os.Getenv("POSTGRES_PORT")
	if port == "" {
		port = "5432"
	}
	dbname := os.Getenv("POSTGRES_DB")
	if dbname == "" {
		panic("Please provide POSTGRES_DB")
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, dbname)

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	log.Println("Successfully connected to Postgres")
}

var GetPostgresConnection = func() *sql.DB {
	once.Do(getPostgresConnection)
	return db
}
