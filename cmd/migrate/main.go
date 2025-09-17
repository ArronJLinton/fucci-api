//go:build migrate
// +build migrate

package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/lib/pq"
)

type Migration struct {
	Version string
	Name    string
	Up      string
	Down    string
}

func main() {
	// Get database URL from environment
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Println("Connected to database successfully")

	// Create migrations table if it doesn't exist
	createMigrationsTable(db)

	// Get all migration files
	migrationFiles, err := getMigrationFiles("sql/schema")
	if err != nil {
		log.Fatal("Failed to get migration files:", err)
	}

	// Get already executed migrations
	executedMigrations, err := getExecutedMigrations(db)
	if err != nil {
		log.Fatal("Failed to get executed migrations:", err)
	}

	// Run pending migrations
	pendingMigrations := getPendingMigrations(migrationFiles, executedMigrations)

	if len(pendingMigrations) == 0 {
		fmt.Println("No pending migrations found")
		return
	}

	fmt.Printf("Found %d pending migrations\n", len(pendingMigrations))

	for _, migration := range pendingMigrations {
		fmt.Printf("Running migration: %s\n", migration.Name)

		// Read migration file
		content, err := ioutil.ReadFile(filepath.Join("sql/schema", migration.Name))
		if err != nil {
			log.Fatalf("Failed to read migration file %s: %v", migration.Name, err)
		}

		// Extract UP migration (everything before -- +goose Down)
		upMigration := extractUpMigration(string(content))

		// Execute migration
		if _, err := db.Exec(upMigration); err != nil {
			log.Fatalf("Migration %s failed: %v", migration.Name, err)
		}

		// Record migration as executed
		if err := recordMigration(db, migration.Version, migration.Name); err != nil {
			log.Fatalf("Failed to record migration %s: %v", migration.Name, err)
		}

		fmt.Printf("Migration %s completed successfully\n", migration.Name)
	}

	fmt.Println("All pending migrations completed successfully!")
}

func createMigrationsTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version VARCHAR(255) PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	if _, err := db.Exec(query); err != nil {
		log.Fatal("Failed to create migrations table:", err)
	}
}

func getMigrationFiles(schemaDir string) ([]Migration, error) {
	files, err := ioutil.ReadDir(schemaDir)
	if err != nil {
		return nil, err
	}

	var migrations []Migration
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			// Extract version from filename (e.g., "001_users.sql" -> "001")
			version := strings.Split(file.Name(), "_")[0]
			migrations = append(migrations, Migration{
				Version: version,
				Name:    file.Name(),
			})
		}
	}

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

func getExecutedMigrations(db *sql.DB) (map[string]bool, error) {
	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	executed := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		executed[version] = true
	}

	return executed, nil
}

func getPendingMigrations(allMigrations []Migration, executed map[string]bool) []Migration {
	var pending []Migration
	for _, migration := range allMigrations {
		if !executed[migration.Version] {
			pending = append(pending, migration)
		}
	}
	return pending
}

func extractUpMigration(content string) string {
	// Find the -- +goose Down marker
	downMarker := "-- +goose Down"
	downIndex := strings.Index(content, downMarker)

	if downIndex == -1 {
		// No down migration, return everything after -- +goose Up
		upMarker := "-- +goose Up"
		upIndex := strings.Index(content, upMarker)
		if upIndex == -1 {
			return content
		}
		return strings.TrimSpace(content[upIndex+len(upMarker):])
	}

	// Return everything between -- +goose Up and -- +goose Down
	upMarker := "-- +goose Up"
	upIndex := strings.Index(content, upMarker)
	if upIndex == -1 {
		return strings.TrimSpace(content[:downIndex])
	}

	upContent := content[upIndex+len(upMarker) : downIndex]
	return strings.TrimSpace(upContent)
}

func recordMigration(db *sql.DB, version, name string) error {
	query := "INSERT INTO schema_migrations (version, name) VALUES ($1, $2)"
	_, err := db.Exec(query, version, name)
	return err
}
