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
			shifted := ((int(c-'A') + shift) % 26) + 'A'
			result.WriteRune(rune(shifted))
		} else if c >= 'a' && c <= 'z' {
			shifted := ((int(c-'a') + shift) % 26) + 'a'
			result.WriteRune(rune(shifted))
		}

	}
	return result.String()
}
