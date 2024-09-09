package model

import "time"

// Reminder represents a reminder
type Reminder struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"user_id"`
	TaskID int64 `json:"task_id"`
	Due    Due   `json:"due"`

	// Meta fields
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ReminderCreationRequest represents a request to create a reminder
type ReminderCreationRequest struct {
	TaskID int64 `json:"task_id"`
	Due    Due   `json:"due"`
}

// ReminderUpdateRequest represents a request to update a reminder
type ReminderUpdateRequest struct {
	Due Due `json:"due"`
}

func (r *ReminderCreationRequest) Patch(userID int64, reminder *Reminder) {
	reminder.UserID = userID
	reminder.TaskID = r.TaskID
	reminder.Due = r.Due
}

func (r *ReminderUpdateRequest) Patch(reminder *Reminder) {
	reminder.Due = r.Due
}
