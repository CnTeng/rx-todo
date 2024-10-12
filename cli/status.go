package cli

import (
	"github.com/CnTeng/rx-todo/model"
)

type status int

const (
	None status = iota
	Add
	Change
	Delete
)

type statusMap map[int64]status

func NewStatusMap[T model.Resource](res []T, status status) *statusMap {
	sm := make(statusMap, len(res))
	for _, r := range res {
		sm[r.GetID()] = status
	}
	return &sm
}

func (sm *statusMap) getStatusIcon(id int64, icons *icons) string {
	s := (*sm)[id]

	switch s {
	case Add:
		return icons.add
	case Change:
		return icons.change
	case Delete:
		return icons.delete
	default:
		return icons.none
	}
}
