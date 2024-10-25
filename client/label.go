package client

import (
	"encoding/json"

	"github.com/CnTeng/rx-todo/model"
)

func (c *Client) CreateLabel(r *model.LabelCreationRequest) (*model.Label, error) {
	response, err := NewRequest(c.Endpoint, c.Token).
		WithPath("labels").
		WithBody(r).
		Post()
	if err != nil {
		return nil, err
	}
	defer response.Close()

	label := &model.Label{}
	if err := json.NewDecoder(response).Decode(label); err != nil {
		return nil, err
	}

	if err := c.patch(label); err != nil {
		return nil, err
	}

	return label, nil
}

func (c *Client) GetLabel(id int64) *model.Label {
	label := c.storage.GetLabel(id)

	return label
}

func (c *Client) GetLabels() []*model.Label {
	labels := c.storage.GetLabels()

	return labels
}

func (c *Client) UpdateLabel(r *model.LabelUpdateRequestWithID) (*model.Label, error) {
	response, err := NewRequest(c.Endpoint, c.Token).
		WithPath("labels").
		WithID(*r.ID).
		WithBody(r.LabelUpdateRequest).
		Put()
	if err != nil {
		return nil, err
	}
	defer response.Close()

	label := &model.Label{}
	if err := json.NewDecoder(response).Decode(label); err != nil {
		return nil, err
	}

	if err := c.patch(label); err != nil {
		return nil, err
	}

	return label, nil
}

func (c *Client) DeleteLabel(id int64) (*model.Label, error) {
	if _, err := NewRequest(c.Endpoint, c.Token).
		WithPath("labels").
		WithID(id).
		Delete(); err != nil {
		return nil, err
	}

	label := c.GetLabel(id)
	label.Deleted = true

	if err := c.patch(label); err != nil {
		return nil, err
	}

	return label, nil
}
