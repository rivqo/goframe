import { DocPagination } from "@/components/doc-pagination"

export default function QueryBuilderPage() {
  return (
    <div className="space-y-6">
      <h1>Query Builder</h1>
      <p>
        GoFrame provides a fluent query builder that makes it easy to build SQL queries. The query builder offers a
        convenient, fluent interface for creating and executing database queries.
      </p>

      <h2>Basic Usage</h2>
      <p>
        To start building a query, use the <code>Table</code> method on the database instance:
      </p>

      <pre>
        <code>{`// Get all users
var users []models.User
err := db.Table("users").Get(&users)

// Get a specific user
var user models.User
err := db.Table("users").Where("id", "=", 1).First(&user)`}</code>
      </pre>

      <h2>Retrieving Results</h2>

      <h3>Retrieving All Rows</h3>
      <pre>
        <code>{`// Get all users
var users []models.User
err := db.Table("users").Get(&users)

// With conditions
var activeUsers []models.User
err := db.Table("users").Where("active", "=", true).Get(&activeUsers)`}</code>
      </pre>

      <h3>Retrieving A Single Row</h3>
      <pre>
        <code>{`// Get first user
var user models.User
err := db.Table("users").First(&user)

// Get user by ID
var user models.User
err := db.Table("users").Where("id", "=", 1).First(&user)`}</code>
      </pre>

      <h3>Retrieving A Single Column</h3>
      <pre>
        <code>{`// Get all user names
var names []string
err := db.Table("users").Select("name").Get(&names)

// Get a single value
var name string
err := db.Table("users").Where("id", "=", 1).Select("name").First(&name)`}</code>
      </pre>

      <h3>Aggregates</h3>
      <pre>
        <code>{`// Count users
count, err := db.Table("users").Count()

// Get max value
var maxAge int
err := db.Table("users").Select("MAX(age) as max_age").First(&maxAge)`}</code>
      </pre>

      <h2>Building Queries</h2>

      <h3>Select Statements</h3>
      <pre>
        <code>{`// Select specific columns
var users []models.User
err := db.Table("users").Select("id", "name", "email").Get(&users)

// Select with alias
var result struct {
    ID int \`db:"id"\`
    FullName string \`db:"full_name"\`
}
err := db.Table("users").Select("id", "name as full_name").First(&result)`}</code>
      </pre>

      <h3>Where Clauses</h3>
      <pre>
        <code>{`// Basic where
var users []models.User
err := db.Table("users").Where("active", "=", true).Get(&users)

// Multiple conditions
var users []models.User
err := db.Table("users").
    Where("active", "=", true).
    Where("age", ">", 18).
    Get(&users)

// OR where
var users []models.User
err := db.Table("users").
    Where("active", "=", true).
    OrWhere("role", "=", "admin").
    Get(&users)

// Where IN
var users []models.User
err := db.Table("users").
    WhereIn("id", 1, 2, 3).
    Get(&users)

// Where NULL
var users []models.User
err := db.Table("users").
    WhereNull("deleted_at").
    Get(&users)

// Raw where
var users []models.User
err := db.Table("users").
    WhereRaw("age > ? AND role = ?", 18, "user").
    Get(&users)`}</code>
      </pre>

      <h3>Ordering, Grouping, Limit & Offset</h3>
      <pre>
        <code>{`// Ordering
var users []models.User
err := db.Table("users").
    OrderBy("name", "asc").
    OrderBy("created_at", "desc").
    Get(&users)

// Grouping
var results []struct {
    Role string \`db:"role"\`
    Count int \`db:"count"\`
}
err := db.Table("users").
    Select("role", "COUNT(*) as count").
    GroupBy("role").
    Get(&results)

// Limit & Offset
var users []models.User
err := db.Table("users").
    Limit(10).
    Offset(20).
    Get(&users)`}</code>
      </pre>

      <h3>Joins</h3>
      <pre>
        <code>{`// Inner join
var results []struct {
    UserName string \`db:"user_name"\`
    PostTitle string \`db:"post_title"\`
}
err := db.Table("users").
    Select("users.name as user_name", "posts.title as post_title").
    Join("posts", "users.id", "=", "posts.user_id").
    Get(&results)

// Left join
var results []struct {
    UserName string \`db:"user_name"\`
    PostTitle string \`db:"post_title"\`
}
err := db.Table("users").
    Select("users.name as user_name", "posts.title as post_title").
    LeftJoin("posts", "users.id", "=", "posts.user_id").
    Get(&results)`}</code>
      </pre>

      <h2>Inserts</h2>
      <pre>
        <code>{`// Insert a single record
id, err := db.Table("users").Insert(map[string]interface{}{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password",
})

// Insert multiple records
ids, err := db.Table("users").Insert([]map[string]interface{}{
    {
        "name": "John Doe",
        "email": "john@example.com",
        "password": "password",
    },
    {
        "name": "Jane Smith",
        "email": "jane@example.com",
        "password": "password",
    },
})`}</code>
      </pre>

      <h2>Updates</h2>
      <pre>
        <code>{`// Update a single record
affected, err := db.Table("users").
    Where("id", "=", 1).
    Update(map[string]interface{}{
        "name": "John Smith",
        "updated_at": time.Now(),
    })

// Update multiple records
affected, err := db.Table("users").
    Where("active", "=", false).
    Update(map[string]interface{}{
        "deleted_at": time.Now(),
    })`}</code>
      </pre>

      <h2>Deletes</h2>
      <pre>
        <code>{`// Delete a single record
affected, err := db.Table("users").
    Where("id", "=", 1).
    Delete()

// Delete multiple records
affected, err := db.Table("users").
    Where("active", "=", false).
    Delete()`}</code>
      </pre>

      <h2>Transactions</h2>
      <p>GoFrame supports database transactions to ensure data integrity:</p>

      <pre>
        <code>{`// Start a transaction
tx, err := db.Begin()
if err != nil {
    return err
}

// Defer rollback in case of error
defer tx.Rollback()

// Perform operations within the transaction
_, err = tx.Table("users").Insert(map[string]interface{}{
    "name": "John Doe",
    "email": "john@example.com",
})
if err != nil {
    return err
}

_, err = tx.Table("profiles").Insert(map[string]interface{}{
    "user_id": id,
    "bio": "A software developer",
})
if err != nil {
    return err
}

// Commit the transaction
return tx.Commit()`}</code>
      </pre>

      <h2>Raw Queries</h2>
      <p>Sometimes you may need to use raw SQL queries. GoFrame provides methods for executing raw queries:</p>

      <pre>
        <code>{`// Raw select
var users []models.User
err := db.Query("SELECT * FROM users WHERE active = ?", true).Scan(&users)

// Raw execute
result, err := db.Exec("UPDATE users SET active = ? WHERE id = ?", true, 1)`}</code>
      </pre>

      <DocPagination
        prev={{
          title: "Migrations",
          href: "/docs/database/migrations",
        }}
        next={{
          title: "Repositories",
          href: "/docs/database/repositories",
        }}
      />
    </div>
  )
}

