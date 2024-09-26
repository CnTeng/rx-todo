package cli

import (
	"fmt"
	"strings"

	"github.com/CnTeng/rx-todo/internal/model"
	"github.com/fatih/color"
)

type TaskSlice []*model.Task

func (ts *TaskSlice) List() string {
	var builder strings.Builder

	donePrint := color.New(color.CrossedOut).SprintFunc()

	for _, t := range *ts {
		if t.Done {
			builder.WriteString(fmt.Sprintf("%s %s\n", color.GreenString(" "), donePrint(t.Content)))
		} else {
			builder.WriteString(fmt.Sprintf("%s %s\n", " ", t.Content))
		}
	}

	return builder.String()
}
