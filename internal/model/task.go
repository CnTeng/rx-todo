package model

import "time"

type Due struct {
	Data      time.Time `json:"data"`
	Recurring bool      `json:"recurring"`
}

type Task struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Content     string    `json:"content"`
	Description string    `json:"description"`
	Due         Due       `json:"due"`
	Duration    int       `json:"duration"`
	Priority    int       `json:"priority"`
	ProjectID   int64     `json:"project_id"`
	ParentID    int64     `json:"parent_id"`
	ChildOrder  int       `json:"child_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Done        bool      `json:"done"`
	DoneAt      time.Time `json:"done_at"`
	Archive     bool      `json:"archive"`
	ArchiveAt   time.Time `json:"archive_at"`
}

type TaskAddRequest struct {
	Content     string `json:"content"`
	Description string `json:"description"`
	Due         Due    `json:"due"`
	Duration    int    `json:"duration"`
	Priority    int    `json:"priority"`
	ProjectID   int64  `json:"project_id"`
	ParentID    int64  `json:"parent_id"`
	ChildOrder  int    `json:"child_order"`
}

type TaskUpdateRequest struct {
	ID          int64  `json:"id"`
	Content     string `json:"content"`
	Description string `json:"description"`
	Due         Due    `json:"due"`
	Duration    int    `json:"duration"`
	Priority    int    `json:"priority"`
}

type TaskMoveRequest struct {
	ID        int64 `json:"id"`
	ProjectID int64 `json:"project_id"`
	ParentID  int64 `json:"parent_id"`
}

type TaskReorderRequest struct {
	Tasks []struct {
		ID         int64 `json:"id"`
		ChildOrder int   `json:"child_order"`
	} `json:"tasks"`
}

type TaskDeleteRequest struct {
	ID int64 `json:"id"`
}

type TaskDoneRequest struct {
	ID int64 `json:"id"`
}

type TaskUnDoneRequest struct {
	ID int64 `json:"id"`
}

type TaskArchiveRequest struct {
	ID int64 `json:"id"`
}

type TaskUnArchiveRequest struct {
	ID int64 `json:"id"`
}

type TaskGetRequest struct {
	ID int64 `json:"id"`
}

func (t *Task) PatchUpdateRequest(r *TaskUpdateRequest) {
	if r.Content == "" {
		r.Content = t.Content
	}

	if r.Description == "" {
		r.Description = t.Description
	}

	if r.Due.Data.IsZero() {
		r.Due = t.Due
	}

	if r.Duration == 0 {
		r.Duration = t.Duration
	}

	if r.Priority == 0 {
		r.Priority = t.Priority
	}
}

func (t *Task) PatchMoveRequest(r *TaskMoveRequest) {
	if r.ProjectID == 0 {
		r.ProjectID = t.ProjectID
	}

	if r.ParentID == 0 {
		r.ParentID = t.ParentID
	}
}
