package model

import "time"

// Resource represents a generic struct for all resources
type resource struct {
	ID        int64     `json:"id"`
	Deleted   bool      `json:"deleted"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Resource interface {
	GetID() int64
	GetUpdatedAt() time.Time
}

func (r *resource) GetUpdatedAt() time.Time {
	return r.UpdatedAt
}

func (r *resource) GetID() int64 {
	return r.ID
}

func getSyncStatus[T Resource](res []T) *time.Time {
	if len(res) == 0 {
		return nil
	}
	time := time.Time{}

	for _, resource := range res {
		if resource.GetUpdatedAt().After(time) {
			time = resource.GetUpdatedAt()
		}
	}

	return &time
}

type ResourceSyncRequest struct {
	LabelSyncedAt    *time.Time `json:"label_synced_at,omitempty"`
	ProjectSyncedAt  *time.Time `json:"project_synced_at,omitempty"`
	ReminderSyncedAt *time.Time `json:"reminder_synced_at,omitempty"`
	TaskSyncedAt     *time.Time `json:"task_synced_at,omitempty"`
	UserSyncedAt     *time.Time `json:"user_synced_at,omitempty"`
}

type ResourceSyncResponse struct {
	Labels     map[int64]*Label    `json:"labels"`
	Projects   map[int64]*Project  `json:"projects"`
	Reminders  map[int64]*Reminder `json:"reminders"`
	Tasks      map[int64]*Task     `json:"tasks"`
	User       *User               `json:"user,omitempty"`
	SyncStatus ResourceSyncRequest `json:"sync_status,omitempty"`
}

func NewResourceSyncResponse(args ...any) *ResourceSyncResponse {
	response := &ResourceSyncResponse{
		Labels:    make(map[int64]*Label),
		Projects:  make(map[int64]*Project),
		Reminders: make(map[int64]*Reminder),
		Tasks:     make(map[int64]*Task),
	}

	for _, option := range args {
		switch v := option.(type) {
		case []*Label:
			for _, label := range v {
				response.Labels[label.ID] = label
			}
			response.SyncStatus.LabelSyncedAt = getSyncStatus(v)
		case []*Project:
			for _, project := range v {
				response.Projects[project.ID] = project
			}
			response.SyncStatus.ProjectSyncedAt = getSyncStatus(v)
		case []*Reminder:
			for _, reminder := range v {
				response.Reminders[reminder.ID] = reminder
			}
			response.SyncStatus.ReminderSyncedAt = getSyncStatus(v)
		case []*Task:
			for _, task := range v {
				response.Tasks[task.ID] = task
			}
			response.SyncStatus.TaskSyncedAt = getSyncStatus(v)
		case *User:
			response.User = v
			response.SyncStatus.UserSyncedAt = &v.UpdatedAt
		}
	}

	return response
}
