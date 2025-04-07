package db

import (
	"fmt"
	"strings"
)

type QueryBuilder struct {
	db         *Database
	table      string
	columns    []string
	wheres     []string
	whereBinds []interface{}
	joins      []string
	orderBys   []string
	groupBys   []string
	havings    []string
	limit      int
	offset     int
	distinct   bool
	unions     []string
	binds      []interface{}
}

func NewQueryBuilder(db *Database, table string) *QueryBuilder {
	return &QueryBuilder{
		db:       db,
		table:    table,
		columns:  []string{"*"},
		wheres:   []string{},
		joins:    []string{},
		orderBys: []string{},
		groupBys: []string{},
		havings:  []string{},
		unions:   []string{},
	}
}

func (q *QueryBuilder) Select(columns ...string) *QueryBuilder {
	if len(columns) > 0 {
		q.columns = columns
	}
	return q
}

func (q *QueryBuilder) Distinct() *QueryBuilder {
	q.distinct = true
	return q
}

func (q *QueryBuilder) Where(column string, operator string, value interface{}) *QueryBuilder {
	q.wheres = append(q.wheres, fmt.Sprintf("%s %s ?", column, operator))
	q.whereBinds = append(q.whereBinds, value)
	return q
}

func (q *QueryBuilder) WhereRaw(sql string, binds ...interface{}) *QueryBuilder {
	q.wheres = append(q.wheres, sql)
	q.whereBinds = append(q.whereBinds, binds...)
	return q
}

func (q *QueryBuilder) OrWhere(column string, operator string, value interface{}) *QueryBuilder {
	if len(q.wheres) == 0 {
		return q.Where(column, operator, value)
	}
	q.wheres = append(q.wheres, fmt.Sprintf("OR %s %s ?", column, operator))
	q.whereBinds = append(q.whereBinds, value)
	return q
}

func (q *QueryBuilder) WhereIn(column string, values ...interface{}) *QueryBuilder {
	if len(values) == 0 {
		q.wheres = append(q.wheres, "1=0")
		return q
	}
	placeholders := make([]string, len(values))
	for i := range values {
		placeholders[i] = "?"
	}
	q.wheres = append(q.wheres, fmt.Sprintf("%s IN (%s)", column, strings.Join(placeholders, ", ")))
	q.whereBinds = append(q.whereBinds, values...)
	return q
}

func (q *QueryBuilder) WhereNotIn(column string, values ...interface{}) *QueryBuilder {
	if len(values) == 0 {
		return q
	}
	placeholders := make([]string, len(values))
	for i := range values {
		placeholders[i] = "?"
	}
	q.wheres = append(q.wheres, fmt.Sprintf("%s NOT IN (%s)", column, strings.Join(placeholders, ", ")))
	q.whereBinds = append(q.whereBinds, values...)
	return q
}

func (q *QueryBuilder) WhereNull(column string) *QueryBuilder {
	q.wheres = append(q.wheres, fmt.Sprintf("%s IS NULL", column))
	return q
}

func (q *QueryBuilder) WhereNotNull(column string) *QueryBuilder {
	q.wheres = append(q.wheres, fmt.Sprintf("%s IS NOT NULL", column))
	return q
}

func (q *QueryBuilder) Join(table string, first string, operator string, second string) *QueryBuilder {
	q.joins = append(q.joins, fmt.Sprintf("JOIN %s ON %s %s %s", table, first, operator, second))
	return q
}

func (q *QueryBuilder) LeftJoin(table string, first string, operator string, second string) *QueryBuilder {
	q.joins = append(q.joins, fmt.Sprintf("LEFT JOIN %s ON %s %s %s", table, first, operator, second))
	return q
}

func (q *QueryBuilder) RightJoin(table string, first string, operator string, second string) *QueryBuilder {
	q.joins = append(q.joins, fmt.Sprintf("RIGHT JOIN %s ON %s %s %s", table, first, operator, second))
	return q
}

func (q *QueryBuilder) OrderBy(column string, direction string) *QueryBuilder {
	dir := strings.ToUpper(direction)
	if dir != "ASC" && dir != "DESC" {
		dir = "ASC"
	}
	q.orderBys = append(q.orderBys, fmt.Sprintf("%s %s", column, dir))
	return q
}

func (q *QueryBuilder) GroupBy(columns ...string) *QueryBuilder {
	q.groupBys = append(q.groupBys, columns...)
	return q
}

func (q *QueryBuilder) Having(column string, operator string, value interface{}) *QueryBuilder {
	q.havings = append(q.havings, fmt.Sprintf("%s %s ?", column, operator))
	q.binds = append(q.binds, value)
	return q
}

func (q *QueryBuilder) Limit(limit int) *QueryBuilder {
	if limit >= 0 {
		q.limit = limit
	}
	return q
}

func (q *QueryBuilder) Offset(offset int) *QueryBuilder {
	if offset >= 0 {
		q.offset = offset
	}
	return q
}

