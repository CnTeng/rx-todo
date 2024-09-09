package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int64   `json:"id"`
	Username string  `json:"username"`
	Password string  `json:"-"`
	Email    string  `json:"email"`
	Timezone *string `json:"timezone"`
	InboxID  int64   `json:"-"`

	// Meta fields
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
	Email    *string `json:"email"`
	Timezone *string `json:"timezone"`
}

type UpdateUserRequest struct {
	Username    *string `json:"username"`
	OldPassword *string `json:"old_password"`
	NewPassword *string `json:"new_password"`
	Email       *string `json:"email"`
	Timezone    *string `json:"timezone"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func (r *CreateUserRequest) Patch(user *User, hashedPassword string) {
	if r.Username != nil {
		user.Username = *r.Username
	}

	user.Password = hashedPassword

	if r.Email != nil {
		user.Email = *r.Email
	}

	if r.Timezone != nil {
		user.Timezone = r.Timezone
	}
}

func (r *UpdateUserRequest) Patch(user *User, hashedPassword string) {
	if r.Username != nil {
		user.Username = *r.Username
	}

	user.Password = hashedPassword

	if r.Email != nil {
		user.Email = *r.Email
	}

	if r.Timezone != nil {
		user.Timezone = r.Timezone
	}
}
