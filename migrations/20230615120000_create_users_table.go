package migrations

import (
	"github.com/example/goframe/db"
)

// Migration_20230615120000 represents the create_users_table migration
type Migration_20230615120000 struct{}

// Up runs the migration
func (m *Migration_20230615120000) Up(migrator *db.Migrator) error {
	// Create table
	sql := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		remember_token VARCHAR(100) NULL,
		email_verified_at TIMESTAMP NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	)
	`

	return migrator.DB().Exec(sql)
}

// Down rolls back the migration
func (m *Migration_20230615120000) Down(migrator *db.Migrator) error {
	// Drop table
	sql := "DROP TABLE IF EXISTS users"
	return migrator.DB().Exec(sql)
}

