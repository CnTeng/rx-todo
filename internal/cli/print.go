package cli

import (
	"fmt"
	"math"
	"strings"

	"github.com/CnTeng/rx-todo/client"
	"github.com/CnTeng/rx-todo/internal/model"
	"github.com/fatih/color"
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
	tbl := NewTable(50, true)
	tbl.AddHeader("#", "Name", "Description", "Progress")
	tbl.SetHeaderStyle(color.Underline)

	for i, t := range ts {
		tbl.AddRow(
			i+1,
			t.Name,
			t.Description,
			c.PrintProgress(t.Progress, 10),
		)
	}

	fmt.Print(tbl.Render())
}

func (c *cli) PrintTasks(ts client.TaskSlice, title string, status client.ResourceStatus) {
	tbl := NewTable(80, true)
	tbl.AddHeader("ID", "  ", "Labels", "Name", "Progress")
	tbl.SetHeaderStyle(color.Underline)

	for _, t := range ts {
		row := []any{}

		row = append(row, t.ID)

		priorityColor := color.FgWhite

		switch t.Priority {
		case model.PriorityNone:
			priorityColor = color.FgWhite
		case model.PriorityLow:
			priorityColor = color.FgBlue
		case model.PriorityMedium:
			priorityColor = color.FgYellow
		case model.PriorityHigh:
			priorityColor = color.FgRed
		}

		if t.Done {
			row = append(row, color.New(priorityColor).Sprint(c.icons.done))
		} else {
			row = append(row, color.New(priorityColor).Sprint(c.icons.undone))
		}

		labels := make([]string, len(t.Labels))
		if t.Labels != nil {
			for i, l := range t.Labels {
				labels[i] = color.BgRGB(strToRGB(l.Color)).Sprint(l.Name)
			}
		}

		row = append(
			row,
			strings.Join(labels, " "),
			t.Name,
			c.PrintProgress(t.Progress, 10),
		)

		tbl.AddRow(row...)
	}

	fmt.Print(tbl.Render())
}

func (c *cli) PrintLabels(ls client.LabelSlice, title string, status client.ResourceStatus) {
	tbl := NewTable(50, true)
	tbl.AddHeader("#", "Label", "Color")

	for _, l := range ls {
		tbl.AddRow(
			color.YellowString("%d", l.ID),
			l.Name,
			color.BgRGB(strToRGB(l.Color)).Sprint(l.Color),
		)
	}

	fmt.Print(tbl.Render())
}
