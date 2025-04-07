package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     int
	Name     string
	User     string
	Password string
}

type Database struct {
	config DatabaseConfig
	db     *sql.DB
}

// Connect creates a new database connection
func (db *Database) Connect(config DatabaseConfig) (*Database, error) {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Name)

	db.db, err = sql.Open(config.Driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Test the connection
	err = db.db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	db.config = config
	return db, nil
}

func NewDatabase(cfg *DatabaseConfig) (*Database, error) {
	// Create connection string based on driver
	var dsn string
	switch cfg.Driver {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	case "postgres":
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	// Open database connection
	sqlDB, err := sql.Open(cfg.Driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	// sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	// sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	// sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Verify connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Database{db: sqlDB}, nil
}

// Close closes the database connection
func (db *Database) Close() error {
	if db.db != nil {
		return db.db.Close()
	}
	return nil
}

// Model represents a database model
type Model struct {
	ID uint `db:"id" json:"id"`
}

// Create inserts a new record into the database
func (db *Database) Create(value interface{}) error {
	// Get table name from type
	tableName := strings.ToLower(reflect.TypeOf(value).Name())

	// Get fields and values using reflection
	val := reflect.ValueOf(value).Elem()
	typ := val.Type()

	var fields []string
	var placeholders []string
	var values []interface{}

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag == "" || dbTag == "-" {
			continue
		}

		fields = append(fields, dbTag)
		placeholders = append(placeholders, "?")
		values = append(values, val.Field(i).Interface())
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(fields, ", "),
		strings.Join(placeholders, ", "),
	)

	_, err := db.db.Exec(query, values...)
	return err
}

// Query represents a database query
type Query struct {
	db         *Database
	model      interface{}
	conditions []string
	values     []interface{}
	limit      int
	offset     int
	orderBy    string
}

// Table creates a new query for the given model
func (db *Database) Table(model interface{}) *Query {
	return &Query{
		db:    db,
		model: model,
	}
}

// Where adds a where condition to the query
func (q *Query) Where(condition string, values ...interface{}) *Query {
	q.conditions = append(q.conditions, condition)
	q.values = append(q.values, values...)
	return q
}

// Limit sets the limit for the query
func (q *Query) Limit(limit int) *Query {
	q.limit = limit
	return q
}

// Offset sets the offset for the query
func (q *Query) Offset(offset int) *Query {
	q.offset = offset
	return q
}

// OrderBy sets the order by clause for the query
func (q *Query) OrderBy(orderBy string) *Query {
	q.orderBy = orderBy
	return q
}

// Find executes the query and scans the result into dest
func (q *Query) Find(dest interface{}) error {
	tableName := strings.ToLower(reflect.TypeOf(q.model).Name())
	query := fmt.Sprintf("SELECT * FROM %s", tableName)

	if len(q.conditions) > 0 {
		query += " WHERE " + strings.Join(q.conditions, " AND ")
	}

	if q.orderBy != "" {
		query += " ORDER BY " + q.orderBy
	}

	if q.limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", q.limit)
	}

	if q.offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", q.offset)
	}

	rows, err := q.db.db.Query(query, q.values...)
	if err != nil {
		return err
	}
	defer rows.Close()

	return rowsScan(rows, dest)
}

// First executes the query and scans the first result into dest
func (q *Query) First(dest interface{}) error {
	q.Limit(1)
	return q.Find(dest)
}

// Update updates the given model with the given values
func (q *Query) Update(values map[string]interface{}) error {
	if len(values) == 0 {
		return errors.New("no values provided for update")
	}

	tableName := strings.ToLower(reflect.TypeOf(q.model).Name())
	var setClauses []string
	var updateValues []interface{}

	for field, value := range values {
		setClauses = append(setClauses, fmt.Sprintf("%s = ?", field))
		updateValues = append(updateValues, value)
	}

	query := fmt.Sprintf("UPDATE %s SET %s", tableName, strings.Join(setClauses, ", "))

	if len(q.conditions) > 0 {
		query += " WHERE " + strings.Join(q.conditions, " AND ")
		updateValues = append(updateValues, q.values...)
	}

	_, err := q.db.db.Exec(query, updateValues...)
	return err
}

// Delete deletes records matching the query
func (q *Query) Delete() error {
	tableName := strings.ToLower(reflect.TypeOf(q.model).Name())
	query := fmt.Sprintf("DELETE FROM %s", tableName)

	if len(q.conditions) > 0 {
		query += " WHERE " + strings.Join(q.conditions, " AND ")
	}

	_, err := q.db.db.Exec(query, q.values...)
	return err
}

// Transaction represents a database transaction
type Transaction struct {
	db *Database
	tx *sql.Tx
}

// Begin starts a new transaction
func (db *Database) Begin() (*Transaction, error) {
	tx, err := db.db.Begin()
	if err != nil {
		return nil, err
	}
	return &Transaction{db: db, tx: tx}, nil
}

// Commit commits the transaction
func (tx *Transaction) Commit() error {
	return tx.tx.Commit()
}

// Rollback rolls back the transaction
func (tx *Transaction) Rollback() error {
	return tx.tx.Rollback()
}

// Exec executes a SQL query and returns only error
func (db *Database) Exec(query string, args ...interface{}) error {
	_, err := db.db.Exec(query, args...)
	return err
}

// Query executes a SQL query that returns rows
func (db *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.db.Query(query, args...)
}

// QueryRow executes a SQL query that returns at most one row
func (db *Database) QueryRow(query string, args ...interface{}) *sql.Row {
	return db.db.QueryRow(query, args...)
}

// scanRows scans database rows into a destination slice
func (db *Database) scanRows(rows *sql.Rows, dest interface{}) error {
	destVal := reflect.ValueOf(dest)
	if destVal.Kind() != reflect.Ptr || destVal.Elem().Kind() != reflect.Slice {
		return errors.New("dest must be a pointer to a slice")
	}

	sliceVal := destVal.Elem()
	elemType := sliceVal.Type().Elem()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	for rows.Next() {
		elem := reflect.New(elemType).Elem()
		fields := make([]interface{}, len(columns))

		for i := 0; i < elem.NumField(); i++ {
			field := elem.Type().Field(i)
			dbTag := field.Tag.Get("db")
			if dbTag == "" {
				continue
			}

			for j, col := range columns {
				if strings.EqualFold(col, dbTag) {
					fields[j] = elem.Field(i).Addr().Interface()
					break
				}
			}
		}

		if err := rows.Scan(fields...); err != nil {
			return err
		}

		sliceVal.Set(reflect.Append(sliceVal, elem))
	}

	return rows.Err()
}

// scanRows scans database rows into a destination slice
func rowsScan(rows *sql.Rows, dest interface{}) error {
	destVal := reflect.ValueOf(dest)
	if destVal.Kind() != reflect.Ptr || destVal.Elem().Kind() != reflect.Slice {
		return errors.New("dest must be a pointer to a slice")
	}

	sliceVal := destVal.Elem()
	elemType := sliceVal.Type().Elem()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	for rows.Next() {
		elem := reflect.New(elemType).Elem()
		fields := make([]interface{}, len(columns))

		for i := 0; i < elem.NumField(); i++ {
			field := elem.Type().Field(i)
			dbTag := field.Tag.Get("db")
			if dbTag == "" {
				continue
			}

			for j, col := range columns {
				if strings.EqualFold(col, dbTag) {
					fields[j] = elem.Field(i).Addr().Interface()
					break
				}
			}
		}

		if err := rows.Scan(fields...); err != nil {
			return err
		}

		sliceVal.Set(reflect.Append(sliceVal, elem))
	}

	return rows.Err()
}