package model

import "time"

// Label represents a label
type Label struct {
	ID     int64  `json:"id"`
	UserID int64  `json:"user_id"`
	Name   string `json:"name"`
	Color  string `json:"color"`

	// CreatedAt is the time the label was created, Read-only.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the time the label was updated, Read-only.
	UpdatedAt time.Time `json:"updated_at"`
}
