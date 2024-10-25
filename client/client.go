package client

import (
	"encoding/json"
	"fmt"

	"github.com/CnTeng/rx-todo/internal/model"
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
	response, err := NewRequest(c.Endpoint, c.Token).
		WithPath("resources/sync").
		WithBody(&model.ResourceSyncRequest{}).
		Post()
	if err != nil {
		return err
	}
	defer response.Close()

	if err := json.NewDecoder(response).Decode(&c.resources); err != nil {
		return err
	}

	fmt.Println(c.resources)

	c.GetProjectProgress()
	c.GetTaskProgress()

	return c.save()
}

func (c *Client) patch(res any) error {
	if err := c.load(); err != nil {
		return err
	}

	c.storage.patch(res)

	return c.save()
}
