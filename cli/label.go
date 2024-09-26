package cli

import (
	"strings"
	"unicode/utf8"

	"github.com/CnTeng/rx-todo/internal/model"
	"github.com/acarl005/stripansi"
	"github.com/fatih/color"
	"github.com/rodaine/table"
)

type LabelSlice []*model.Label

func (ls *LabelSlice) List(s *StatusMap) string {
	var builder strings.Builder
	var tbl table.Table
	hasStatus := false

	if s != nil {
		hasStatus = true
		tbl = table.New(" ", "Label", "Color", "Updated At")
	} else {
		tbl = table.New("Label", "Color", "Updated At")
	}

	tbl.WithWriter(&builder).
		WithHeaderFormatter(color.New(color.FgGreen, color.Underline).SprintfFunc()).
		WithWidthFunc(func(s string) int {
			return utf8.RuneCountInString(stripansi.Strip(s))
		})

	for _, c := range *ls {
		if hasStatus {
			status := (*s)[c.ID]
			tbl.AddRow(
				status.String(),
				c.Name,
				BgHexRGB(c.Color).Sprint(c.Color),
				c.UpdatedAt.Format("2006-01-02 15:04:05"),
			)
		} else {
			tbl.AddRow(
				c.Name,
				BgHexRGB(c.Color).Sprint(c.Color),
				c.UpdatedAt.Format("2006-01-02 15:04:05"),
			)
		}
	}

	builder.WriteString("\n")
	tbl.Print()
	builder.WriteString("\n")

	return builder.String()
}
