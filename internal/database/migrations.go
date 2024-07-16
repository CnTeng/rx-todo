package database

import (
	"database/sql"
	_ "embed"
	"fmt"
)

//go:embed migrations/0001_create.sql
var initSQL string

var migrations = []func(tx *sql.Tx) error{
	func(tx *sql.Tx) error {
		sql := initSQL

		_, err := tx.Exec(sql)
		if err != nil {
			return fmt.Errorf("database: unable to create table: %v", err)
		}

		return nil
	},
}
