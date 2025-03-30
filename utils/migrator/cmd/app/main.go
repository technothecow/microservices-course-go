package main

import (
	"fmt"
	"log"
	"os"
	"sn/libraries/postgres"

	"sn/utils/migrator/internal" // update to match your module path
)

func main() {
	migrationsFolder := os.Getenv("MIGRATIONS_FOLDER")
	if migrationsFolder == "" {
		log.Fatalf("MIGRATIONS_FOLDER environment variable not set, skipping migrations")
		return
	}

	db := postgres.GetPostgresConnection()
	defer db.Close()

	// Test the connection.
	if err := db.Ping(); err != nil {
		log.Fatalf("unable to connect: %v", err)
	}
	fmt.Println("Connected to database")

	// Apply migrations; make sure migrations folder contains .sql files.
	if err := migrator.ApplyMigrations(db, migrationsFolder); err != nil {
		log.Fatalf("migration error: %v", err)
	}

	fmt.Println("Migrations applied successfully!")
}
