package client

import (
	"encoding/json"

	"github.com/CnTeng/rx-todo/internal/model"
)

type TaskUpdateRequestWithID struct {
	ID int64 `json:"id"`
	model.TaskUpdateRequest
}

type TaskMoveRequestWithID struct {
	ID int64 `json:"id"`
	model.TaskMoveRequest
}

func (c *Client) CreateTask(r *model.TaskCreationRequest) (*model.Task, error) {
	response, err := NewRequest(c.Endpoint, c.Token).
		WithPath("tasks").
		WithBody(r).
		Post()
	if err != nil {
		return nil, err
	}
	defer response.Close()

	task := &model.Task{}
	if err := json.NewDecoder(response).Decode(task); err != nil {
		return nil, err
	}

	if err := c.patch(task); err != nil {
		return nil, err
	}

	return task, nil
}

func (c *Client) GetTask(id int64) *model.Task {
	return c.storage.GetTask(id)
}

func (c *Client) GetTasks() TaskSlice {
	return c.storage.GetTasks()
}

func (c *Client) UpdateTask(r *TaskUpdateRequestWithID) (*model.Task, error) {
	response, err := NewRequest(c.Endpoint, c.Token).
		WithPath("tasks").
		WithID(r.ID).
		WithBody(r.TaskUpdateRequest).
		Put()
	if err != nil {
		return nil, err
	}
	defer response.Close()

	task := &model.Task{}
	if err := json.NewDecoder(response).Decode(task); err != nil {
		return nil, err
	}

	if err := c.patch(task); err != nil {
		return nil, err
	}

	return task, nil
}

func (c *Client) MoveTask(r *TaskMoveRequestWithID) (*model.Task, error) {
	response, err := NewRequest(c.Endpoint, c.Token).
		WithPath("tasks").
		WithID(r.ID).
		WithPath("move").
		WithParameter("previous_id", r.PreviousID).
		WithParameter("project_id", r.ProjectID).
		WithParameter("parent_id", r.ParentID).
		WithBody(r).
		Put()
	if err != nil {
		return nil, err
	}
	defer response.Close()

	task := &model.Task{}
	if err := json.NewDecoder(response).Decode(task); err != nil {
		return nil, err
	}

	if err := c.patch(task); err != nil {
		return nil, err
	}

	return task, nil
}

func (c *Client) operateTask(id int64, operation string) (*model.Task, error) {
	response, err := NewRequest(c.Endpoint, c.Token).
		WithPath("tasks").
		WithID(id).
		WithPath(operation).
		Put()
	if err != nil {
		return nil, err
	}
	defer response.Close()

	task := &model.Task{}
	if err := json.NewDecoder(response).Decode(task); err != nil {
		return nil, err
	}

	if err := c.patch(task); err != nil {
		return nil, err
	}

	return task, nil
}

func (c *Client) OpenTask(id int64) (*model.Task, error) {
	return c.operateTask(id, "open")
}

func (c *Client) CloseTask(id int64) (*model.Task, error) {
	return c.operateTask(id, "close")
}

func (c *Client) ArchiveTask(id int64) (*model.Task, error) {
	return c.operateTask(id, "archive")
}

func (c *Client) UnarchiveTask(id int64) (*model.Task, error) {
	return c.operateTask(id, "unarchive")
}

func (c *Client) DeleteTask(id int64) (*model.Task, error) {
	if _, err := NewRequest(c.Endpoint, c.Token).
		WithPath("tasks").
		WithID(id).
		Delete(); err != nil {
		return nil, err
	}

	task := c.GetTask(id)
	task.Deleted = true

	if err := c.patch(task); err != nil {
		return nil, err
	}

	return task, nil
}
