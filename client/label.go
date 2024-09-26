package client

import (
	"encoding/json"

	"github.com/CnTeng/rx-todo/internal/model"
)

func (c *Client) CreateLabel(name string, color string) (*model.Label, error) {
	request := newRequest(
		c.Endpoint,
		c.Token,
		&model.LabelCreationRequest{Name: name, Color: getHexColor(color)},
	)

	response, err := request.withPath("labels").post()
	if err != nil {
		return nil, err
	}
	defer response.Close()

	var label model.Label
	if err := json.NewDecoder(response).Decode(&label); err != nil {
		return nil, err
	}

	return &label, nil
}

func (c *Client) UpdateLabel(id int64, name *string, color *string) (*model.Label, error) {
	request := newRequest(
		c.Endpoint,
		c.Token,
		&model.LabelUpdateRequest{Name: name, Color: color},
	)

	response, err := request.withPath(c.Endpoint).withID(id).put()
	if err != nil {
		return nil, err
	}
	defer response.Close()

	var label model.Label
	if err := json.NewDecoder(response).Decode(&label); err != nil {
		return nil, err
	}

	return &label, nil
}

func (c *Client) DeleteLabel(id int64) error {
	request := newRequest(c.Endpoint, c.Token, nil)

	_, err := request.withPath("labels").withID(id).delete()
	if err != nil {
		return err
	}

	return nil
}
