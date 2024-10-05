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

func NewStatusMap[T model.Syncable](res []T, status status) *StatusMap {
	sm := make(StatusMap, len(res))
	for _, r := range res {
		sm[r.GetID()] = status
	}
	return &sm
}
