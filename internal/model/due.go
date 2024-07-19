package model

import "time"

type Due struct {
	Date      *time.Time `json:"date"`
	Recurring *bool      `json:"recurring"`
}
