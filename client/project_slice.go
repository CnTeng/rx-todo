package client

import (
	"cmp"
	"fmt"
	"path/filepath"
	"slices"

	"github.com/CnTeng/rx-todo/internal/model"
)

type ProjectSlice []*model.Project

func (ps ProjectSlice) SortByID() ProjectSlice {
	slices.SortStableFunc(ps, func(a, b *model.Project) int {
		return cmp.Compare(a.ID, b.ID)
	})
	return ps
}

func (ps ProjectSlice) SortByName() ProjectSlice {
	slices.SortStableFunc(ps, func(a, b *model.Project) int {
		return cmp.Compare(a.Name, b.Name)
	})
	return ps
}

func (ps ProjectSlice) SortByPosition() ProjectSlice {
	slices.SortStableFunc(ps, func(a, b *model.Project) int {
		return cmp.Compare(a.Position, b.Position)
	})
	return ps
}

func (ps ProjectSlice) filter(condition func(*model.Project) bool) ProjectSlice {
	var filtered ProjectSlice
	for _, p := range ps {
		if condition(p) {
			filtered = append(filtered, p)
		}
	}
	return filtered
}

func (ps ProjectSlice) FilterByName(pattern string) ProjectSlice {
	return ps.filter(func(p *model.Project) bool {
		matched, _ := filepath.Match(pattern, p.Name)
		return matched
	})
}

func (ps ProjectSlice) FilterByIDs(ids []int64) ProjectSlice {
	return ps.filter(func(p *model.Project) bool {
		if ids == nil {
			return true
		}
		return slices.Contains(ids, p.ID)
	})
}

func (ps ProjectSlice) GetIDs() []int64 {
	var ids []int64
	for _, t := range ps {
		ids = append(ids, t.ID)
	}
	return ids
}

func (ps ProjectSlice) GetIDByIndex(idx int) (int64, error) {
	if idx < 0 || idx >= len(ps) {
		return 0, fmt.Errorf("index out of range")
	}
	return ps[idx].ID, nil
}

func (ps ProjectSlice) GetIDsByIndexs(idxs []int) ([]int64, error) {
	if idxs == nil {
		return ps.GetIDs(), nil
	}
	ids := make([]int64, len(idxs))
	for i, idx := range idxs {
		id, err := ps.GetIDByIndex(idx)
		if err != nil {
			return nil, err
		}
		ids[i] = id
	}
	return ids, nil
}

func (ps ProjectSlice) GetProjectByIndex(idx int) (*model.Project, error) {
	if idx < 0 || idx >= len(ps) {
		return nil, fmt.Errorf("index out of range")
	}

	return ps[idx], nil
}
