package model

import "time"

type Project struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	ProjectID   int64     `json:"project_id"`
	Content     string    `json:"content"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Archive     bool      `json:"archive"`
	ArchiveAt   time.Time `json:"archive_at"`
	ChildOrder  int       `json:"child_order"`
}
