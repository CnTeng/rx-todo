package cli

import "github.com/fatih/color"

type icons struct {
	done   string
	undone string

	none   string
	add    string
	change string
	delete string
}

var nerdIcons = icons{
	done:   " ",
	undone: " ",

	none:   " ",
	add:    "✚",
	change: "",
	delete: "✖",
}

var textIcons = icons{
	done:   "[x]",
	undone: "[ ]",

	none:   " ",
	add:    "+",
	change: "~",
	delete: "-",
}

type iconType int

const (
	Nerd iconType = iota
	Text
)

func newIcons(t iconType) *icons {
	var icons icons

	if t == Nerd {
		icons = nerdIcons
	} else {
		icons = textIcons
	}

	return &icons
}

func (i *icons) withColor() *icons {
	i.done = color.New(color.FgGreen).Sprint(i.done)
	i.add = color.New(color.FgGreen).Sprint(i.add)
	i.change = color.New(color.FgYellow).Sprint(i.change)
	i.delete = color.New(color.FgRed).Sprint(i.delete)

	return i
}
