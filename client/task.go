package client

import (
	"encoding/json"

	"github.com/CnTeng/rx-todo/model"
)

func (c *Client) CreateTask(r *model.TaskCreationRequest) (*model.Task, error) {
	request := newRequest(c.Endpoint, c.Token, r)

	response, err := request.withPath("tasks").post()
	if err != nil {
		return nil, err
	}
	defer response.Close()

	task := &model.Task{}
	if err := json.NewDecoder(response).Decode(task); err != nil {
		return nil, err
	}

	if err := c.patch([]*model.Task{task}); err != nil {
		return nil, err
	}

	return task, nil
}

func (c *Client) GetTasks() []*model.Task {
	tasks := c.storage.GetTasks()

	return tasks
}

func (c *Client) UpdateTask(r *model.TaskUpdateRequestWithID) (*model.Task, error) {
	request := newRequest(c.Endpoint, c.Token, r.TaskUpdateRequest)

	response, err := request.withPath("tasks").withID(r.ID).put()
	if err != nil {
		return nil, err
	}
	defer response.Close()

	task := &model.Task{}
	if err := json.NewDecoder(response).Decode(task); err != nil {
		return nil, err
	}

	if err := c.patch([]*model.Task{task}); err != nil {
		return nil, err
	}

	return task, nil
}

func (c *Client) DeleteTask(id int64) (*model.Task, error) {
	request := newRequest(c.Endpoint, c.Token, nil)

	if _, err := request.withPath("tasks").withID(id).delete(); err != nil {
		return nil, err
	}

	task := c.GetTask(id)
	task.Deleted = true

	if err := c.patch([]*model.Task{task}); err != nil {
		return nil, err
	}

	return task, nil
}
