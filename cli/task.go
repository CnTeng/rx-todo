package cli

import (
	"strings"
	"unicode/utf8"

	"github.com/CnTeng/rx-todo/model"
	"github.com/acarl005/stripansi"
	"github.com/fatih/color"
	"github.com/rodaine/table"
)

type TaskSlice []*model.Task

func (ts *TaskSlice) List(s *StatusMap) string {
	var builder strings.Builder
	var tbl table.Table
	hasStatus := false
	var doneIcon string

	if s != nil {
		hasStatus = true
		tbl = table.New(" ", " ", "ID", "Content", "Description", "Priority", "Labels")
	} else {
		tbl = table.New(" ", "ID", "Content", "Description", "Priority", "Labels")
	}

	tbl.WithWriter(&builder).
		WithHeaderFormatter(color.New(color.FgGreen, color.Underline).SprintfFunc()).
		WithWidthFunc(func(s string) int {
			return utf8.RuneCountInString(stripansi.Strip(s))
		})

	for _, t := range *ts {
		if hasStatus {
			status := (*s)[t.ID]
			if t.Done {
				doneIcon = color.New(color.FgGreen).Sprint(" ")
			} else {
				doneIcon = " "
			}

			tbl.AddRow(
				status.String(),
				doneIcon,
				color.New(color.FgYellow).Sprint(t.ID),
				t.Content,
				t.Description,
				t.Priority,
				strings.Join(t.Labels, ", "),
			)
		} else {
			if t.Done {
				doneIcon = color.New(color.FgGreen).Sprint(" ")
			} else {
				doneIcon = " "
			}

			tbl.AddRow(
				doneIcon,
				color.New(color.FgYellow).Sprint(t.ID),
				t.Content,
				t.Description,
				t.Priority,
				strings.Join(t.Labels, ", "),
			)
		}
	}

	builder.WriteString("\n")
	tbl.Print()
	builder.WriteString("\n")

	return builder.String()
}
