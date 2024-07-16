package model

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)

type Token struct {
	ID     int64  `json:"id"`
	UserID int64  `json:"user_id"`
	Token  string `json:"token"`

	// Meta fields
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateTokenRequest struct {
	UserID   int64  `json:"user_id"`
	Password string `json:"password"`
}

type UpdateTokenRequest CreateTokenRequest

func NewToken(user int64) (string, error) {
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}

func (r *CreateTokenRequest) Patch(token *Token) {
	token.UserID = r.UserID
}

func (r *UpdateTokenRequest) Patch(token *Token) {
	token.UserID = r.UserID
}
