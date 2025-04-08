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
