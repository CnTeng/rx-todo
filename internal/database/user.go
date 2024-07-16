package database

import (
	"fmt"

	"github.com/CnTeng/rx-todo/internal/model"
	"golang.org/x/crypto/bcrypt"
)

func (db *DB) VerifyUser(id int64, password string) error {
	var hashedPassword string
	query := `SELECT password FROM users WHERE id = $1`

	err := db.QueryRow(query, id).Scan(&hashedPassword)
	if err != nil {
		return fmt.Errorf("database: unable to verify user: %v", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return fmt.Errorf("database: unable to verify user: error password")
	}

	return nil
}

func (db *DB) CreateUser(user *model.User) (*model.User, error) {
	query := `
    INSERT INTO users 
			(username, password, email, timezone)
    VALUES 
			(LOWER($1), $2, $3, COALESCE($4, 'UTC'))
		RETURNING 
			id, timezone, created_at, updated_at
  `

	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("database: unable to create user: %v", err)
	}
	defer func() { _ = tx.Rollback() }()

	// Create user
	err = tx.QueryRow(
		query,
		user.Username,
		user.Password,
		user.Email,
		user.Timezone,
	).Scan(
		&user.ID,
		&user.Timezone,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("database: unable to create user: %v", err)
	}

	// Create inbox for user
	err = tx.QueryRow(
		`
			INSERT INTO projects 
				(user_id, content, inbox) 
			VALUES 
				($1, $2, $3)
			RETURNING
				id
		`, user.ID, "Inbox", true).Scan(&user.InboxID)
	if err != nil {
		return nil, fmt.Errorf("database: unable to create user inbox: %v", err)
	}

	// Update user's inbox_id
	_, err = tx.Exec(`UPDATE users	SET inbox_id = $2 WHERE id = $1`, user.ID, user.InboxID)
	if err != nil {
		return nil, fmt.Errorf("database: unable to create user inbox: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("database: unable to commit transaction: %v", err)
	}

	return user, nil
}

func (db *DB) GetUserByID(id int64) (*model.User, error) {
	query := `
		SELECT
			id,
			username,
			password,
			email,
			timezone,
			created_at,
			updated_at
		FROM
			users
		WHERE
			id = $1
	`

	return db.fetchUser(query, id)
}

func (db *DB) GetUserByEmail(email string) (*model.User, error) {
	query := `
		SELECT
			id,
			username,
			password,
			email,
			timezone,
			created_at,
			updated_at
		FROM
			users
		WHERE
			email = $1
	`

	return db.fetchUser(query, email)
}

func (db *DB) GetUserInboxID(id int64) (int64, error) {
	var inboxID int64
	query := `SELECT inbox_id FROM users WHERE id = $1`

	err := db.QueryRow(query, id).Scan(&inboxID)
	if err != nil {
		return 0, fmt.Errorf("database: unable to get user inbox id: %v", err)
	}

	return inboxID, nil
}

func (db *DB) fetchUser(query string, args ...any) (*model.User, error) {
	user := new(model.User)

	err := db.QueryRow(query, args...).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.Timezone,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("database: unable to get user: %v", err)
	}

	return user, nil
}

func (db *DB) UpdateUser(user *model.User) (*model.User, error) {
	query := `
		UPDATE 
			users
		SET
			username = LOWER($2),
			password = $3,
			email = $4,
			timezone = $5
			updated_at = NOW()
		WHERE
			id = $1
		RETURNING 
			created_at,
			updated_at
	`

	err := db.QueryRow(
		query,
		user.ID,
		user.Username,
		user.Password,
		user.Email,
		user.Timezone,
	).Scan(&user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("database: unable to update user: %v", err)
	}

	return user, nil
}

func (db *DB) DeleteUser(id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("database: unable to delete user: %v", err)
	}

	return nil
}
