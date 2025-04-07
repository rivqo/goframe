import { DocPagination } from "@/components/doc-pagination"

export default function MigrationsPage() {
  return (
    <div className="space-y-6">
      <h1>Migrations</h1>
      <p>
        Migrations are like version control for your database, allowing your team to define and share the application's
        database schema definition. GoFrame provides a robust migration system inspired by Laravel's migrations.
      </p>

      <h2>Creating Migrations</h2>
      <p>
        To create a new migration, use the <code>make:migration</code> command:
      </p>

      <pre>
        <code>{`goframe make:migration create_users_table`}</code>
      </pre>

      <p>
        This will create a new migration file in the <code>migrations</code> directory with a timestamp prefix:
      </p>

      <pre>
        <code>{`// migrations/20230615120000_create_users_table.go
package migrations

import (
"github.com/example/goframe/db"
)

// Migration_20230615120000 represents the create_users_table migration
type Migration_20230615120000 struct{}

// Up runs the migration
func (m *Migration_20230615120000) Up(migrator *db.Migrator) error {
// Create table
sql := \`
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	)
\`

return migrator.DB().Exec(sql)
}

// Down rolls back the migration
func (m *Migration_20230615120000) Down(migrator *db.Migrator) error {
// Drop table
sql := "DROP TABLE IF EXISTS users"
return migrator.DB().Exec(sql)
}
`}</code>
      </pre>

      <h2>Using the Schema Builder</h2>
      <p>
        GoFrame now includes a fluent schema builder that makes it easier to define your database schema. Here's how to
        use it:
      </p>

      <pre>
        <code>{`// migrations/20230615120000_create_users_table.go
package migrations

import (
"github.com/example/goframe/db"
)

// Migration_20230615120000 represents the create_users_table migration
type Migration_20230615120000 struct{}

// Up runs the migration
func (m *Migration_20230615120000) Up(migrator *db.Migrator) error {
    schema := db.NewSchema(migrator.DB())
    
    return schema.Create("users", func(table *db.Blueprint) {
        table.ID()
        table.String("name", 255)
        table.String("email", 255).Unique()
        table.String("password", 255)
        table.Boolean("active").Default(true)
        table.DateTime("email_verified_at").Nullable()
        table.Timestamps()
    })
}

// Down rolls back the migration
func (m *Migration_20230615120000) Down(migrator *db.Migrator) error {
    schema := db.NewSchema(migrator.DB())
    return schema.Drop("users")
}
`}</code>
      </pre>

      <h2>Running Migrations</h2>
      <p>
        To run all pending migrations, use the <code>migrate</code> command:
      </p>

      <pre>
        <code>{`goframe migrate`}</code>
      </pre>

      <h2>Rolling Back Migrations</h2>
      <p>GoFrame provides several commands for rolling back migrations:</p>
      <pre>
        <code>{`# Rollback the last batch of migrations
goframe migrate --rollback

# Rollback all migrations
goframe migrate --reset

# Rollback and re-run all migrations
goframe migrate --refresh

# Rollback a specific number of migrations
goframe migrate --step=2`}</code>
      </pre>

      <h2>Migration Structure</h2>
      <p>Each migration file contains two methods:</p>
      <ul>
        <li>
          <code>Up</code>: Performs the migration (create tables, add columns, etc.)
        </li>
        <li>
          <code>Down</code>: Reverses the migration (drop tables, remove columns, etc.)
        </li>
      </ul>
      <p>
        Migrations are tracked in a <code>migrations</code> table in your database, which stores the name of each
        migration file that has been run.
      </p>

      <h2>Schema Builder Methods</h2>
      <p>The schema builder provides a variety of methods for defining your database schema:</p>

      <h3>Table Operations</h3>
      <pre>
        <code>{`// Create a new table
schema.Create("users", func(table *db.Blueprint) {
    // Define columns
})

// Modify an existing table
schema.Table("users", func(table *db.Blueprint) {
    // Add or modify columns
})

// Drop a table
schema.Drop("users")

// Rename a table
schema.Rename("users", "people")`}</code>
      </pre>

      <h3>Column Types</h3>
      <pre>
        <code>{`// ID column (auto-incrementing primary key)
table.ID()

// String column
table.String("name", 255)

// Text column
table.Text("description")

// Integer column
table.Integer("age", false) // Not auto-incrementing

// Big integer column
table.BigInteger("views", false)

// Boolean column
table.Boolean("active")

// Date column
table.Date("birth_date")

// DateTime column
table.DateTime("published_at")

// Decimal column
table.Decimal("price", 8, 2) // 8 digits, 2 decimal places

// Float column
table.Float("rating")

// JSON column
table.JSON("metadata")

// Timestamps (created_at and updated_at)
table.Timestamps()

// Soft deletes (deleted_at)
table.SoftDeletes()`}</code>
      </pre>

      <h3>Column Modifiers</h3>
      <pre>
        <code>{`// Nullable column
table.String("middle_name", 255).Nullable()

// Default value
table.Boolean("active").Default(true)

// Unsigned integer
table.Integer("age", false).Unsigned()

// Unique column
table.String("email", 255).Unique()

// Add an index
table.String("name", 255).Index()

// Add a comment
table.String("name", 255).Comment("User's full name")

// Position a column after another column
table.String("middle_name", 255).After("first_name")

// Position a column first
table.String("id", 36).First()

// Make a column the primary key
table.String("uuid", 36).Primary()`}</code>
      </pre>

      <h3>Indexes</h3>
      <pre>
        <code>{`// Add an index
table.Index("name")

// Add a unique index
table.Unique("email")

// Add a composite index
table.Index("first_name", "last_name")

// Set the primary key
table.Primary("id")`}</code>
      </pre>

      <h3>Table Options</h3>
      <pre>
        <code>{`// Set the storage engine
table.Engine("InnoDB")

// Set the character set
table.Charset("utf8mb4")

// Set the collation
table.Collation("utf8mb4_unicode_ci")

// Make the table temporary
table.Temporary()`}</code>
      </pre>

      <h2>Common Migration Examples</h2>

      <h3>Creating a Table</h3>
      <pre>
        <code>{`func (m *Migration_20230615120000) Up(migrator *db.Migrator) error {
    schema := db.NewSchema(migrator.DB())
    
    return schema.Create("posts", func(table *db.Blueprint) {
        table.ID()
        table.String("title", 255)
        table.String("slug", 255).Unique()
        table.Text("content")
        table.BigInteger("user_id", false)
        table.Boolean("published").Default(false)
        table.DateTime("published_at").Nullable()
        table.Timestamps()
        
        table.Index("user_id")
        table.Index("published")
    })
}`}</code>
      </pre>

      <h3>Adding a Column</h3>
      <pre>
        <code>{`func (m *Migration_20230615120100) Up(migrator *db.Migrator) error {
    schema := db.NewSchema(migrator.DB())
    
    return schema.Table("users", func(table *db.Blueprint) {
        table.String("phone", 20).Nullable().After("email")
    })
}`}</code>
      </pre>

      <h3>Creating a Pivot Table</h3>
      <pre>
        <code>{`func (m *Migration_20230615120200) Up(migrator *db.Migrator) error {
    schema := db.NewSchema(migrator.DB())
    
    return schema.Create("post_tag", func(table *db.Blueprint) {
        table.BigInteger("post_id", false)
        table.BigInteger("tag_id", false)
        table.DateTime("created_at").Nullable()
        
        table.Primary("post_id", "tag_id")
        table.Index("post_id")
        table.Index("tag_id")
    })
}`}</code>
      </pre>

      <DocPagination
        prev={{
          title: "Database",
          href: "/docs/database",
        }}
        next={{
          title: "Query Builder",
          href: "/docs/database/query-builder",
        }}
      />
    </div>
  )
}

