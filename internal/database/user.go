package database

import (
	"fmt"

	"github.com/CnTeng/rx-todo/internal/model"
)

func (db *DB) AddUser(r *model.User) error {
	query := `
    INSERT INTO users (
      username,
      password,
      email,
      timezone
    )
    VALUES ($1, $2, $3, $4)
  `
	password, err := model.HashPassword(r.Password)
	if err != nil {
		return fmt.Errorf("database: unable to hash password: %v", err)
	}

	_, err = db.Exec(query, r.Username, password, r.Email, r.Timezone)
	if err != nil {
		return fmt.Errorf("database: unable to add user: %v", err)
	}

	return nil
}
