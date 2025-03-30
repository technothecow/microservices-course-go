package migrator

import (
	"database/sql"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"
)

// Migration represents a migration file.
type Migration struct {
	Name    string // file name (including ordering prefix)
	Content string // the SQL commands inside the file
}

// ApplyMigrations applies all pending migrations found in folder to the given PostgreSQL database.
func ApplyMigrations(db *sql.DB, migrationsFolder string) error {
	// Ensure the migration tracking table exists.
	if err := ensureSchemaMigrationsTable(db); err != nil {
		return err
	}

	// Get a list of migration files.
	migrations, err := readMigrations(migrationsFolder)
	if err != nil {
		return err
	}
	if len(migrations) == 0 {
		log.Println("no migrations found")
		return nil
	}

	// Get already applied migrations.
	applied, err := appliedMigrations(db)
	if err != nil {
		return err
	}

	// Loop through migrations in order, and apply only those that have not been applied.
	for _, m := range migrations {
		if applied[m.Name] {
			continue
		}
		log.Printf("applying migration: %s", m.Name)
		if err := applyMigration(db, m); err != nil {
			return fmt.Errorf("failed applying migration %s: %w", m.Name, err)
		}
		log.Printf("migration %s applied", m.Name)
	}

	return nil
}

// ensureSchemaMigrationsTable creates a table to track migrations.
func ensureSchemaMigrationsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		id SERIAL PRIMARY KEY,
		filename VARCHAR(255) NOT NULL UNIQUE,
		applied_at TIMESTAMP WITH TIME ZONE DEFAULT now()
	);`
	_, err := db.Exec(query)
	return err
}

// appliedMigrations retrieves all migration filenames already applied.
func appliedMigrations(db *sql.DB) (map[string]bool, error) {
	rows, err := db.Query(`SELECT filename FROM schema_migrations;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		applied[name] = true
	}
	return applied, nil
}

// readMigrations reads migration files from the folder and orders them.
func readMigrations(folder string) ([]Migration, error) {
	var migrations []Migration

	// Walk through the folder and look for .sql files.
	err := filepath.WalkDir(folder, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".sql") {
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		m := Migration{
			Name:    d.Name(),
			Content: string(data),
		}
		migrations = append(migrations, m)
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Sort migrations by file name (assumes a valid naming convention).
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Name < migrations[j].Name
	})
	return migrations, nil
}

// applyMigration applies the given migration within a transaction.
func applyMigration(db *sql.DB, m Migration) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// It's a good practice to run the migration in a transaction, if possible.
	// Note: Some migration statements (like certain DDL commands) might not work in transactions.
	if _, err := tx.Exec(m.Content); err != nil {
		tx.Rollback()
		return err
	}

	// Record that the migration was applied.
	query := `INSERT INTO schema_migrations (filename) VALUES ($1);`
	if _, err := tx.Exec(query, m.Name); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
