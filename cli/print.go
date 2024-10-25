package cli

import (
	"fmt"
	"math"
	"strings"

	"github.com/CnTeng/rx-todo/client"
	"github.com/CnTeng/rx-todo/model"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func (c *cli) PrintProgress(progress *model.Progress, width int) string {
	if progress == nil {
		return c.icons.progress.none
	}

	var b strings.Builder

	doneWidth := int(math.Round(float64(width) * float64(progress.Done) / float64(progress.Total)))

	b.WriteString(strings.Repeat(c.icons.progress.done, doneWidth))

	if undoneWidth := width - doneWidth; undoneWidth > 0 {
		b.WriteString(c.icons.progress.separator)
		b.WriteString(strings.Repeat(c.icons.progress.undone, undoneWidth))
	}

	b.WriteString(fmt.Sprintf(" %d/%d", progress.Done, progress.Total))

	return b.String()
}

func (c *cli) PrintSlice(rs client.ResourceSlice, title string, status client.ResourceStatus) {
	switch rs := rs.(type) {
	case client.ProjectSlice:
		c.PrintProjects(rs, title, status)
	case client.TaskSlice:
		c.PrintTasks(rs, title, status)
	default:
		fmt.Println("Unknown resource type")
	}
}

func (c *cli) PrintProjects(ts client.ProjectSlice, title string, status client.ResourceStatus) {
	tbl := table.NewWriter()
	tbl.SetTitle(title)
	tbl.AppendHeader(table.Row{"#", "Name", "Description", "Progress"})
	tbl.SetStyle(tableConfig)

	for i, t := range ts {
		row := table.Row{}

		row = append(
			row,
			i+1,
			t.Name,
			t.Description,
			c.PrintProgress(t.Progress, 10),
		)

		tbl.AppendRow(row)
	}

	fmt.Println(tbl.Render())
}

func (c *cli) PrintTasks(ts client.TaskSlice, title string, status client.ResourceStatus) {
	tbl := table.NewWriter()
	tbl.SetTitle(title)
	tbl.AppendHeader(table.Row{"#", "", "Name", "Description", "Progress", "Labels"})
	tbl.SetStyle(tableConfig)

	for i, t := range ts {
		row := table.Row{}

		row = append(row, i+1)

		priorityColor := text.FgWhite

		switch t.Priority {
		case model.PriorityNone:
			priorityColor = text.FgWhite
		case model.PriorityLow:
			priorityColor = text.FgBlue
		case model.PriorityMedium:
			priorityColor = text.FgYellow
		case model.PriorityHigh:
			priorityColor = text.FgRed
		}

		if t.Done {
			row = append(row, priorityColor.Sprint(c.icons.done))
		} else {
			row = append(row, priorityColor.Sprint(c.icons.undone))
		}

		labels := make([]string, len(t.Labels))

		if t.Labels != nil {
			for i, l := range t.Labels {
				labels[i] = color.BgRGB(strToRGB(l.Color)).Sprint(l.Name)
			}
		}

		row = append(
			row,
			t.Name,
			t.Description,
			c.PrintProgress(t.Progress, 10),
			strings.Join(labels, " "),
		)

		tbl.AppendRow(row)
	}

	fmt.Println(tbl.Render())
}

func (c *cli) PrintLabels(ls client.LabelSlice, title string, status client.ResourceStatus) {
	tbl := table.NewWriter()
	tbl.SetTitle(title)
	tbl.AppendHeader(table.Row{"#", "Label", "Color"})
	tbl.SetStyle(tableConfig)

	for _, l := range ls {
		row := table.Row{}

		row = append(
			row,
			text.FgYellow.Sprint(l.ID),
			l.Name,
			color.BgRGB(strToRGB(l.Color)).Sprint(l.Color),
		)

		tbl.AppendRow(row)
	}

	fmt.Println(tbl.Render())
}
