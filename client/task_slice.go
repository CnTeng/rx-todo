package client

import (
	"cmp"
	"fmt"
	"path/filepath"
	"slices"

	"github.com/CnTeng/rx-todo/internal/model"
)

type TaskSlice []*model.Task

func (ts TaskSlice) SortByID() TaskSlice {
	slices.SortStableFunc(ts, func(a, b *model.Task) int {
		return cmp.Compare(a.ID, b.ID)
	})
	return ts
}

func (ts TaskSlice) SortByName() TaskSlice {
	slices.SortStableFunc(ts, func(a, b *model.Task) int {
		return cmp.Compare(a.Name, b.Name)
	})
	return ts
}

func (ts TaskSlice) SortByPosition() TaskSlice {
	slices.SortStableFunc(ts, func(a, b *model.Task) int {
		return cmp.Compare(a.Position, b.Position)
	})
	return ts
}

func (ts TaskSlice) filter(condition func(*model.Task) bool) TaskSlice {
	var filtered TaskSlice
	for _, t := range ts {
		if condition(t) {
			filtered = append(filtered, t)
		}
	}
	return filtered
}

func (ts TaskSlice) FilterByName(pattern *string) TaskSlice {
	if pattern == nil {
		return ts
	}
	return ts.filter(func(t *model.Task) bool {
		matched, _ := filepath.Match(*pattern, t.Name)
		return matched
	})
}

func (ts TaskSlice) FilterByIDs(ids []int64) TaskSlice {
	return ts.filter(func(t *model.Task) bool {
		if ids == nil {
			return true
		}
		return slices.Contains(ids, t.ID)
	})
}

func (ts TaskSlice) FilterByPriority(priority model.Priority) TaskSlice {
	return ts.filter(func(t *model.Task) bool {
		return t.Priority == priority
	})
}

func (ts TaskSlice) FilterByLessPriority(priority model.Priority) TaskSlice {
	return ts.filter(func(t *model.Task) bool {
		return t.Priority <= priority
	})
}

func (ts TaskSlice) FilterByGreaterPriority(priority model.Priority) TaskSlice {
	return ts.filter(func(t *model.Task) bool {
		return t.Priority >= priority
	})
}

func (ts TaskSlice) FilterByProjectID(projectID int64) TaskSlice {
	return ts.filter(func(t *model.Task) bool {
		if t.ProjectID == nil {
			return false
		}
		return *t.ProjectID == projectID
	})
}

func (ts TaskSlice) GetIDs() []int64 {
	ids := make([]int64, len(ts))
	for i, t := range ts {
		ids[i] = t.ID
	}
	return ids
}

func (ts TaskSlice) GetIDByIndex(idx int) (int64, error) {
	if idx < 0 || idx >= len(ts) {
		return 0, fmt.Errorf("index out of range")
	}
	return ts[idx].ID, nil
}

func (ts TaskSlice) GetIDsByIndexs(idxs []int) ([]int64, error) {
	if idxs == nil {
		return ts.GetIDs(), nil
	}
	ids := make([]int64, len(idxs))
	for i, idx := range idxs {
		id, err := ts.GetIDByIndex(idx)
		if err != nil {
			return nil, err
		}
		ids[i] = id
	}
	return ids, nil
}

func (ts TaskSlice) GetTaskByIndex(idx int) (*model.Task, error) {
	if idx < 0 || idx >= len(ts) {
		return nil, fmt.Errorf("index out of range")
	}

	return ts[idx], nil
}

func (ts TaskSlice) GetProgress() {
	progressMap := make(map[int64]model.Progress)

	for _, t := range ts {
		if t.ParentID == nil {
			continue
		}

		progress := progressMap[*t.ParentID]
		progress.Total++

		if t.Done {
			progress.Done++
		}

		progressMap[*t.ParentID] = progress
	}

	for _, t := range ts {
		if progress, ok := progressMap[t.ID]; ok {
			t.Progress = &progress
		}
	}
}
