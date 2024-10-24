package database

import (
	"database/sql"
	_ "embed"
	"fmt"
	"time"

	"github.com/CnTeng/rx-todo/model"
	"golang.org/x/crypto/bcrypt"
)

var (
	//go:embed sql/user_verify.sql
	verifyUserQuery string

	//go:embed sql/user_create.sql
	createUserQuery string

	//go:embed sql/user_create_inbox.sql
	createUserInboxQuery string

	//go:embed sql/user_update_inbox_id.sql
	updateUserInboxIDQuery string

	//go:embed sql/user_get_by_id.sql
	getUserByIDQuery string

	//go:embed sql/user_get_by_email.sql
	getUserByEmailQuery string

	//go:embed sql/user_get_inbox_id.sql
	getUserInboxIDQuery string

	//go:embed sql/user_get_by_updated_at.sql
	getUserByUpdateTimeQuery string

	//go:embed sql/user_update.sql
	updateUserQuery string

	//go:embed sql/user_delete.sql
	deleteUserQuery string
)

func (db *DB) VerifyUser(id int64, password string) error {
	var hashedPassword string

	if err := db.QueryRow(verifyUserQuery, id).Scan(&hashedPassword); err != nil {
		return fmt.Errorf("failed to verify user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return fmt.Errorf("failed to verify user: error password")
	}

	return nil
}

func (db *DB) CreateUser(user *model.User) (*model.User, error) {
	return user, db.withTx(func(tx *sql.Tx) error {
		if err := tx.QueryRow(
			createUserQuery,
			user.Username,
			user.Password,
			user.Email,
			user.Timezone,
		).Scan(
			&user.ID,
			&user.Timezone,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		if err := tx.QueryRow(createUserInboxQuery, user.ID, "Inbox", true).Scan(&user.InboxID); err != nil {
			return fmt.Errorf("failed to create user inbox: %w", err)
		}

		if _, err := tx.Exec(updateUserInboxIDQuery, user.ID, user.InboxID); err != nil {
			return fmt.Errorf("failed to create user inbox: %w", err)
		}

		return nil
	})
}

func (db *DB) GetUserByID(id int64) (*model.User, error) {
	return db.fetchUser(getUserByIDQuery, id)
}

func (db *DB) GetUserByEmail(email string) (*model.User, error) {
	return db.fetchUser(getUserByEmailQuery, email)
}

func (db *DB) GetUserByUpdateTime(id int64, updateTime *time.Time) (*model.User, error) {
	return db.fetchUser(getUserByUpdateTimeQuery, id, updateTime)
}

func (db *DB) GetUserInboxID(id int64) (int64, error) {
	var inboxID int64

	if err := db.QueryRow(getUserInboxIDQuery, id).Scan(&inboxID); err != nil {
		return 0, fmt.Errorf("failed to get user inbox id: %w", err)
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
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (db *DB) UpdateUser(user *model.User) (*model.User, error) {
	if err := db.QueryRow(
		updateUserQuery,
		user.ID,
		user.Username,
		user.Password,
		user.Email,
		user.Timezone,
	).Scan(&user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (db *DB) DeleteUser(id int64) error {
	if _, err := db.Exec(deleteUserQuery, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
