package db

import (
	"database/sql"
	"fmt"
)

func Clean(driver string, db *sql.DB, schema string) error {
	switch driver {
	case "postgres":
		return cleanPostgres(db, schema)
	default:
		return fmt.Errorf("clean not implemented for driver: %s", driver)
	}
}

func cleanPostgres(db *sql.DB, schema string) error {
	if schema == "" {
		schema = "public"
	}

	_, err := db.Exec("SET CONSTRAINTS ALL DEFERRED")
	if err != nil {
		return fmt.Errorf("failed to disable foreign key checks: %v", err)
	}

	tables, err := getTables(db, schema, "postgres")
	if err != nil {
		return fmt.Errorf("failed to retrieve tables: %v", err)
	}

	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s.\"%s\" RESTART IDENTITY CASCADE;", schema, table))
		if err != nil {
			return fmt.Errorf("failed to truncate table %s.%s: %v", schema, table, err)
		}
	}

	_, err = db.Exec("SET CONSTRAINTS ALL IMMEDIATE")
	if err != nil {
		return fmt.Errorf("failed to re-enable foreign key checks: %v", err)
	}

	return nil
}

func getTables(db *sql.DB, schema string, driver string) ([]string, error) {
	var query string
	var rows *sql.Rows
	var err error

	switch driver {
	case "postgres":
		query = `SELECT table_name FROM information_schema.tables WHERE table_schema = $1`
	// not sure about mysql / sqlite implementation yet
	// case "mysql":
	// case "sqlite":
	default:
		return nil, fmt.Errorf("unsupported driver: %s", driver)
	}

	rows, err = db.Query(query, schema)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tables, nil
}
