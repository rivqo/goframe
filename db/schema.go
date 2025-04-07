package db

import (
	"fmt"
	"strings"
)

// Schema provides methods for creating and modifying database tables
type Schema struct {
	db *Database
}

// NewSchema creates a new schema builder
func NewSchema(db *Database) *Schema {
	return &Schema{db: db}
}

// Blueprint defines the structure of a table
type Blueprint struct {
	table     string
	columns   []Column
	indexes   []Index
	primary   string
	engine    string
	charset   string
	collation string
	temporary bool
}

// Column represents a table column
type Column struct {
	name       string
	columnType string // Changed from 'type' to 'columnType' since 'type' is a reserved word
	length     int
	nullable   bool
	default_   interface{}
	unsigned   bool
	unique     bool
	index      bool
	comment    string
	after      string
	first      bool
}

// Index represents a table index
type Index struct {
	name    string
	columns []string
	unique  bool
}

// Create creates a new table
func (s *Schema) Create(table string, callback func(*Blueprint)) error {
	blueprint := &Blueprint{
		table:     table,
		columns:   []Column{},
		indexes:   []Index{},
		engine:    "InnoDB",
		charset:   "utf8mb4",
		collation: "utf8mb4_unicode_ci",
		temporary: false,
	}
	
	callback(blueprint)
	
	return s.createTable(blueprint)
}

// Table modifies an existing table
func (s *Schema) Table(table string, callback func(*Blueprint)) error {
	blueprint := &Blueprint{
		table:   table,
		columns: []Column{},
		indexes: []Index{},
	}
	
	callback(blueprint)
	
	return s.alterTable(blueprint)
}

// Drop drops a table
func (s *Schema) Drop(table string) error {
	sql := fmt.Sprintf("DROP TABLE IF EXISTS `%s`", table)
	return s.db.Exec(sql)
}

// DropIfExists drops a table if it exists
func (s *Schema) DropIfExists(table string) error {
	return s.Drop(table)
}

// Rename renames a table
func (s *Schema) Rename(from, to string) error {
	sql := fmt.Sprintf("ALTER TABLE `%s` RENAME TO `%s`", from, to)
	return s.db.Exec(sql)
}

// createTable creates a new table
func (s *Schema) createTable(blueprint *Blueprint) error {
	var sql strings.Builder
	
	if blueprint.temporary {
		sql.WriteString("CREATE TEMPORARY TABLE ")
	} else {
		sql.WriteString("CREATE TABLE ")
	}
	
	sql.WriteString(fmt.Sprintf("`%s` (", blueprint.table))
	
	// Add columns
	var columnDefs []string
	for _, column := range blueprint.columns {
		columnDefs = append(columnDefs, s.getColumnDefinition(column))
	}
	
	// Add primary key
	if blueprint.primary != "" {
		columnDefs = append(columnDefs, fmt.Sprintf("PRIMARY KEY (`%s`)", blueprint.primary))
	}
	
	// Add indexes
	for _, index := range blueprint.indexes {
		cols := make([]string, len(index.columns))
		for i, col := range index.columns {
			cols[i] = fmt.Sprintf("`%s`", col)
		}
		if index.unique {
			columnDefs = append(columnDefs, fmt.Sprintf("UNIQUE INDEX `%s` (%s)", index.name, strings.Join(cols, ", ")))
		} else {
			columnDefs = append(columnDefs, fmt.Sprintf("INDEX `%s` (%s)", index.name, strings.Join(cols, ", ")))
		}
	}
	
	sql.WriteString(strings.Join(columnDefs, ", "))
	sql.WriteString(")")
	
	// Add engine
	if blueprint.engine != "" {
		sql.WriteString(fmt.Sprintf(" ENGINE=%s", blueprint.engine))
	}
	
	// Add charset
	if blueprint.charset != "" {
		sql.WriteString(fmt.Sprintf(" DEFAULT CHARSET=%s", blueprint.charset))
	}
	
	// Add collation
	if blueprint.collation != "" {
		sql.WriteString(fmt.Sprintf(" COLLATE=%s", blueprint.collation))
	}
	
	return s.db.Exec(sql.String())
}

