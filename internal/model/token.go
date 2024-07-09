package model

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)

type Token struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	Token      string    `json:"token"`
	LastUsedAt time.Time `json:"last_used_at"`
	CreateAt   time.Time `json:"create_at"`
}

type TokenAddRequest struct {
	UserID int64 `json:"user_id"`
}

func NewToken(user int64) (string, error) {
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}
