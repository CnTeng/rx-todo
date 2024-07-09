package model

import "time"

type Reminder struct {
	ID       int64     `json:"id"`
	TaskID   int64     `json:"task_id"`
	Reminder time.Time `json:"reminder_time"`
	Created  time.Time `json:"created_at"`
}
