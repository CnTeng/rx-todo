package model

// Reminder represents a reminder
type Reminder struct {
	Resource
	UserID int64 `json:"user_id"`
	TaskID int64 `json:"task_id"`
	Due    Due   `json:"due"`
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

func (r *ReminderCreationRequest) Patch(reminder *Reminder) {
	reminder.TaskID = r.TaskID
	reminder.Due = r.Due
}

func (r *ReminderUpdateRequest) Patch(reminder *Reminder) {
	reminder.Due = r.Due
}
