package encrypt

import "strings"

type Caeser struct {
	shift int
}

func NewCaeser(shift int) *Caeser {
	return &Caeser{
		shift: shift,
	}
}

func (c *Caeser) Encrypt(data string) string {
	return c.applyShift(data, c.shift)
}

func (c *Caeser) Decrypt(data string) string {
	rev_shift := (26 - c.shift%26) % 26
	return c.applyShift(data, rev_shift)
}

func (c *Caeser) applyShift(data string, shift int) string {
	var result strings.Builder

	for _, c := range data {

		if c >= 'A' && c <= 'Z' {
			result.WriteRune('A' + (c - 'A' + rune(shift)%26))
		} else if c >= 'a' && c <= 'z' {
			result.WriteRune('a' + (c - 'a' + rune(shift)%26))
		} else {
			result.WriteRune(c)
		}
	}
	return result.String()
}
