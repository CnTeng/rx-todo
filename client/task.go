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
