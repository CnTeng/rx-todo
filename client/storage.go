package client

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/CnTeng/rx-todo/internal/model"
	"github.com/adrg/xdg"
)

type storage struct {
	Path string
}

func NewStorage(path string) (*storage, error) {
	path, err := xdg.CacheFile(path)
	if err != nil {
		return nil, err
	}

	return &storage{Path: path}, nil
}

func (c *Client) Store(resource *model.ResourceSyncResponse) error {
	file, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(c.StoragePath, file, 0o644); err != nil {
		return err
	}

	return nil
}

func pathLabel(res *model.ResourceSyncResponse, labels []*model.Label) {
	labelMap := make(map[int64]*model.Label)
	for _, rl := range *res.Labels {
		labelMap[rl.ID] = rl
	}

	for _, l := range labels {
		labelMap[l.ID] = l
	}

	updatedLabels := make([]*model.Label, 0, len(labelMap))
	for _, label := range labelMap {
		updatedLabels = append(updatedLabels, label)
	}

	*res.Labels = updatedLabels
}

func (c *Client) Patch(res any) (*model.ResourceSyncResponse, error) {
	resources := &model.ResourceSyncResponse{}
	file, err := os.ReadFile(c.StoragePath)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(file, resources); err != nil {
		return nil, err
	}

	switch res := res.(type) {
	case []*model.Label:
		pathLabel(resources, res)
	default:
		return nil, fmt.Errorf("unsupported type")
	}

	return resources, c.Store(resources)
}
