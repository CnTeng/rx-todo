package model

import (
	"fmt"
	"time"
)

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

// LabelCreationRequest represents a request to create a label
type LabelCreationRequest struct {
	Name  string `json:"name" validate:"required"`
	Color string `json:"color" validate:"required,hexcolor"`
}

// LabelUpdateRequest represents a request to update a label
type LabelUpdateRequest struct {
	Name  *string `json:"name" validate:"omitnil,notempty"`
	Color *string `json:"color" validate:"omitnil,hexcolor"`
}

// LabelUpdateRequest at least needs one of name or color to be set
func (r *LabelUpdateRequest) validate() error {
	if (r.Name == nil) && (r.Color == nil) {
		return fmt.Errorf("validate: at least one of name or color should be set")
	}

	return nil
}

func (r *LabelCreationRequest) Patch(label *Label) {
	label.Name = r.Name
	label.Color = r.Color
}

func (r *LabelUpdateRequest) Patch(label *Label) {
	if r.Name != nil {
		label.Name = *r.Name
	}

	if r.Color != nil {
		label.Color = *r.Color
	}
}
