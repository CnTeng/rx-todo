package cli

import "strconv"

func strToRGB(hex string) (int, int, int) {
	if len(hex) != 7 || hex[0] != '#' {
		return 0, 0, 0
	}

	value, err := strconv.ParseUint(hex[1:], 16, 32)
	if err != nil {
		return 0, 0, 0
	}

	r := int(value >> 16)
	g := int(value >> 8 & 0xFF)
	b := int(value & 0xFF)

	return r, g, b
}
