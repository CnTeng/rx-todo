package cli

import (
	"strconv"

	"github.com/fatih/color"
)

func RGB(r, g, b uint8) *color.Color {
	return color.New(
		color.FgWhite+1,
		2,
		color.Attribute(r),
		color.Attribute(g),
		color.Attribute(b),
	)
}

func BgRGB(r, g, b uint8) *color.Color {
	return color.New(
		color.BgWhite+1,
		2,
		color.Attribute(r),
		color.Attribute(g),
		color.Attribute(b),
	)
}

func HexRGB(hex string) *color.Color {
	if hex[0] == '#' {
		hex = hex[1:]
	}

	if len(hex) != 6 {
		return color.New()
	}

	value, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		return color.New()
	}

	r := uint8(value >> 16)
	g := uint8(value >> 8 & 0xFF)
	b := uint8(value & 0xFF)

	return RGB(r, g, b)
}

func BgHexRGB(hex string) *color.Color {
	if hex[0] == '#' {
		hex = hex[1:]
	}

	if len(hex) != 6 {
		return color.New()
	}

	value, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		return color.New()
	}

	r := uint8(value >> 16)
	g := uint8(value >> 8 & 0xFF)
	b := uint8(value & 0xFF)

	return BgRGB(r, g, b)
}