// alterTable alters an existing table
func (s *Schema) alterTable(blueprint *Blueprint) error {
	var sql strings.Builder
	
	sql.WriteString("ALTER TABLE ")
	sql.WriteString(fmt.Sprintf("`%s`", blueprint.table))
	
	var alterations []string
	
	// Add columns
	for _, column := range blueprint.columns {
		var alteration string
		
		if column.after != "" {
			alteration = fmt.Sprintf("ADD COLUMN %s AFTER `%s`", s.getColumnDefinition(column), column.after)
		} else if column.first {
			alteration = fmt.Sprintf("ADD COLUMN %s FIRST", s.getColumnDefinition(column))
		} else {
			alteration = fmt.Sprintf("ADD COLUMN %s", s.getColumnDefinition(column))
		}
		
		alterations = append(alterations, alteration)
	}
	
	// Add indexes
	for _, index := range blueprint.indexes {
		var alteration string
		cols := make([]string, len(index.columns))
		for i, col := range index.columns {
			cols[i] = fmt.Sprintf("`%s`", col)
		}
		
		if index.unique {
			alteration = fmt.Sprintf("ADD UNIQUE INDEX `%s` (%s)", index.name, strings.Join(cols, ", "))
		} else {
			alteration = fmt.Sprintf("ADD INDEX `%s` (%s)", index.name, strings.Join(cols, ", "))
		}
		
		alterations = append(alterations, alteration)
	}
	
	sql.WriteString(" ")
	sql.WriteString(strings.Join(alterations, ", "))
	
	return s.db.Exec(sql.String())
}

// getColumnDefinition returns the SQL definition for a column
func (s *Schema) getColumnDefinition(column Column) string {
	var def strings.Builder
	
	def.WriteString(fmt.Sprintf("`%s` ", column.name))
	def.WriteString(column.columnType)
	
	if column.length > 0 {
		def.WriteString(fmt.Sprintf("(%d)", column.length))
	}
	
	if column.unsigned {
		def.WriteString(" UNSIGNED")
	}
	
	if !column.nullable {
		def.WriteString(" NOT NULL")
	} else {
		def.WriteString(" NULL")
	}
	
	if column.default_ != nil {
		switch v := column.default_.(type) {
		case string:
			def.WriteString(fmt.Sprintf(" DEFAULT '%s'", v))
		default:
			def.WriteString(fmt.Sprintf(" DEFAULT %v", v))
		}
	}
	
	if column.unique {
		def.WriteString(" UNIQUE")
	}
	
	if column.comment != "" {
		def.WriteString(fmt.Sprintf(" COMMENT '%s'", strings.ReplaceAll(column.comment, "'", "''")))
	}
	
	return def.String()
}

// Column methods for Blueprint

// ID adds an auto-incrementing ID column
func (b *Blueprint) ID() *Column {
	b.Primary("id")
	return b.BigInteger("id", true)
}

// String adds a VARCHAR column
func (b *Blueprint) String(name string, length int) *Column {
	column := Column{
		name:       name,
		columnType: "VARCHAR",
		length:     length,
		nullable:   false,
	}
	
	b.columns = append(b.columns, column)
	return &b.columns[len(b.columns)-1]
}

// Text adds a TEXT column
func (b *Blueprint) Text(name string) *Column {
	column := Column{
		name:       name,
		columnType: "TEXT",
		nullable:   false,
	}
	
	b.columns = append(b.columns, column)
	return &b.columns[len(b.columns)-1]
}

// Integer adds an INTEGER column
func (b *Blueprint) Integer(name string, autoIncrement bool) *Column {
	column := Column{
		name:       name,
		columnType: "INTEGER",
		nullable:   false,
	}
	
	if autoIncrement {
		column.columnType = "SERIAL"
	}
	
	b.columns = append(b.columns, column)
	return &b.columns[len(b.columns)-1]
}

// BigInteger adds a BIGINT column
func (b *Blueprint) BigInteger(name string, autoIncrement bool) *Column {
	column := Column{
		name:       name,
		columnType: "BIGINT",
		nullable:   false,
	}
	
	if autoIncrement {
		column.columnType = "BIGSERIAL"
	}
	
	b.columns = append(b.columns, column)
	return &b.columns[len(b.columns)-1]
}

// Boolean adds a BOOLEAN column
func (b *Blueprint) Boolean(name string) *Column {
	column := Column{
		name:       name,
		columnType: "BOOLEAN",
		nullable:   false,
	}
	
	b.columns = append(b.columns, column)
	return &b.columns[len(b.columns)-1]
}

