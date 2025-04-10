package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // MySQL
	_ "github.com/lib/pq"              // PostgreSQL
	_ "github.com/mattn/go-sqlite3"    // SQLite
)

func Connect(driver, dsn string, schema string) (*sql.DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	// set schema manually if PostgreSQL
	if driver == "postgres" && schema != "" {
		_, err := db.Exec(fmt.Sprintf("SET search_path TO %s", schema))
		if err != nil {
			return nil, fmt.Errorf("failed to set schema: %w", err)
		}
	}
	return db, nil
}

func CheckUniqueValue(db *sql.DB, driver, table, field string, value any) (bool, error) {
	var query string

	switch driver {
	case "postgres":
		query = fmt.Sprintf(`SELECT 1 FROM "%s" WHERE "%s" = $1 LIMIT 1`, table, field)
	case "mysql", "sqlite":
		query = fmt.Sprintf("SELECT 1 FROM `%s` WHERE `%s` = ? LIMIT 1", table, field)
	default:
		return false, fmt.Errorf("unsupported driver: %s", driver)
	}

	var exists bool
	err := db.QueryRow(query, value).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, fmt.Errorf("failed to check unique value for field '%s' in table '%s': %w", field, table, err)
	}

	return exists, nil
}
