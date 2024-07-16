package model

import "time"

type Project struct {
	ID          int64      `json:"id"`
	UserID      int64      `json:"user_id"`
	Content     string     `json:"content"`
	Description string     `json:"description"`
	ParentID    *int64     `json:"parent_id"`
	ChildOrder  int64      `json:"child_order"`
	Inbox       bool       `json:"inbox"`
	Favorite    bool       `json:"favorite"`
	Archived    bool       `json:"archived"`
	ArchivedAt  *time.Time `json:"archived_at,omitempty"`

	// Meta fields
	Position  float64   `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateProjectRequest struct {
	Content     *string `json:"content"`
	Description *string `json:"description"`
	ParentID    *int64  `json:"parent_id"`
	ChildOrder  *int64  `json:"child_order"`
	Favorite    *bool   `json:"favorite"`
}

type UpdateProjectRequest struct {
	ID          *int64  `json:"id"`
	Content     *string `json:"content"`
	Description *string `json:"description"`
}

type MoveProjectRequest struct {
	ID       int64  `json:"id"`
	ParentID *int64 `json:"parent_id"`
}

type ReorderProjectMap struct {
	ID         int64 `json:"id"`
	ChildOrder int64 `json:"child_order"`
}

type ReorderProjectRequest struct {
	ParentID *int64              `json:"parent_id"`
	Children []ReorderProjectMap `json:"children"`
}

func (r *CreateProjectRequest) Patch(project *Project) {
	if r.Content != nil {
		project.Content = *r.Content
	}

	if r.Description != nil {
		project.Description = *r.Description
	}

	if r.ParentID != nil {
		project.ParentID = r.ParentID
	}

	if r.ChildOrder != nil {
		project.ChildOrder = *r.ChildOrder
	}

	if r.Favorite != nil {
		project.Favorite = *r.Favorite
	}
}

func (r *UpdateProjectRequest) Patch(project *Project) {
	if r.Content != nil {
		project.Content = *r.Content
	}

	if r.Description != nil {
		project.Description = *r.Description
	}
}

func (m *ReorderProjectMap) Patch(project *Project) {
	project.ChildOrder = m.ChildOrder
}