// Date adds a DATE column
func (b *Blueprint) Date(name string) *Column {
	column := Column{
		name:       name,
		columnType: "DATE",
		nullable:   false,
	}
	
	b.columns = append(b.columns, column)
	return &b.columns[len(b.columns)-1]
}

// DateTime adds a TIMESTAMP column
func (b *Blueprint) DateTime(name string) *Column {
	column := Column{
		name:       name,
		columnType: "TIMESTAMP",
		nullable:   false,
	}
	
	b.columns = append(b.columns, column)
	return &b.columns[len(b.columns)-1]
}

// Decimal adds a DECIMAL column
func (b *Blueprint) Decimal(name string, precision, scale int) *Column {
	column := Column{
		name:       name,
		columnType: fmt.Sprintf("DECIMAL(%d,%d)", precision, scale),
		nullable:   false,
	}
	
	b.columns = append(b.columns, column)
	return &b.columns[len(b.columns)-1]
}

// Float adds a FLOAT column
func (b *Blueprint) Float(name string) *Column {
	column := Column{
		name:       name,
		columnType: "FLOAT",
		nullable:   false,
	}
	
	b.columns = append(b.columns, column)
	return &b.columns[len(b.columns)-1]
}

// JSON adds a JSON column
func (b *Blueprint) JSON(name string) *Column {
	column := Column{
		name:       name,
		columnType: "JSON",
		nullable:   false,
	}
	
	b.columns = append(b.columns, column)
	return &b.columns[len(b.columns)-1]
}

// Timestamps adds created_at and updated_at columns
func (b *Blueprint) Timestamps() {
	b.DateTime("created_at").Nullable()
	b.DateTime("updated_at").Nullable()
}

// TimestampsTz adds created_at and updated_at columns with timezone
func (b *Blueprint) TimestampsTz() {
	b.DateTime("created_at").Nullable()
	b.DateTime("updated_at").Nullable()
}

// SoftDeletes adds a deleted_at column
func (b *Blueprint) SoftDeletes() *Column {
	return b.DateTime("deleted_at").Nullable()
}

// Primary sets the primary key
func (b *Blueprint) Primary(columns ...string) {
	quoted := make([]string, len(columns))
	for i, col := range columns {
		quoted[i] = fmt.Sprintf("`%s`", col)
	}
	b.primary = strings.Join(quoted, ", ")
}

// Index adds an index
func (b *Blueprint) Index(columns ...string) {
	name := fmt.Sprintf("idx_%s_%s", b.table, strings.Join(columns, "_"))
	
	b.indexes = append(b.indexes, Index{
		name:    name,
		columns: columns,
		unique:  false,
	})
}

// Unique adds a unique index
func (b *Blueprint) Unique(columns ...string) {
	name := fmt.Sprintf("unq_%s_%s", b.table, strings.Join(columns, "_"))
	
	b.indexes = append(b.indexes, Index{
		name:    name,
		columns: columns,
		unique:  true,
	})
}

// Engine sets the storage engine
func (b *Blueprint) Engine(engine string) *Blueprint {
	b.engine = engine
	return b
}

// Charset sets the character set
func (b *Blueprint) Charset(charset string) *Blueprint {
	b.charset = charset
	return b
}

// Collation sets the collation
func (b *Blueprint) Collation(collation string) *Blueprint {
	b.collation = collation
	return b
}

// Temporary makes the table temporary
func (b *Blueprint) Temporary() *Blueprint {
	b.temporary = true
	return b
}

// Column modifiers

// Nullable makes the column nullable
func (c *Column) Nullable() *Column {
	c.nullable = true
	return c
}

// Default sets the default value
func (c *Column) Default(value interface{}) *Column {
	c.default_ = value
	return c
}

// Unsigned makes the column unsigned
func (c *Column) Unsigned() *Column {
	c.unsigned = true
	return c
}

// Unique makes the column unique
func (c *Column) Unique() *Column {
	c.unique = true
	return c
}

// Index adds an index to the column
func (c *Column) Index() *Column {
	c.index = true
	return c
}

// Comment adds a comment to the column
func (c *Column) Comment(comment string) *Column {
	c.comment = comment
	return c
}

// After positions the column after another column
func (c *Column) After(column string) *Column {
	c.after = column
	return c
}

// First positions the column first
func (c *Column) First() *Column {
	c.first = true
	return c
}

// Primary makes the column the primary key
func (c *Column) Primary() *Column {
	c.unique = true
	return c
}