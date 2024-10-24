package client

import (
	"encoding/json"

	"github.com/CnTeng/rx-todo/model"
)

func (c *Client) CreateProject(r *model.ProjectCreationRequest) (*model.Project, error) {
	request := newRequest(c.Endpoint, c.Token, r)

	response, err := request.withPath("projects").post()
	if err != nil {
		return nil, err
	}
	defer response.Close()

	project := &model.Project{}
	if err := json.NewDecoder(response).Decode(project); err != nil {
		return nil, err
	}

	if err := c.patch([]*model.Project{project}); err != nil {
		return nil, err
	}

	return project, nil
}
