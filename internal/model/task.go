package model

import (
	"fmt"
	"time"
)

type Duration struct {
	Amount *int    `json:"amount"`
	Unit   *string `json:"unit"`
}

type Priority int

const (
	PriorityNone Priority = iota
	PriorityLow
	PriorityMedium
	PriorityHigh
)

type Task struct {
	resource
	UserID      int64      `json:"user_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Due         *Due       `json:"due,omitempty"`
	Duration    *Duration  `json:"duration,omitempty"`
	Priority    Priority   `json:"priority"`
	ProjectID   *int64     `json:"project_id,omitempty"`
	ParentID    *int64     `json:"parent_id,omitempty"`
	Position    int64      `json:"position"`
	Labels      []*Label   `json:"labels"`
	Progress    *Progress  `json:"Progress,omitempty"`
	Done        bool       `json:"done"`
	DoneAt      *time.Time `json:"done_at,omitempty"`
	Archived    bool       `json:"archived"`
	ArchivedAt  *time.Time `json:"archived_at,omitempty"`
}

// TaskCreationRequest represents a request to create a task
type TaskCreationRequest struct {
	Name        *string   `json:"name"        toml:"name"       validate:"required,notempty"`
	Description *string   `json:"description" toml:"description"`
	Due         *Due      `json:"due"         toml:"due"`
	Duration    *Duration `json:"duration"    toml:"duration"`
	Priority    *Priority `json:"priority"    toml:"priority"`
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
	Priority    *Priority `json:"priority"    toml:"priority"`
	Labels      *[]string `json:"labels"      toml:"labels"`
}

type TaskMoveRequest struct {
	PreviousID *int64 `json:"previous_id"`
	ProjectID  *int64 `json:"project_id"`
	ParentID   *int64 `json:"parent_id"`
}

// TaskUpdateRequest at least needs one of task attributes to be set
func (r *TaskUpdateRequest) validate() error {
	if r.Name == nil && r.Description == nil &&
		r.Due == nil && r.Duration == nil &&
		r.Priority == nil &&
		r.Labels == nil {
		return fmt.Errorf("validate: at least one of task attributes should be set")
	}

	return nil
}

func (r *TaskMoveRequest) validate() error {
	if r.ProjectID == nil && r.ParentID == nil {
		return fmt.Errorf("validate: at least one of project_id or parent_id should be set")
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
		task.Labels = make([]*Label, len(*r.Labels))
		for i, l := range *r.Labels {
			task.Labels[i] = &Label{Name: l}
		}
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

	if r.Labels != nil {
		task.Labels = make([]*Label, len(*r.Labels))
		for i, l := range *r.Labels {
			task.Labels[i] = &Label{Name: l}
		}
	}
}

func (r *TaskMoveRequest) Patch(task *Task) {
	task.ProjectID = r.ProjectID
	task.ParentID = r.ParentID
}
