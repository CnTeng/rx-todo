package client

import (
	"encoding/json"

	"github.com/CnTeng/rx-todo/internal/model"
)

type Client struct {
	Endpoint    string
	Token       string
	StoragePath string
}

func NewClient(server, token, storage string) *Client {
	return &Client{
		Endpoint:    server,
		Token:       token,
		StoragePath: storage,
	}
}

func (c *Client) Sync() (*model.ResourceSyncResponse, error) {
	request := newRequest(
		c.Endpoint+"/resources/sync",
		c.Token,
		&model.ResourceSyncRequest{})

	response, err := request.post()
	if err != nil {
		return nil, err
	}
	defer response.Close()

	var resource model.ResourceSyncResponse
	if err := json.NewDecoder(response).Decode(&resource); err != nil {
		return nil, err
	}

	return &resource, nil
}
