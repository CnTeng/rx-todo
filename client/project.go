package client

import (
	"encoding/json"

	"github.com/CnTeng/rx-todo/internal/model"
)

type ProjectUpdateRequestWithID struct {
	ID int64 `json:"id"`
	model.ProjectUpdateRequest
}

type ProjectMoveRequestWithID struct {
	ID         int64  `json:"id"`
	PreviousID *int64 `json:"previous_id"`
}

func (c *Client) CreateProject(r *model.ProjectCreationRequest) (*model.Project, error) {
	response, err := NewRequest(c.Endpoint, c.Token).
		WithPath("projects").
		WithBody(r).
		Post()
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

func (c *Client) GetProject(id int64) *model.Project {
	return c.storage.GetProject(id)
}

func (c *Client) GetProjects() []*model.Project {
	return c.storage.GetProjects()
}

func (c *Client) UpdateProject(r *ProjectUpdateRequestWithID) (*model.Project, error) {
	response, err := NewRequest(c.Endpoint, c.Token).
		WithPath("projects").
		WithID(r.ID).
		WithBody(r.ProjectUpdateRequest).
		Put()
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

func (c *Client) MoveProject(r *ProjectMoveRequestWithID) (*model.Project, error) {
	response, err := NewRequest(c.Endpoint, c.Token).
		WithPath("projects").
		WithID(r.ID).
		WithParameter("previous_id", r.PreviousID).
		Put()
	if err != nil {
		return nil, err
	}
	defer response.Close()

	project := &model.Project{}
	if err := json.NewDecoder(response).Decode(project); err != nil {
		return nil, err
	}

	if err := c.patch(project); err != nil {
		return nil, err
	}

	return project, nil
}

func (c *Client) operateProject(id int64, operation string) (*model.Project, error) {
	response, err := NewRequest(c.Endpoint, c.Token).
		WithPath("projects").
		WithID(id).
		WithPath(operation).
		Put()
	if err != nil {
		return nil, err
	}
	defer response.Close()

	project := &model.Project{}
	if err := json.NewDecoder(response).Decode(project); err != nil {
		return nil, err
	}

	if err := c.patch(project); err != nil {
		return nil, err
	}

	return project, nil
}

func (c *Client) ArchiveProject(id int64) (*model.Project, error) {
	return c.operateProject(id, "archive")
}

func (c *Client) UnarchiveProject(id int64) (*model.Project, error) {
	return c.operateProject(id, "unarchive")
}

func (c *Client) DeleteProject(id int64) (*model.Project, error) {
	if _, err := NewRequest(c.Endpoint, c.Token).
		WithPath("projects").
		WithID(id).
		Delete(); err != nil {
		return nil, err
	}

	project := c.GetProject(id)
	project.Deleted = true

	if err := c.patch(project); err != nil {
		return nil, err
	}

	return project, nil
}
