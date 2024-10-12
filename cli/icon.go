package cli

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

func getIcons(t iconType) *icons {
	switch t {
	case Nerd:
		return &nerdIcons
	case Text:
		return &textIcons
	default:
		return &textIcons
	}
}
