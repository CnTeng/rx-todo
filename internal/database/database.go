package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func NewDB(dsn string) (DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return DB{nil}, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(0)

	return DB{db}, nil
}

func (db *DB) Migrate() error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	for _, migration := range migrations {
		if err := migration(tx); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (db *DB) execSimpleQuery(query string, userID, id int64) error {
	_, err := db.Exec(query, userID, id)
	if err != nil {
		return fmt.Errorf("database: unable to exec query: %v", err)
	}

	return nil
}
