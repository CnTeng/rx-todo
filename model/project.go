package model

import (
	"fmt"
	"time"
)

// Project represents a project
type Project struct {
	Resource
	UserID      int64      `json:"user_id"`
	Content     string     `json:"content"`
	Description string     `json:"description"`
	ParentID    *int64     `json:"parent_id"`
	ChildOrder  int64      `json:"child_order"`
	Inbox       bool       `json:"inbox"`
	Favorite    bool       `json:"favorite"`
	Archived    bool       `json:"archived"`
	ArchivedAt  *time.Time `json:"archived_at,omitempty"`
}

// ProjectCreationRequest represents a request to create a project
type ProjectCreationRequest struct {
	Content     *string `json:"content" validate:"required,notempty"`
	Description *string `json:"description"`
	ParentID    *int64  `json:"parent_id"`
	Favorite    *bool   `json:"favorite"`
}

// ProjectUpdateRequest represents a request to update a project
type ProjectUpdateRequest struct {
	Content     *string `json:"content" validate:"omitnil,notempty"`
	Description *string `json:"description" validate:"omitnil,notempty"`
}

// ProjectMoveRequest represents a request to move a project
type ProjectMoveRequest struct {
	ParentID int64 `json:"parent_id" validate:"required"`
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
	if r.Content == nil && r.Description == nil {
		return fmt.Errorf("validate: at least one of content or description should be set")
	}

	return nil
}

func (r *ProjectCreationRequest) Patch(project *Project) {
	if r.Content != nil {
		project.Content = *r.Content
	}

	if r.Description != nil {
		project.Description = *r.Description
	}

	if r.ParentID != nil {
		project.ParentID = r.ParentID
	}

	if r.Favorite != nil {
		project.Favorite = *r.Favorite
	}
}

func (r *ProjectUpdateRequest) Patch(project *Project) {
	if r.Content != nil {
		project.Content = *r.Content
	}

	if r.Description != nil {
		project.Description = *r.Description
	}
}

func (m *ProjectReorderMap) Patch(project *Project) {
	project.ChildOrder = m.ChildOrder
}
