package db

import (
	"errors"
	"fmt"
	"path/filepath"
	// "plugin"
	"sort"
	"strings"
	"time"
)

// Migration represents a database migration
type Migration struct {
	ID        uint      `db:"id"`
	Name      string    `db:"name"`
	Batch     int       `db:"batch"`
	CreatedAt time.Time `db:"created_at"`
}

// MigrationInterface defines the interface for migrations
type MigrationInterface interface {
	Up(migrator *Migrator) error
	Down(migrator *Migrator) error
}

// Migrator handles database migrations
type Migrator struct {
	db *Database
}

// NewMigrator creates a new migrator
func NewMigrator(db *Database) *Migrator {
	return &Migrator{db: db}
}

// DB returns the database connection
func (m *Migrator) DB() *Database {
	return m.db
}

// RunMigrations runs all pending migrations
func (m *Migrator) RunMigrations(files []string) (int, error) {
	// Create migrations table if it doesn't exist
	if err := m.createMigrationsTable(); err != nil {
		return 0, err
	}

	// Get already run migrations
	var migrations []Migration
	if err := m.db.Table(Migration{}).Find(&migrations); err != nil {
		return 0, err
	}

	// Get the current batch number
	batch := 1
	if len(migrations) > 0 {
		// Find the maximum batch number
		for _, migration := range migrations {
			if migration.Batch > batch {
				batch = migration.Batch
			}
		}
		batch++
	}

	// Get pending migrations
	pendingMigrations := m.getPendingMigrations(files, migrations)

	// Sort migrations by name (which includes timestamp)
	sort.Strings(pendingMigrations)

	// Run pending migrations
	count := 0
	for _, file := range pendingMigrations {
		// Load migration
		migration, err := m.loadMigration(file)
		if err != nil {
			return count, err
		}

		// Run migration
		if err := migration.Up(m); err != nil {
			return count, err
		}

		// Record migration
		name := filepath.Base(file)
		if err := m.recordMigration(name, batch); err != nil {
			return count, err
		}

		count++
	}

	return count, nil
}

// RollbackMigrations rolls back the last batch of migrations
func (m *Migrator) RollbackMigrations(files []string, step int) (int, error) {
	// Get already run migrations
	var migrations []Migration
	if err := m.db.Table(Migration{}).OrderBy("batch DESC, id DESC").Find(&migrations); err != nil {
		return 0, err
	}

	if len(migrations) == 0 {
		return 0, nil
	}

	// Group migrations by batch
	batchMigrations := make(map[int][]Migration)
	for _, migration := range migrations {
		batchMigrations[migration.Batch] = append(batchMigrations[migration.Batch], migration)
	}

	// Get batches to rollback
	var batches []int
	for batch := range batchMigrations {
		batches = append(batches, batch)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(batches)))

	// Limit by step if provided
	if step > 0 && step < len(batches) {
		batches = batches[:step]
	}

	// Rollback migrations
	count := 0
	for _, batch := range batches {
		for _, migration := range batchMigrations[batch] {
			// Find migration file
			var migrationFile string
			for _, file := range files {
				if filepath.Base(file) == migration.Name {
					migrationFile = file
					break
				}
			}

			if migrationFile == "" {
				return count, fmt.Errorf("migration file not found: %s", migration.Name)
			}

			// Load migration
			migrationObj, err := m.loadMigration(migrationFile)
			if err != nil {
				return count, err
			}

			// Run down migration
			if err := migrationObj.Down(m); err != nil {
				return count, err
			}

			// Remove migration record
			if err := m.removeMigration(migration.ID); err != nil {
				return count, err
			}

			count++
		}
	}

	return count, nil
}

// ResetMigrations rolls back all migrations
func (m *Migrator) ResetMigrations(files []string) (int, error) {
	// Get already run migrations
	var migrations []Migration
	if err := m.db.Table(Migration{}).OrderBy("batch DESC, id DESC").Find(&migrations); err != nil {
		return 0, err
	}

	if len(migrations) == 0 {
		return 0, nil
	}

	// Rollback all migrations
	count := 0
	for _, migration := range migrations {
		// Find migration file
		var migrationFile string
		for _, file := range files {
			if filepath.Base(file) == migration.Name {
				migrationFile = file
				break
			}
		}

		if migrationFile == "" {
			return count, fmt.Errorf("migration file not found: %s", migration.Name)
		}

		// Load migration
		migrationObj, err := m.loadMigration(migrationFile)
		if err != nil {
			return count, err
		}

		// Run down migration
		if err := migrationObj.Down(m); err != nil {
			return count, err
		}

		// Remove migration record
		if err := m.removeMigration(migration.ID); err != nil {
			return count, err
		}

		count++
	}

	return count, nil
}

// createMigrationsTable creates the migrations table if it doesn't exist
func (m *Migrator) createMigrationsTable() error {
	sql := `
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			batch INTEGER NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`
	return m.db.Exec(sql)
}

// getPendingMigrations returns migrations that haven't been run yet
func (m *Migrator) getPendingMigrations(files []string, migrations []Migration) []string {
	// Create a map of already run migrations
	migrationMap := make(map[string]bool)
	for _, migration := range migrations {
		migrationMap[migration.Name] = true
	}

	// Get pending migrations
	var pendingMigrations []string
	for _, file := range files {
		name := filepath.Base(file)
		if !migrationMap[name] {
			pendingMigrations = append(pendingMigrations, file)
		}
	}

	return pendingMigrations
}

// loadMigration loads a migration from a file
func (m *Migrator) loadMigration(file string) (MigrationInterface, error) {
	// In a real implementation, this would load the migration from the file
	// For simplicity, we'll just return a dummy migration
	
	// This is a simplified version - in a real implementation, 
	// you would use Go's plugin system or reflection to load the migration
	
	// Extract migration name from file
	name := strings.TrimSuffix(filepath.Base(file), ".go")
	parts := strings.SplitN(name, "_", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid migration filename format")
	}
	
	// In a real implementation, you would load the migration class
	// For now, we'll return a dummy migration
	return &dummyMigration{name: parts[1]}, nil
}

// dummyMigration is a placeholder for a real migration
type dummyMigration struct {
	name string
}

// Up runs the migration
func (m *dummyMigration) Up(migrator *Migrator) error {
	fmt.Printf("Running migration: %s\n", m.name)
	return nil
}

// Down rolls back the migration
func (m *dummyMigration) Down(migrator *Migrator) error {
	fmt.Printf("Rolling back migration: %s\n", m.name)
	return nil
}

// recordMigration records a migration in the migrations table
func (m *Migrator) recordMigration(name string, batch int) error {
	sql := `
		INSERT INTO migrations (name, batch, created_at)
		VALUES (?, ?, ?)
	`
	return m.db.Exec(sql, name, batch, time.Now())
}

// removeMigration removes a migration from the migrations table
func (m *Migrator) removeMigration(id uint) error {
	sql := "DELETE FROM migrations WHERE id = ?"
	return m.db.Exec(sql, id)
}

