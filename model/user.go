package model

import "golang.org/x/crypto/bcrypt"

type User struct {
	resource
	Username string  `json:"username"`
	Password string  `json:"-"`
	Email    string  `json:"email"`
	Timezone *string `json:"timezone"`
	InboxID  int64   `json:"-"`
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

func (r *CreateUserRequest) Patch(user *User) {
	if r.Username != nil {
		user.Username = *r.Username
	}

	if r.Email != nil {
		user.Email = *r.Email
	}

	if r.Timezone != nil {
		user.Timezone = r.Timezone
	}
}

func (r *UpdateUserRequest) Patch(user *User) {
	if r.Username != nil {
		user.Username = *r.Username
	}

	if r.NewPassword != nil {
		user.Password = *r.NewPassword
	}

	if r.Email != nil {
		user.Email = *r.Email
	}

	if r.Timezone != nil {
		user.Timezone = r.Timezone
	}
}
