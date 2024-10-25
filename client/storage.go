package client

import (
	"encoding/json"
	"os"

	"github.com/CnTeng/rx-todo/model"
	"github.com/adrg/xdg"
)

type storage struct {
	Path string
	*resources
}

func newStorage(path string) (*storage, error) {
	path, err := xdg.CacheFile(path)
	if err != nil {
		return nil, err
	}

	return &storage{
		Path:      path,
		resources: newResources(),
	}, nil
}

func (s *storage) load() error {
	file, err := os.ReadFile(s.Path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(file, s); err != nil {
		return err
	}

	return nil
}

func (s *storage) save() error {
	file, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(s.Path, file, 0o644); err != nil {
		return err
	}

	return nil
}

func (s *storage) patch(res any) {
	switch res := res.(type) {
	case *model.Label:
		s.Labels[res.ID] = res
	case []*model.Label:
		for _, label := range res {
			s.Labels[label.ID] = label
		}
	case *model.Project:
		s.Projects[res.ID] = res
	case []*model.Project:
		for _, project := range res {
			s.Projects[project.ID] = project
		}
	case *model.Reminder:
		s.Reminders[res.ID] = res
	case []*model.Reminder:
		for _, reminder := range res {
			s.Reminders[reminder.ID] = reminder
		}
	case *model.Task:
		if res.Deleted {
			delete(s.Tasks, res.ID)
		} else {
			s.Tasks[res.ID] = res
		}
	case []*model.Task:
		for _, task := range res {
			if task.Deleted {
				delete(s.Tasks, task.ID)
			} else {
				s.Tasks[task.ID] = task
			}
		}
	}
}
