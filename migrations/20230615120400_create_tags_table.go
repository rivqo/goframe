package migrations

import (
	"github.com/example/goframe/db"
)

// Migration_20230615120400 represents the create_tags_table migration
type Migration_20230615120400 struct{}

// Up runs the migration
func (m *Migration_20230615120400) Up(migrator *db.Migrator) error {
	// Create tags table
	tagsSql := `
	CREATE TABLE IF NOT EXISTS tags (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL UNIQUE,
		slug VARCHAR(255) NOT NULL UNIQUE,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	)
	`

	// Create post_tag pivot table
	pivotSql := `
	CREATE TABLE IF NOT EXISTS post_tag (
		post_id INTEGER NOT NULL,
		tag_id INTEGER NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		PRIMARY KEY (post_id, tag_id),
		FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
		FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
	)
	`

	// Create indexes
	indexSql := `
	CREATE INDEX idx_tags_slug ON tags(slug);
	CREATE INDEX idx_post_tag_post_id ON post_tag(post_id);
	CREATE INDEX idx_post_tag_tag_id ON post_tag(tag_id);
	`

	if err := migrator.DB().Exec(tagsSql); err != nil {
		return err
	}

	if err := migrator.DB().Exec(pivotSql); err != nil {
		return err
	}

	return migrator.DB().Exec(indexSql)
}

// Down rolls back the migration
func (m *Migration_20230615120400) Down(migrator *db.Migrator) error {
	// Drop indexes
	indexSql := `
	DROP INDEX IF EXISTS idx_tags_slug;
	DROP INDEX IF EXISTS idx_post_tag_post_id;
	DROP INDEX IF EXISTS idx_post_tag_tag_id;
	`

	if err := migrator.DB().Exec(indexSql); err != nil {
		return err
	}

	// Drop pivot table
	if err := migrator.DB().Exec("DROP TABLE IF EXISTS post_tag"); err != nil {
		return err
	}

	// Drop tags table
	return migrator.DB().Exec("DROP TABLE IF EXISTS tags")
}

