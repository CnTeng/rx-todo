package client

import (
	"encoding/json"

	"github.com/CnTeng/rx-todo/model"
)

type Client struct {
	Endpoint string
	Token    string
	*storage
}

func NewClient(server, token, storagePath string) (*Client, error) {
	storage, err := newStorage(storagePath)
	if err != nil {
		return nil, err
	}

	return &Client{
		Endpoint: server,
		Token:    token,
		storage:  storage,
	}, nil
}

func (c *Client) Sync() error {
	request := newRequest(
		c.Endpoint,
		c.Token,
		&model.ResourceSyncRequest{},
	)

	response, err := request.withPath("resources/sync").post()
	if err != nil {
		return err
	}
	defer response.Close()

	if err := json.NewDecoder(response).Decode(&c.Resources); err != nil {
		return err
	}

	return c.save()
}

func (c *Client) patch(res any) error {
	if err := c.load(); err != nil {
		return err
	}

	c.storage.patch(res)

	return c.save()
}
