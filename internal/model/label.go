package model

import "time"

// Label represents a label
type Label struct {
	ID     int64  `json:"id"`
	UserID int64  `json:"user_id"`
	Name   string `json:"name"`
	Color  string `json:"color"`

	// Meta fields
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateLabelRequest struct {
	Name  *string `json:"name"`
	Color *string `json:"color"`
}

func (r *UpdateLabelRequest) Patch(label *Label) {
	if r.Name != nil {
		label.Name = *r.Name
	}
	if r.Color != nil {
		label.Color = *r.Color
	}
}
