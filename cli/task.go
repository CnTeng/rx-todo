package cli

import (
	"cmp"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/CnTeng/rx-todo/model"
	"github.com/acarl005/stripansi"
	"github.com/fatih/color"
	"github.com/rodaine/table"
)

type TaskSlice []*model.Task

func (ts *TaskSlice) SortByID() *TaskSlice {
	slices.SortStableFunc(*ts, func(a, b *model.Task) int {
		return cmp.Compare(a.ID, b.ID)
	})
	return ts
}

func (ts *TaskSlice) SortByName() *TaskSlice {
	slices.SortStableFunc(*ts, func(a, b *model.Task) int {
		return cmp.Compare(a.Content, b.Content)
	})
	return ts
}

func (ts *TaskSlice) Filter(f func(*model.Task) bool) *TaskSlice {
	filtered := make(TaskSlice, 0)
	for _, t := range *ts {
		if f(t) {
			filtered = append(filtered, t)
		}
	}
	return &filtered
}

func (ts *TaskSlice) FilterByName(name string) *TaskSlice {
	return ts.Filter(func(t *model.Task) bool {
		return t.Content == name
	})
}

func (c *cli) PrintTasks(ts *TaskSlice) {
	headers := make([]any, 0, 7)

	headers = append(headers, " ", "ID", "Content", "Description", "Priority", "Labels")

	tbl := table.New(headers...).
		WithHeaderFormatter(color.New(color.FgGreen, color.Underline).SprintfFunc()).
		WithWidthFunc(func(s string) int {
			return utf8.RuneCountInString(stripansi.Strip(s))
		})

	for _, t := range *ts {
		vals := make([]any, 0, 7)

		if t.Done {
			vals = append(vals, c.icons.done)
		} else {
			vals = append(vals, c.icons.undone)
		}

		vals = append(
			vals,
			color.New(color.FgYellow).Sprint(t.ID),
			t.Content,
			t.Description,
			t.Priority,
			strings.Join(t.Labels, ", "),
		)

		tbl.AddRow(vals...)
	}

	tbl.Print()
}
