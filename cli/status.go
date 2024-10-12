package cli

import (
	"github.com/CnTeng/rx-todo/model"
	"github.com/fatih/color"
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

func (sm *statusMap) getStatusIcon(id int64, iconType iconType) string {
	s := (*sm)[id]

	icons := getIcons(iconType)

	switch s {
	case Add:
		return color.New(color.FgGreen).Sprint(icons.add)
	case Change:
		return color.New(color.FgYellow).Sprint(icons.change)
	case Delete:
		return color.New(color.FgRed).Sprint(icons.delete)
	default:
		return icons.none
	}
}
