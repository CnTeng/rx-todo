package database

import (
	"database/sql"
	_ "embed"
	"fmt"

	_ "github.com/lib/pq"
)

//go:embed sql/deletion_log_create.sql
var createDeletionLogQuery string

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

type txFunc func(tx *sql.Tx) error

func (db *DB) withTx(fns ...txFunc) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	for _, fn := range fns {
		if err := fn(tx); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
