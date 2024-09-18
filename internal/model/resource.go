package model

import "time"

type Resource struct {
	ID int64 `json:"id"`

	// Meta fields
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Syncable interface {
	GetUpdatedAt() time.Time
}

func (r Resource) GetUpdatedAt() time.Time {
	return r.UpdatedAt
}

type SyncStatus struct {
	LabelSyncedAt    *time.Time `json:"label_synced_at,omitempty"`
	ProjectSyncedAt  *time.Time `json:"project_synced_at,omitempty"`
	ReminderSyncedAt *time.Time `json:"reminder_synced_at,omitempty"`
	TaskSyncedAt     *time.Time `json:"task_synced_at,omitempty"`
	UserSyncedAt     *time.Time `json:"user_synced_at,omitempty"`
}

type ResourceSyncRequest SyncStatus

type ResourceSyncResponse struct {
	Labels     *[]*Label    `json:"labels,omitempty"`
	Projects   *[]*Project  `json:"projects,omitempty"`
	Reminders  *[]*Reminder `json:"reminders,omitempty"`
	Tasks      *[]*Task     `json:"tasks,omitempty"`
	User       *User        `json:"user,omitempty"`
	SyncStatus SyncStatus   `json:"sync_status,omitempty"`
}

func NewResponse(args ...any) *ResourceSyncResponse {
	response := &ResourceSyncResponse{}

	for _, option := range args {
		switch v := option.(type) {
		case []*Label:
			response.Labels = &v
			response.SyncStatus.LabelSyncedAt = GetSyncStatus(v)
		case []*Project:
			response.Projects = &v
			response.SyncStatus.ProjectSyncedAt = GetSyncStatus(v)
		case []*Reminder:
			response.Reminders = &v
			response.SyncStatus.ReminderSyncedAt = GetSyncStatus(v)
		case []*Task:
			response.Tasks = &v
			response.SyncStatus.TaskSyncedAt = GetSyncStatus(v)
		case *User:
			response.User = v
			response.SyncStatus.UserSyncedAt = &v.UpdatedAt
		}
	}

	return response
}

func GetSyncStatus[T Syncable](resources []T) *time.Time {
	if len(resources) == 0 {
		return nil
	}
	time := time.Time{}

	for _, resource := range resources {
		if resource.GetUpdatedAt().After(time) {
			time = resource.GetUpdatedAt()
		}
	}

	return &time
}
