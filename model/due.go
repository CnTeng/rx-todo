package model

import (
	"fmt"
	"time"
)

type Due struct {
	Date      *time.Time `json:"date"`
	Recurring *bool      `json:"recurring"`
}

var timeLayouts = []string{
	time.RFC3339,
	time.DateTime,
	time.DateOnly,
	time.TimeOnly,
}

func ParseDueDate(date string) (*time.Time, error) {
	for _, layout := range timeLayouts {
		if t, err := time.Parse(layout, date); err == nil {
			return &t, nil
		}
	}

	return nil, fmt.Errorf("invalid due format")
}
