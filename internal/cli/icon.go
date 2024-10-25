package cli

import "github.com/jedib0t/go-pretty/v6/text"

type progressIcons struct {
	none      string
	done      string
	separator string
	undone    string
}

type icons struct {
	done   string
	undone string

	none   string
	add    string
	change string
	delete string

	progress progressIcons
}

var nerdIcons = icons{
	done:   " ",
	undone: " ",

	none:   " ",
	add:    "+",
	change: "~",
	delete: "-",

	progress: progressIcons{
		none:      "─",
		done:      "━",
		separator: "╺",
		undone:    "━",
	},
}

var textIcons = icons{
	done:   "[x]",
	undone: "[ ]",

	none:   " ",
	add:    "+",
	change: "~",
	delete: "-",

	progress: progressIcons{
		none:      "─",
		done:      "━",
		separator: "╺",
		undone:    "━",
	},
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
	i.done = text.FgGreen.Sprint(i.done)

	i.add = text.FgGreen.Sprint(i.add)
	i.change = text.FgYellow.Sprint(i.change)
	i.delete = text.FgRed.Sprint(i.delete)

	i.progress.done = text.FgGreen.Sprint(i.progress.done)
	i.progress.undone = text.Faint.Sprint(i.progress.undone)
	i.progress.separator = text.Faint.Sprint(i.progress.separator)
	i.progress.none = text.Faint.Sprint(i.progress.none)

	return i
}
