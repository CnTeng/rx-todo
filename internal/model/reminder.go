package model

import "time"

type Reminder struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"user_id"`
	TaskID int64 `json:"task_id"`
	Due    Due   `json:"due"`

	// Meta fields
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateReminderRequest struct {
	UserID int64 `json:"user_id"`
	TaskID int64 `json:"task_id"`
	Due    Due   `json:"due"`
}

type UpdateReminderRequest struct {
	Due Due `json:"due"`
}

func (r *CreateReminderRequest) Patch(userID int64, reminder *Reminder) {
	reminder.UserID = userID
	reminder.TaskID = r.TaskID
	reminder.Due = r.Due
}

func (r *UpdateReminderRequest) Patch(userID int64, reminder *Reminder) {
	reminder.Due = r.Due
}
