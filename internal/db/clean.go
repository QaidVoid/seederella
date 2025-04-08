package db

import (
	"database/sql"
	"fmt"
)

func Clean(driver string, db *sql.DB, schema string) error {
	switch driver {
	case "postgres":
		if schema == "" {
			schema = "public"
		}
		_, err := db.Exec(fmt.Sprintf(`DROP SCHEMA %s CASCADE; CREATE SCHEMA %s;`, schema, schema))
		return err
	default:
		return fmt.Errorf("clean not implemented for driver: %s", driver)
	}
}
