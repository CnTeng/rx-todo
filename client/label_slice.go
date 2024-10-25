package client

import (
	"cmp"
	"slices"

	"github.com/CnTeng/rx-todo/model"
)

type LabelSlice []*model.Label

func (ls LabelSlice) SortByID() LabelSlice {
	slices.SortStableFunc(ls, func(a, b *model.Label) int {
		return cmp.Compare(a.ID, b.ID)
	})
	return ls
}

func (ls LabelSlice) SortByName() LabelSlice {
	slices.SortStableFunc(ls, func(a, b *model.Label) int {
		return cmp.Compare(a.Name, b.Name)
	})
	return ls
}
