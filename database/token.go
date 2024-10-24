package database

import (
	_ "embed"
	"fmt"

	"github.com/CnTeng/rx-todo/model"
)

var (
	//go:embed sql/token_create.sql
	createTokenQuery string

	//go:embed sql/token_get_user_id.sql
	getUserIDByTokenQuery string

	//go:embed sql/token_get_by_id.sql
	getTokenByIDQuery string

	//go:embed sql/token_get_all.sql
	getTokensQuery string

	//go:embed sql/token_update.sql
	updateTokenQuery string

	//go:embed sql/token_delete.sql
	deleteTokenQuery string
)

func (db *DB) CreateToken(token *model.Token) (*model.Token, error) {
	newToken, err := model.NewToken(token.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to create token: %w", err)
	}

	if err := db.
		QueryRow(createTokenQuery, token.UserID, newToken).Scan(
		&token.ID,
		&token.Token,
		&token.CreatedAt,
		&token.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to create token: %w", err)
	}

	return token, nil
}

func (db *DB) GetUserIDByToken(token *string) (int64, error) {
	var userID int64

	if err := db.QueryRow(getUserIDByTokenQuery, token).Scan(&userID); err != nil {
		return 0, fmt.Errorf("failed to authenticate token: %w", err)
	}

	return userID, nil
}

func (db *DB) GetTokenByID(userID, id int64) (*model.Token, error) {
	t := &model.Token{}

	if err := db.QueryRow(getTokenByIDQuery, userID, id).Scan(
		&t.ID,
		&t.UserID,
		&t.Token,
		&t.CreatedAt,
		&t.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	return t, nil
}

func (db *DB) GetTokens(userID int64) ([]*model.Token, error) {
	var tokens []*model.Token

	rows, err := db.Query(getTokensQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tokens: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		token := &model.Token{}

		if err := rows.Scan(
			&token.ID,
			&token.UserID,
			&token.Token,
			&token.CreatedAt,
			&token.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to get tokens: %w", err)
		}

		tokens = append(tokens, token)
	}

	return tokens, nil
}

func (db *DB) UpdateToken(token *model.Token) (*model.Token, error) {
	newToken, err := model.NewToken(token.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to create token: %w", err)
	}

	if err := db.QueryRow(
		updateTokenQuery,
		token.ID,
		token.UserID,
		newToken,
	).Scan(
		&token.ID,
		&token.Token,
		&token.CreatedAt,
		&token.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to update token: %w", err)
	}

	return token, nil
}

func (db *DB) DeleteToken(token *model.Token) error {
	if _, err := db.Exec(deleteTokenQuery, token.ID, token.UserID); err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}

	return nil
}
