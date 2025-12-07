package database

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLClient struct {
	db *sql.DB
}

func NewMySQLClient(dsn string) (*MySQLClient, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexão: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao pingar banco: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return &MySQLClient{db: db}, nil
}

func (c *MySQLClient) Close() error {
	return c.db.Close()
}

func (c *MySQLClient) Create(tableName string, data map[string]interface{}) (int64, error) {
	if len(data) == 0 {
		return 0, fmt.Errorf("dados vazios para inserção")
	}

	columns := make([]string, 0, len(data))
	placeholders := make([]string, 0, len(data))
	values := make([]interface{}, 0, len(data))

	for col, val := range data {
		columns = append(columns, col)
		placeholders = append(placeholders, "?")
		values = append(values, val)
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	result, err := c.db.Exec(query, values...)
	if err != nil {
		return 0, fmt.Errorf("erro ao executar INSERT: %w", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("erro ao obter último ID: %w", err)
	}

	return lastID, nil
}

func (c *MySQLClient) Read(tableName string, columns []string, filter map[string]interface{}) ([]map[string]interface{}, error) {
	columnsSQL := "*"
	if len(columns) > 0 {
		columnsSQL = strings.Join(columns, ", ")
	}

	whereClause := ""
	values := make([]interface{}, 0)

	if len(filter) > 0 {
		whereParts := make([]string, 0, len(filter))
		for col, val := range filter {
			whereParts = append(whereParts, fmt.Sprintf("%s = ?", col))
			values = append(values, val)
		}
		whereClause = " WHERE " + strings.Join(whereParts, " AND ")
	}

	query := fmt.Sprintf("SELECT %s FROM %s%s", columnsSQL, tableName, whereClause)

	rows, err := c.db.Query(query, values...)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar SELECT: %w", err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("erro ao obter colunas: %w", err)
	}

	results := make([]map[string]interface{}, 0)

	for rows.Next() {
		columnValues := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columnValues {
			columnPointers[i] = &columnValues[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, fmt.Errorf("erro ao escanear linha: %w", err)
		}

		row := make(map[string]interface{})
		for i, colName := range cols {
			val := columnValues[i]
			if b, ok := val.([]byte); ok {
				row[colName] = string(b)
			} else {
				row[colName] = val
			}
		}

		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar linhas: %w", err)
	}

	return results, nil
}

func (c *MySQLClient) Update(tableName string, data map[string]interface{}, filter map[string]interface{}) (int64, error) {
	if len(data) == 0 {
		return 0, fmt.Errorf("dados vazios para atualização")
	}

	if len(filter) == 0 {
		return 0, fmt.Errorf("condições vazias - UPDATE sem WHERE não é permitido")
	}

	setParts := make([]string, 0, len(data))
	values := make([]interface{}, 0, len(data)+len(filter))

	for col, val := range data {
		setParts = append(setParts, fmt.Sprintf("%s = ?", col))
		values = append(values, val)
	}

	whereParts := make([]string, 0, len(filter))
	for col, val := range filter {
		whereParts = append(whereParts, fmt.Sprintf("%s = ?", col))
		values = append(values, val)
	}

	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE %s",
		tableName,
		strings.Join(setParts, ", "),
		strings.Join(whereParts, " AND "),
	)

	result, err := c.db.Exec(query, values...)
	if err != nil {
		return 0, fmt.Errorf("erro ao executar UPDATE: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("erro ao obter linhas afetadas: %w", err)
	}

	return rowsAffected, nil
}

func (c *MySQLClient) Delete(tableName string, filter map[string]interface{}) (int64, error) {
	if len(filter) == 0 {
		return 0, fmt.Errorf("condições vazias - DELETE sem WHERE não é permitido")
	}

	whereParts := make([]string, 0, len(filter))
	values := make([]interface{}, 0, len(filter))

	for col, val := range filter {
		whereParts = append(whereParts, fmt.Sprintf("%s = ?", col))
		values = append(values, val)
	}

	query := fmt.Sprintf(
		"DELETE FROM %s WHERE %s",
		tableName,
		strings.Join(whereParts, " AND "),
	)

	result, err := c.db.Exec(query, values...)
	if err != nil {
		return 0, fmt.Errorf("erro ao executar DELETE: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("erro ao obter linhas afetadas: %w", err)
	}

	return rowsAffected, nil
}
