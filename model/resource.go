package model

import "time"

type Resource struct {
	ID int64 `json:"id"`

	// Meta fields
	Deleted   bool      `json:"deleted"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Syncable interface {
	GetID() int64
	GetUpdatedAt() time.Time
}

func (r Resource) GetUpdatedAt() time.Time {
	return r.UpdatedAt
}

func (r Resource) GetID() int64 {
	return r.ID
}

type SyncStatus struct {
	LabelSyncedAt    *time.Time `json:"label_synced_at,omitempty"`
	ProjectSyncedAt  *time.Time `json:"project_synced_at,omitempty"`
	ReminderSyncedAt *time.Time `json:"reminder_synced_at,omitempty"`
	TaskSyncedAt     *time.Time `json:"task_synced_at,omitempty"`
	UserSyncedAt     *time.Time `json:"user_synced_at,omitempty"`
}

type Resources struct {
	Labels     map[int64]*Label    `json:"labels"`
	Projects   map[int64]*Project  `json:"projects"`
	Reminders  map[int64]*Reminder `json:"reminders"`
	Tasks      map[int64]*Task     `json:"tasks"`
	User       *User               `json:"user,omitempty"`
	SyncStatus SyncStatus          `json:"sync_status,omitempty"`
}

func NewResources() *Resources {
	return &Resources{
		Labels:    make(map[int64]*Label),
		Projects:  make(map[int64]*Project),
		Reminders: make(map[int64]*Reminder),
		Tasks:     make(map[int64]*Task),
	}
}

type ResourceSyncRequest SyncStatus

type ResourceSyncResponse Resources

func NewResponse(args ...any) *ResourceSyncResponse {
	response := (*ResourceSyncResponse)(NewResources())

	for _, option := range args {
		switch v := option.(type) {
		case []*Label:
			for _, label := range v {
				response.Labels[label.ID] = label
			}
			response.SyncStatus.LabelSyncedAt = GetSyncStatus(v)
		case []*Project:
			for _, project := range v {
				response.Projects[project.ID] = project
			}
			response.SyncStatus.ProjectSyncedAt = GetSyncStatus(v)
		case []*Reminder:
			for _, reminder := range v {
				response.Reminders[reminder.ID] = reminder
			}
			response.SyncStatus.ReminderSyncedAt = GetSyncStatus(v)
		case []*Task:
			for _, task := range v {
				response.Tasks[task.ID] = task
			}
			response.SyncStatus.TaskSyncedAt = GetSyncStatus(v)
		case *User:
			response.User = v
			response.SyncStatus.UserSyncedAt = &v.UpdatedAt
		}
	}

	return response
}

func (r *Resources) GetLabels() []*Label {
	labels := make([]*Label, 0, len(r.Labels))

	for _, label := range r.Labels {
		labels = append(labels, label)
	}

	return labels
}

func (r *Resources) GetProjects() []*Project {
	projects := make([]*Project, 0, len(r.Projects))

	for _, project := range r.Projects {
		projects = append(projects, project)
	}

	return projects
}

func (r *Resources) GetReminders() []*Reminder {
	reminders := make([]*Reminder, 0, len(r.Reminders))

	for _, reminder := range r.Reminders {
		reminders = append(reminders, reminder)
	}

	return reminders
}

func (r *Resources) GetTasks() []*Task {
	tasks := make([]*Task, 0, len(r.Tasks))

	for _, task := range r.Tasks {
		tasks = append(tasks, task)
	}

	return tasks
}

func GetSyncStatus[T Syncable](res []T) *time.Time {
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
