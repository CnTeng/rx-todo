package database

import (
	"fmt"

	"github.com/CnTeng/rx-todo/internal/model"
)

func (db *DB) CreateToken(token *model.Token) (*model.Token, error) {
	query := `
    INSERT INTO tokens 
			(user_id, token)
    VALUES 
			($1, $2)
		RETURNING
			id, token, created_at, updated_at
  `

	t, err := model.NewToken(token.UserID)
	if err != nil {
		return nil, fmt.Errorf("database: unable to create token: %v", err)
	}

	err = db.QueryRow(query, token.UserID, t).Scan(&token.ID, &token.Token, &token.CreatedAt, &token.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("database: unable to add token: %v", err)
	}

	return token, nil
}

func (db *DB) GetUserIDByToken(token *string) (int64, error) {
	var userID int64
	query := `SELECT user_id FROM tokens WHERE token = $1`

	err := db.QueryRow(query, token).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("database: unable to authenticate token: %v", err)
	}

	return userID, nil
}

func (db *DB) GetTokenByID(userID, id int64) (*model.Token, error) {
	t := &model.Token{}
	query := `
		SELECT
			id,
			user_id,
			token,
			created_at,
			updated_at
		FROM
			tokens
		WHERE
			user_id = $1 AND id = $2
	`

	err := db.QueryRow(query, userID, id).Scan(&t.ID, &t.UserID, &t.Token, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("database: unable to get token: %v", err)
	}

	return t, nil
}

func (db *DB) GetTokens(userID int64) ([]*model.Token, error) {
	var tokens []*model.Token
	query := `
		SELECT
			id,
			user_id,
			token,
			created_at,
			updated_at
		FROM
			tokens
		WHERE
			user_id = $1
	`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("database: unable to get tokens: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		t := &model.Token{}

		if err := rows.Scan(&t.ID, &t.UserID, &t.Token, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("database: unable to get tokens: %v", err)
		}

		tokens = append(tokens, t)
	}

	return tokens, nil
}

func (db *DB) UpdateToken(token *model.Token) (*model.Token, error) {
	query := `
		UPDATE 
			tokens
		SET
			token = $3,
			updated_at = NOW()
		WHERE
			id = $1 AND user_id = $2
		RETURNING
			id,
			token,
			created_at,
			updated_at
	`

	t, err := model.NewToken(token.UserID)
	if err != nil {
		return nil, fmt.Errorf("database: unable to create token: %v", err)
	}

	err = db.QueryRow(
		query,
		token.ID,
		token.UserID,
		t,
	).Scan(&token.ID, &token.Token, &token.CreatedAt, &token.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("database: unable to update token: %v", err)
	}

	return token, nil
}

func (db *DB) DeleteToken(userID, id int64) error {
	query := `DELETE FROM tokens WHERE id = $1 AND user_id = $2`

	_, err := db.Exec(query, id, userID)
	if err != nil {
		return fmt.Errorf("database: unable to delete token: %v", err)
	}

	return nil
}
