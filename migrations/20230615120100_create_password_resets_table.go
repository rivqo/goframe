package migrations

import (
	"github.com/example/goframe/db"
)

// Migration_20230615120100 represents the create_password_resets_table migration
type Migration_20230615120100 struct{}

// Up runs the migration
func (m *Migration_20230615120100) Up(migrator *db.Migrator) error {
	// Create table
	sql := `
	CREATE TABLE IF NOT EXISTS password_resets (
		email VARCHAR(255) NOT NULL,
		token VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		PRIMARY KEY (email, token)
	)
	`

	return migrator.DB().Exec(sql)
}

// Down rolls back the migration
func (m *Migration_20230615120100) Down(migrator *db.Migrator) error {
	// Drop table
	sql := "DROP TABLE IF EXISTS password_resets"
	return migrator.DB().Exec(sql)
}

