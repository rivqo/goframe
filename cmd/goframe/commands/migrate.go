package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/example/goframe/config"
	"github.com/example/goframe/db"
)

// Migrate runs all pending migrations
func Migrate(cfg *config.Config) {

	dbConfig := db.DatabaseConfig{
		Driver:   cfg.Database.Driver,
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		Name:     cfg.Database.Name,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
	}
	
	fmt.Println("Running migrations...")

	database, err := db.NewDatabase(&dbConfig)
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	// Get all migration files
	migrations, err := getMigrationFiles()
	if err != nil {
		fmt.Printf("Failed to get migration files: %v\n", err)
		os.Exit(1)
	}

	// Run migrations
	migrator := db.NewMigrator(database)
	count, err := migrator.RunMigrations(migrations)
	if err != nil {
		fmt.Printf("Failed to run migrations: %v\n", err)
		os.Exit(1)
	}

	if count == 0 {
		fmt.Println("No migrations to run")
	} else {
		fmt.Printf("Ran %d migrations\n", count)
	}
}

// MigrateRollback rolls back the last batch of migrations
func MigrateRollback(cfg *config.Config, step int) {
	fmt.Println("Rolling back migrations...")

	dbConfig := db.DatabaseConfig{
		Driver:   cfg.Database.Driver,
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		Name:     cfg.Database.Name,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
	}

	database, err := db.NewDatabase(&dbConfig)
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	// Get all migration files
	migrations, err := getMigrationFiles()
	if err != nil {
		fmt.Printf("Failed to get migration files: %v\n", err)
		os.Exit(1)
	}

	// Rollback migrations
	migrator := db.NewMigrator(database)
	count, err := migrator.RollbackMigrations(migrations, step)
	if err != nil {
		fmt.Printf("Failed to rollback migrations: %v\n", err)
		os.Exit(1)
	}

	if count == 0 {
		fmt.Println("No migrations to rollback")
	} else {
		fmt.Printf("Rolled back %d migrations\n", count)
	}
}

// MigrateReset rolls back all migrations
func MigrateReset(cfg *config.Config) {
	fmt.Println("Resetting migrations...")

	dbConfig := db.DatabaseConfig{
		Driver:   cfg.Database.Driver,
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		Name:     cfg.Database.Name,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
	}

	database, err := db.NewDatabase(&dbConfig)
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	// Get all migration files
	migrations, err := getMigrationFiles()
	if err != nil {
		fmt.Printf("Failed to get migration files: %v\n", err)
		os.Exit(1)
	}

	// Reset migrations
	migrator := db.NewMigrator(database)
	count, err := migrator.ResetMigrations(migrations)
	if err != nil {
		fmt.Printf("Failed to reset migrations: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Reset %d migrations\n", count)
}

// MigrateRefresh resets and reruns all migrations
func MigrateRefresh(cfg *config.Config) {
	fmt.Println("Refreshing migrations...")

	dbConfig := db.DatabaseConfig{
		Driver:   cfg.Database.Driver,
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		Name:     cfg.Database.Name,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
	}

	database, err := db.NewDatabase(&dbConfig)
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	// Get all migration files
	migrations, err := getMigrationFiles()
	if err != nil {
		fmt.Printf("Failed to get migration files: %v\n", err)
		os.Exit(1)
	}

	// Refresh migrations
	migrator := db.NewMigrator(database)
	resetCount, err := migrator.ResetMigrations(migrations)
	if err != nil {
		fmt.Printf("Failed to reset migrations: %v\n", err)
		os.Exit(1)
	}

	_, err = migrator.RunMigrations(migrations)
	if err != nil {
		fmt.Printf("Failed to run migrations: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Refreshed %d migrations\n", resetCount)
}

// MakeMigration creates a new migration file
func MakeMigration(name string) {
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("%s_%s.go", timestamp, name)
	path := filepath.Join("migrations", filename)

	// Create migrations directory if it doesn't exist
	if err := os.MkdirAll("migrations", 0755); err != nil {
		fmt.Printf("Failed to create migrations directory: %v\n", err)
		os.Exit(1)
	}

	// Create migration file
	file, err := os.Create(path)
	if err != nil {
		fmt.Printf("Failed to create migration file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Write migration template
	template := fmt.Sprintf(`package migrations

import (
	"github.com/example/goframe/db"
)

// Migration_%s represents the %s migration
type Migration_%s struct{}

// Up runs the migration
func (m *Migration_%s) Up(migrator *db.Migrator) error {
	// Create table
	sql := ` + "`" + `
		CREATE TABLE IF NOT EXISTS table_name (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	` + "`" + `

	return migrator.DB().Exec(sql)
}

// Down rolls back the migration
func (m *Migration_%s) Down(migrator *db.Migrator) error {
	// Drop table
	sql := "DROP TABLE IF EXISTS table_name"
	return migrator.DB().Exec(sql)
}
`, timestamp, name, timestamp, timestamp, timestamp)

	if _, err := file.WriteString(template); err != nil {
		fmt.Printf("Failed to write migration file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created migration: %s\n", path)
}

// getMigrationFiles returns all migration files
func getMigrationFiles() ([]string, error) {
	// Check if migrations directory exists
	if _, err := os.Stat("migrations"); os.IsNotExist(err) {
		return []string{}, nil
	}

	// Get all migration files
	var files []string
	err := filepath.Walk("migrations", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

