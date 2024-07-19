package model

import (
	"time"
)

type Duration struct {
	Amount *int    `json:"amount"`
	Unit   *string `json:"unit"`
}

type Task struct {
	ID          int64      `json:"id"`
	UserID      int64      `json:"user_id"`
	Content     string     `json:"content"`
	Description string     `json:"description"`
	Due         *Due       `json:"due,omitempty"`
	Duration    *Duration  `json:"duration,omitempty"`
	Priority    int        `json:"priority"`
	ProjectID   *int64     `json:"project_id,omitempty"`
	ParentID    *int64     `json:"parent_id,omitempty"`
	ChildOrder  int        `json:"child_order"`
	Done        bool       `json:"done"`
	DoneAt      *time.Time `json:"done_at,omitempty"`
	Archived    bool       `json:"archived"`
	ArchivedAt  *time.Time `json:"archived_at,omitempty"`

	// Meta fields
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateTaskRequest struct {
	Content     string    `json:"content"`
	Description *string   `json:"description"`
	Due         *Due      `json:"due,omitempty"`
	Duration    *Duration `json:"duration,omitempty"`
	Priority    *int      `json:"priority"`
	ProjectID   *int64    `json:"project_id,omitempty"`
	ParentID    *int64    `json:"parent_id,omitempty"`
	ChildOrder  *int      `json:"child_order"`
}

type UpdateTaskRequest struct {
	Content     *string   `json:"content"`
	Description *string   `json:"description"`
	Due         *Due      `json:"due"`
	Duration    *Duration `json:"duration"`
	Priority    *int      `json:"priority"`
	ProjectID   *int64    `json:"project_id"`
	ParentID    *int64    `json:"parent_id"`
	ChildOrder  *int      `json:"child_order"`
}

func (r *CreateTaskRequest) Patch(task *Task, userID, inboxID int64) {
	task.UserID = userID
	task.Content = r.Content

	if r.Description != nil {
		task.Description = *r.Description
	}

	if r.Due != nil {
		task.Due = r.Due
	}

	if r.Duration != nil {
		task.Duration = r.Duration
	}

	if r.Priority != nil {
		task.Priority = *r.Priority
	}

	if r.ProjectID != nil {
		task.ProjectID = r.ProjectID
	} else if r.ParentID != nil {
		task.ParentID = r.ParentID
	} else {
		task.ProjectID = &inboxID
	}

	if r.ChildOrder != nil {
		task.ChildOrder = *r.ChildOrder
	}
}

func (r *UpdateTaskRequest) Patch(task *Task) {
	if r.Content != nil {
		task.Content = *r.Content
	}

	if r.Description != nil {
		task.Description = *r.Description
	}

	if r.Due != nil {
		task.Due = r.Due
	}

	if r.Duration != nil {
		task.Duration = r.Duration
	}

	if r.Priority != nil {
		task.Priority = *r.Priority
	}

	if r.ProjectID != nil {
		task.ProjectID = r.ProjectID
	}

	if r.ParentID != nil {
		task.ParentID = r.ParentID
	}

	if r.ChildOrder != nil {
		task.ChildOrder = *r.ChildOrder
	}
}
