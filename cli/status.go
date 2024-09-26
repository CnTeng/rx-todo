package cli

import (
	"github.com/CnTeng/rx-todo/internal/model"
	"github.com/fatih/color"
)

type status int

const (
	None status = iota
	Add
	Change
	Delete
)

type StatusMap map[int64]status

func (s *status) String() string {
	switch *s {
	case Add:
		return color.New(color.FgGreen).Sprint("+")
	case Change:
		return color.New(color.FgYellow).Sprint("~")
	case Delete:
		return color.New(color.FgRed).Sprint("-")
	default:
		return " "
	}
}

func NewStatusMap(res []*model.Resource, status status) *StatusMap {
	sm := make(StatusMap, len(res))
	for _, r := range res {
		sm[r.ID] = status
	}
	return &sm
}
