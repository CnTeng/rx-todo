package client

import (
	"encoding/json"

	"github.com/CnTeng/rx-todo/model"
)

func (c *Client) CreateLabel(r *model.LabelCreationRequest) (*model.Label, error) {
	request := newRequest(c.Endpoint, c.Token, r)

	response, err := request.withPath("labels").post()
	if err != nil {
		return nil, err
	}
	defer response.Close()

	label := &model.Label{}
	if err := json.NewDecoder(response).Decode(label); err != nil {
		return nil, err
	}

	if err := c.patch([]*model.Label{label}); err != nil {
		return nil, err
	}

	return label, nil
}

func (c *Client) GetLabels() []*model.Label {
	labels := c.storage.GetLabels()

	return labels
}

func (c *Client) UpdateLabel(r *model.LabelUpdateRequestWithID) (*model.Label, error) {
	request := newRequest(c.Endpoint, c.Token, r.LabelUpdateRequest)

	response, err := request.withPath("labels").withID(*r.ID).put()
	if err != nil {
		return nil, err
	}
	defer response.Close()

	label := &model.Label{}
	if err := json.NewDecoder(response).Decode(label); err != nil {
		return nil, err
	}

	if err := c.patch([]*model.Label{label}); err != nil {
		return nil, err
	}

	return label, nil
}

func (c *Client) DeleteLabel(id int64) (*model.Label, error) {
	request := newRequest(c.Endpoint, c.Token, nil)

	if _, err := request.withPath("labels").withID(id).delete(); err != nil {
		return nil, err
	}

	label := c.GetLabel(id)
	label.Deleted = true

	if err := c.patch([]*model.Label{label}); err != nil {
		return nil, err
	}

	return label, nil
}
