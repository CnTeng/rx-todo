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
	resource
	UserID      int64          `json:"user_id"`
	Name        string         `json:"name"`
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
}

// TaskCreationRequest represents a request to create a task
type TaskCreationRequest struct {
	Name        *string   `json:"name"        toml:"name"       validate:"required,notempty"`
	Description *string   `json:"description" toml:"description"`
	Due         *Due      `json:"due"         toml:"due"`
	Duration    *Duration `json:"duration"    toml:"duration"`
	Priority    *int      `json:"priority"    toml:"priority"`
	ProjectID   *int64    `json:"project_id"  toml:"project_id"`
	ParentID    *int64    `json:"parent_id"   toml:"parent_id"`
	Labels      *[]string `json:"labels"      toml:"labels"`
}

// TaskUpdateRequest represents a request to update a task
type TaskUpdateRequest struct {
	Name        *string   `json:"name"        toml:"name"`
	Description *string   `json:"description" toml:"description"`
	Due         *Due      `json:"due"         toml:"due"`
	Duration    *Duration `json:"duration"    toml:"duration"`
	Priority    *int      `json:"priority"    toml:"priority"`
	ProjectID   *int64    `json:"project_id"  toml:"project_id"`
	ParentID    *int64    `json:"parent_id"   toml:"parent_id"`
	ChildOrder  *int      `json:"child_order" toml:"child_order"`
	Labels      *[]string `json:"labels"      toml:"labels"`
}

// TaskUpdateRequestWithID represents a request to update a task with an ID
type TaskUpdateRequestWithID struct {
	ID int64 `json:"id" validate:"required,notempty"`
	TaskUpdateRequest
}

type TaskDeleteRequestWithID struct {
	ID int64 `json:"id" validate:"required,notempty"`
}

// TaskUpdateRequest at least needs one of task attributes to be set
func (r *TaskUpdateRequest) Validate() error {
	if r.Name == nil && r.Description == nil &&
		r.Due == nil && r.Duration == nil &&
		r.Priority == nil && r.ProjectID == nil &&
		r.ParentID == nil && r.ChildOrder == nil &&
		r.Labels == nil {
		return fmt.Errorf("validate: at least one of task attributes should be set")
	}

	if r.ProjectID != nil && r.ParentID != nil {
		return fmt.Errorf("validate: only one of project_id or parent_id should be set")
	}

	return nil
}

func (r *TaskCreationRequest) Patch(task *Task) {
	if r.Name != nil {
		task.Name = *r.Name
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
		task.ProjectID = nil
		task.ParentID = r.ParentID
	}

	if r.Labels != nil {
		task.Labels = *r.Labels
	}
}

func (r *TaskUpdateRequest) Patch(task *Task) {
	if r.Name != nil {
		task.Name = *r.Name
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
