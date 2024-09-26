package cli

import "github.com/fatih/color"

func PrintErr(e error) {
	color.New(color.FgRed, color.Bold).Println("Error:", e)
}
