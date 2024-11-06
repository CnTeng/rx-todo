package cli

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

var tableConfig = table.Style{
	Name: "todoStyle",
	Box:  table.StyleBoxLight,
	Color: table.ColorOptions{
		IndexColumn:  text.Colors{},
		Footer:       text.Colors{},
		Header:       text.Colors{text.FgGreen, text.Underline},
		Row:          text.Colors{},
		RowAlternate: text.Colors{},
	},
	Format:  table.FormatOptions{},
	Options: table.OptionsNoBordersAndSeparators,
}

type cli struct {
	icons *icons
}

func NewCLI(iconType iconType) *cli {
	return &cli{icons: newIcons(iconType).withColor()}
}
