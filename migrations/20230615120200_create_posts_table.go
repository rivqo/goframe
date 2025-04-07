package migrations

import (
	"github.com/example/goframe/db"
)

// Migration_20230615120200 represents the create_posts_table migration
type Migration_20230615120200 struct{}

// Up runs the migration
func (m *Migration_20230615120200) Up(migrator *db.Migrator) error {
	// Create table
	sql := `
	CREATE TABLE IF NOT EXISTS posts (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		slug VARCHAR(255) NOT NULL UNIQUE,
		content TEXT NOT NULL,
		excerpt TEXT NULL,
		user_id INTEGER NOT NULL,
		published BOOLEAN NOT NULL DEFAULT FALSE,
		published_at TIMESTAMP NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	)
	`

	// Create indexes
	indexSql := `
	CREATE INDEX idx_posts_user_id ON posts(user_id);
	CREATE INDEX idx_posts_slug ON posts(slug);
	CREATE INDEX idx_posts_published ON posts(published);
	`

	if err := migrator.DB().Exec(sql); err != nil {
		return err
	}

	return migrator.DB().Exec(indexSql)
}

// Down rolls back the migration
func (m *Migration_20230615120200) Down(migrator *db.Migrator) error {
	// Drop indexes
	indexSql := `
	DROP INDEX IF EXISTS idx_posts_user_id;
	DROP INDEX IF EXISTS idx_posts_slug;
	DROP INDEX IF EXISTS idx_posts_published;
	`

	if err := migrator.DB().Exec(indexSql); err != nil {
		return err
	}

	// Drop table
	sql := "DROP TABLE IF EXISTS posts"
	return migrator.DB().Exec(sql)
}

