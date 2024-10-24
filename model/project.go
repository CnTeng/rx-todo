package model

import (
	"fmt"
	"time"
)

// Project represents a project
type Project struct {
	resource
	UserID      int64      `json:"user_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Position    float64    `json:"position"`
	Inbox       bool       `json:"inbox"`
	Favorite    bool       `json:"favorite"`
	SubTasks    SubTasks   `json:"sub_tasks"`
	Archived    bool       `json:"archived"`
	ArchivedAt  *time.Time `json:"archived_at,omitempty"`
}

// ProjectCreationRequest represents a request to create a project
type ProjectCreationRequest struct {
	Name        *string `json:"name"       validate:"required,notempty"`
	Description *string `json:"description"`
	Favorite    *bool   `json:"favorite"`
}

// ProjectUpdateRequest represents a request to update a project
type ProjectUpdateRequest struct {
	Name        *string `json:"name"        validate:"omitnil,notempty"`
	Description *string `json:"description" validate:"omitnil,notempty"`
}

// ProjectMoveRequest represents a request to reorder projects
type ProjectMoveRequest struct {
	PreviousID *int64 `json:"previous_id" validate:"required,notempty"`
}

func (r *ProjectUpdateRequest) validate() error {
	if r.Name == nil && r.Description == nil {
		return fmt.Errorf("validate: at least one of name or description should be set")
	}

	return nil
}

func (r *ProjectCreationRequest) Patch(project *Project) {
	if r.Name != nil {
		project.Name = *r.Name
	}

	if r.Description != nil {
		project.Description = *r.Description
	}

	if r.Favorite != nil {
		project.Favorite = *r.Favorite
	}
}

func (r *ProjectUpdateRequest) Patch(project *Project) {
	if r.Name != nil {
		project.Name = *r.Name
	}

	if r.Description != nil {
		project.Description = *r.Description
	}
}