func (q *QueryBuilder) Union(query *QueryBuilder) *QueryBuilder {
	sql, binds := query.ToSql()
	q.unions = append(q.unions, fmt.Sprintf("UNION (%s)", sql))
	q.binds = append(q.binds, binds...)
	return q
}

func (q *QueryBuilder) UnionAll(query *QueryBuilder) *QueryBuilder {
	sql, binds := query.ToSql()
	q.unions = append(q.unions, fmt.Sprintf("UNION ALL (%s)", sql))
	q.binds = append(q.binds, binds...)
	return q
}

func (q *QueryBuilder) ToSql() (string, []interface{}) {
	var query strings.Builder
	var binds []interface{}

	query.WriteString("SELECT ")
	if q.distinct {
		query.WriteString("DISTINCT ")
	}
	query.WriteString(strings.Join(q.columns, ", "))
	query.WriteString(" FROM ")
	query.WriteString(q.table)

	if len(q.joins) > 0 {
		query.WriteString(" ")
		query.WriteString(strings.Join(q.joins, " "))
	}

	if len(q.wheres) > 0 {
		query.WriteString(" WHERE ")
		query.WriteString(q.wheres[0])
		binds = append(binds, q.whereBinds[0])
		for i := 1; i < len(q.wheres); i++ {
			clause := q.wheres[i]
			if strings.HasPrefix(clause, "OR ") {
				query.WriteString(" ")
				query.WriteString(clause)
			} else {
				query.WriteString(" AND ")
				query.WriteString(clause)
			}
			binds = append(binds, q.whereBinds[i])
		}
	}

	if len(q.groupBys) > 0 {
		query.WriteString(" GROUP BY ")
		query.WriteString(strings.Join(q.groupBys, ", "))
	}

	if len(q.havings) > 0 {
		query.WriteString(" HAVING ")
		query.WriteString(strings.Join(q.havings, " AND "))
		binds = append(binds, q.binds...)
	}

	if len(q.orderBys) > 0 {
		query.WriteString(" ORDER BY ")
		query.WriteString(strings.Join(q.orderBys, ", "))
	}

	if q.limit > 0 {
		query.WriteString(fmt.Sprintf(" LIMIT %d", q.limit))
	}

	if q.offset > 0 {
		query.WriteString(fmt.Sprintf(" OFFSET %d", q.offset))
	}

	if len(q.unions) > 0 {
		query.WriteString(" ")
		query.WriteString(strings.Join(q.unions, " "))
	}

	return query.String(), binds
}

func (q *QueryBuilder) Get(dest interface{}) error {
	sql, binds := q.ToSql()
	rows, err := q.db.Query(sql, binds...)
	if err != nil {
		return err
	}
	defer rows.Close()
	return q.db.scanRows(rows, dest)
}

func (q *QueryBuilder) First(dest interface{}) error {
	q.Limit(1)
	sql, binds := q.ToSql()
	row := q.db.QueryRow(sql, binds...)
	return row.Scan(dest)
}

func (q *QueryBuilder) Count() (int, error) {
	var count int
	countQuery := NewQueryBuilder(q.db, q.table)
	countQuery.wheres = q.wheres
	countQuery.whereBinds = q.whereBinds
	countQuery.joins = q.joins
	countQuery.columns = []string{"COUNT(*) as count"}
	countQuery.orderBys = []string{}
	sql, binds := countQuery.ToSql()
	row := q.db.QueryRow(sql, binds...)
	err := row.Scan(&count)
	return count, err
}

func (q *QueryBuilder) Insert(values map[string]interface{}) error {
	if len(values) == 0 {
		return fmt.Errorf("no values provided for insert")
	}
	var columns []string
	var placeholders []string
	var binds []interface{}
	for column, value := range values {
		columns = append(columns, column)
		placeholders = append(placeholders, "?")
		binds = append(binds, value)
	}
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", q.table, strings.Join(columns, ", "), strings.Join(placeholders, ", "))
	err := q.db.Exec(query, binds...)
	return err
}

func (q *QueryBuilder) Update(values map[string]interface{}) error {
	if len(values) == 0 {
		return fmt.Errorf("no values provided for update")
	}
	var sets []string
	var binds []interface{}
	for column, value := range values {
		sets = append(sets, fmt.Sprintf("%s = ?", column))
		binds = append(binds, value)
	}
	query := fmt.Sprintf("UPDATE %s SET %s", q.table, strings.Join(sets, ", "))
	if len(q.wheres) > 0 {
		query += " WHERE " + strings.Join(q.wheres, " AND ")
		binds = append(binds, q.whereBinds...)
	}
	err := q.db.Exec(query, binds...)
	return err
}

func (q *QueryBuilder) Delete() error {
	query := fmt.Sprintf("DELETE FROM %s", q.table)
	var binds []interface{}
	if len(q.wheres) > 0 {
		query += " WHERE " + strings.Join(q.wheres, " AND ")
		binds = append(binds, q.whereBinds...)
	}
	err := q.db.Exec(query, binds...)
	return err
}