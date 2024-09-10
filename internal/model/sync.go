package model

import "time"

type (
	ObjectType string
	Operation  string
)

const (
	DeleteOperation Operation = "delete"
	CreateOperation Operation = "create"
	UpdateOperation Operation = "update"
)

type SyncStatus struct {
	ID         int64      `json:"id"`
	UserID     int64      `json:"user_id"`
	ObjectIDs  []int64    `json:"object_id"`
	ObjectType ObjectType `json:"object_type"`
	Operation  Operation  `json:"operation"`

	// Meta fields
	CreatedAt time.Time `json:"created_at"`
}

type SyncRequest struct{}

type SyncResponse struct {
	CompletedInfo map[string]string     `json:"completed_info"`
	Labels        map[string][]Label    `json:"labels"`
	Projects      map[string][]Project  `json:"projects"`
	Reminders     map[string][]Reminder `json:"reminders"`
	Tasks         map[string][]Task     `json:"tasks"`
	User          User                  `json:"user"`
}

type SyncObject interface {
	ToSyncStatus(opt Operation) SyncStatus
}

func (s *SyncStatus) ClearObjectIDs() *SyncStatus {
	s.ObjectIDs = []int64{}
	return s
}
