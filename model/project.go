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
	ChildOrder  int64      `json:"child_order"`
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

// ProjectReorderMap represents a map of project id to child order
type ProjectReorderMap struct {
	ID         int64 `json:"id"`
	ChildOrder int64 `json:"child_order"`
}

// ProjectReorderRequest represents a request to reorder projects
type ProjectReorderRequest struct {
	ParentID *int64              `json:"parent_id"`
	Children []ProjectReorderMap `json:"children"`
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

func (m *ProjectReorderMap) Patch(project *Project) {
	project.ChildOrder = m.ChildOrder
}
