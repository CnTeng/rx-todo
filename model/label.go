package model

import (
	"fmt"
)

// Label represents a label
type Label struct {
	resource
	UserID int64  `json:"user_id"`
	Name   string `json:"name"`
	Color  string `json:"color"`
}

// LabelCreationRequest represents a request to create a label
type LabelCreationRequest struct {
	Name  *string `json:"name" validate:"required"`
	Color *string `json:"color" validate:"required,hexcolor"`
}

// LabelUpdateRequest represents a request to update a label
type LabelUpdateRequest struct {
	Name  *string `json:"name" validate:"omitnil,notempty"`
	Color *string `json:"color" validate:"omitnil,hexcolor"`
}

type LabelDeleteRequestWithID struct {
	ID int64 `json:"id" validate:"required,notempty"`
}

// LabelUpdateRequestWithID represents a request to update a label with an ID
type LabelUpdateRequestWithID struct {
	ID *int64 `json:"id" validate:"required,notempty"`
	*LabelUpdateRequest
}

// LabelUpdateRequest at least needs one of name or color to be set
func (r *LabelUpdateRequest) validate() error {
	if (r.Name == nil) && (r.Color == nil) {
		return fmt.Errorf("validate: at least one of name or color should be set")
	}

	return nil
}

func (r *LabelCreationRequest) Patch(label *Label) {
	if r.Name != nil {
		label.Name = *r.Name
	}

	if r.Color != nil {
		label.Color = *r.Color
	}
}

func (r *LabelUpdateRequest) Patch(label *Label) {
	if r.Name != nil {
		label.Name = *r.Name
	}

	if r.Color != nil {
		label.Color = *r.Color
	}
}
