package cli

import (
	"sort"
	"unicode/utf8"

	"github.com/CnTeng/rx-todo/model"
	"github.com/acarl005/stripansi"
	"github.com/fatih/color"
	"github.com/rodaine/table"
)

type LabelSlice []*model.Label

func (ls *LabelSlice) SortByID() *LabelSlice {
	sort.Slice(*ls, func(i, j int) bool {
		return (*ls)[i].ID < (*ls)[j].ID
	})
	return ls
}

func (ls *LabelSlice) SortByName() *LabelSlice {
	sort.Slice(*ls, func(i, j int) bool {
		return (*ls)[i].Name < (*ls)[j].Name
	})
	return ls
}

func (c *cli) PrintLabels(ls *LabelSlice, sm *statusMap) {
	headers := make([]any, 0, 4)

	if sm != nil {
		headers = append(headers, " ")
	}
	headers = append(headers, "ID", "Label", "Color")

	tbl := table.New(headers...).
		WithHeaderFormatter(color.New(color.FgGreen, color.Underline).SprintfFunc()).
		WithWidthFunc(func(s string) int {
			return utf8.RuneCountInString(stripansi.Strip(s))
		})

	for _, l := range *ls {
		vals := make([]any, 0, 4)
		if sm != nil {
			vals = append(vals, sm.getStatusIcon(l.ID, c.icons))
		}

		vals = append(
			vals,
			color.New(color.FgYellow).Sprint(l.ID),
			l.Name,
			BgHexRGB(l.Color).Sprint(l.Color),
		)

		tbl.AddRow(vals...)
	}

	tbl.Print()
}
