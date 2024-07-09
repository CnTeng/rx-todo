package database

import (
	"fmt"

	"github.com/CnTeng/rx-todo/internal/model"
)

func (db *DB) AddToken(r *model.TokenAddRequest) error {
	query := `
    INSERT INTO tokens (
      user_id,
      token
    )
    VALUES ($1, $2)
  `

	token, err := model.NewToken(r.UserID)
	if err != nil {
		return fmt.Errorf("database: unable to create token: %v", err)
	}

	_, err = db.Exec(query, r.UserID, token)
	if err != nil {
		return fmt.Errorf("database: unable to add token: %v", err)
	}

	return nil
}

func (db *DB) AuthToken(token *string) (int64, error) {
	query := `
		SELECT user_id
		FROM tokens
		WHERE token = $1
	`

	var userID int64
	err := db.QueryRow(query, token).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("database: unable to authenticate token: %v", err)
	}

	return userID, nil
}
