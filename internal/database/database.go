package database

import (
	"database/sql"

	"github.com/CnTeng/rx-todo/internal/model"
	_ "github.com/lib/pq"
)

type DB struct {
	sql.DB
}

func NewDB(dsn string) (DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return DB{}, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(0)

	return DB{*db}, nil
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

func (db *DB) CreateUser(user *model.User) error {
	var hashedPassword string
	if user.Password != "" {
		var err error
		hashedPassword, err = model.HashPassword(user.Password)
		if err != nil {
			return err
		}
	}

	query := `
		INSERT INTO users (username, PASSWORD, email, timezone)
		  VALUES (LOWER($1), $2, $3, $4)
  `

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(query, user.Username, hashedPassword, user.Email, user.Timezone)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
