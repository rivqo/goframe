package migrations

import (
	"github.com/example/goframe/db"
)

// Migration_20230615120500 represents the create_settings_table migration
type Migration_20230615120500 struct{}

// Up runs the migration
func (m *Migration_20230615120500) Up(migrator *db.Migrator) error {
	// Create table
	sql := `
	CREATE TABLE IF NOT EXISTS settings (
		id SERIAL PRIMARY KEY,
		key VARCHAR(255) NOT NULL UNIQUE,
		value TEXT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	)
	`

	// Create index
	indexSql := `
	CREATE INDEX idx_settings_key ON settings(key);
	`

	// Insert default settings
	defaultsSql := `
	INSERT INTO settings (key, value, created_at, updated_at) VALUES
	('site_name', 'GoFrame Blog', NOW(), NOW()),
	('site_description', 'A blog built with GoFrame', NOW(), NOW()),
	('site_logo', '/assets/images/logo.png', NOW(), NOW()),
	('site_favicon', '/assets/images/favicon.ico', NOW(), NOW()),
	('site_email', 'admin@example.com', NOW(), NOW()),
	('posts_per_page', '10', NOW(), NOW()),
	('comments_enabled', 'true', NOW(), NOW()),
	('registration_enabled', 'true', NOW(), NOW()),
	('maintenance_mode', 'false', NOW(), NOW()),
	('theme', 'default', NOW(), NOW())
	ON CONFLICT (key) DO NOTHING;
	`

	if err := migrator.DB().Exec(sql); err != nil {
		return err
	}

	if err := migrator.DB().Exec(indexSql); err != nil {
		return err
	}

	return migrator.DB().Exec(defaultsSql)
}

// Down rolls back the migration
func (m *Migration_20230615120500) Down(migrator *db.Migrator) error {
	// Drop index
	if err := migrator.DB().Exec("DROP INDEX IF EXISTS idx_settings_key"); err != nil {
		return err
	}

	// Drop table
	return migrator.DB().Exec("DROP TABLE IF EXISTS settings")
}

