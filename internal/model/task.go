package model

import (
	"fmt"
	"time"

	"github.com/lib/pq"
)

type Duration struct {
	Amount *int    `json:"amount"`
	Unit   *string `json:"unit"`
}

type Task struct {
	ID          int64          `json:"id"`
	UserID      int64          `json:"user_id"`
	Content     string         `json:"content"`
	Description string         `json:"description"`
	Due         *Due           `json:"due,omitempty"`
	Duration    *Duration      `json:"duration,omitempty"`
	Priority    int            `json:"priority"`
	ProjectID   *int64         `json:"project_id,omitempty"`
	ParentID    *int64         `json:"parent_id,omitempty"`
	ChildOrder  int            `json:"child_order"`
	Labels      pq.StringArray `json:"labels"`
	Done        bool           `json:"done"`
	DoneAt      *time.Time     `json:"done_at,omitempty"`
	Archived    bool           `json:"archived"`
	ArchivedAt  *time.Time     `json:"archived_at,omitempty"`

	// Meta fields
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TaskCreationRequest struct {
	Content     *string   `json:"content" validate:"required,notempty"`
	Description *string   `json:"description"`
	Due         *Due      `json:"due,omitempty"`
	Duration    *Duration `json:"duration,omitempty"`
	Priority    *int      `json:"priority"`
	ProjectID   *int64    `json:"project_id,omitempty"`
	ParentID    *int64    `json:"parent_id,omitempty"`
	ChildOrder  *int      `json:"child_order"`
	Labels      *[]string `json:"labels"`
}

type TaskUpdateRequest struct {
	Content     *string   `json:"content"`
	Description *string   `json:"description"`
	Due         *Due      `json:"due"`
	Duration    *Duration `json:"duration"`
	Priority    *int      `json:"priority"`
	ProjectID   *int64    `json:"project_id"`
	ParentID    *int64    `json:"parent_id"`
	ChildOrder  *int      `json:"child_order"`
	Labels      *[]string `json:"labels"`
}

func (r *TaskUpdateRequest) Validate() error {
	if r.Content == nil && r.Description == nil &&
		r.Due == nil && r.Duration == nil &&
		r.Priority == nil && r.ProjectID == nil &&
		r.ParentID == nil && r.ChildOrder == nil &&
		r.Labels == nil {
		return fmt.Errorf("validate: at least one of task attributes should be set")
	}

	return nil
}

func (r *TaskCreationRequest) Patch(task *Task, userID, inboxID int64) {
	task.UserID = userID

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
	} else if r.ParentID != nil {
		task.ParentID = r.ParentID
	} else {
		task.ProjectID = &inboxID
	}

	if r.ChildOrder != nil {
		task.ChildOrder = *r.ChildOrder
	}

	if r.Labels != nil {
		task.Labels = *r.Labels
	}
}

func (r *TaskUpdateRequest) Patch(task *Task) {
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

	if r.Labels != nil {
		task.Labels = *r.Labels
	}
}
