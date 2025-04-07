package migrations

import (
	"github.com/example/goframe/db"
)

// Migration_20230615120300 represents the create_comments_table migration
type Migration_20230615120300 struct{}

// Up runs the migration
func (m *Migration_20230615120300) Up(migrator *db.Migrator) error {
	// Create table
	sql := `
	CREATE TABLE IF NOT EXISTS comments (
		id SERIAL PRIMARY KEY,
		content TEXT NOT NULL,
		user_id INTEGER NOT NULL,
		post_id INTEGER NOT NULL,
		parent_id INTEGER NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
		FOREIGN KEY (parent_id) REFERENCES comments(id) ON DELETE CASCADE
	)
	`

	// Create indexes
	indexSql := `
	CREATE INDEX idx_comments_user_id ON comments(user_id);
	CREATE INDEX idx_comments_post_id ON comments(post_id);
	CREATE INDEX idx_comments_parent_id ON comments(parent_id);
	`

	if err := migrator.DB().Exec(sql); err != nil {
		return err
	}

	return migrator.DB().Exec(indexSql)
}

// Down rolls back the migration
func (m *Migration_20230615120300) Down(migrator *db.Migrator) error {
	// Drop indexes
	indexSql := `
	DROP INDEX IF EXISTS idx_comments_user_id;
	DROP INDEX IF EXISTS idx_comments_post_id;
	DROP INDEX IF EXISTS idx_comments_parent_id;
	`

	if err := migrator.DB().Exec(indexSql); err != nil {
		return err
	}

	// Drop table
	sql := "DROP TABLE IF EXISTS comments"
	return migrator.DB().Exec(sql)
}

