package hw02

import (
	"fmt"
	"strings"
	"unicode"
)

func SimpleUnpackString( code string ) (string, error) {
	var b strings.Builder
	codeRunes := []rune(code)

	for i := 0; i < len(codeRunes); i++ {
		var s string

		switch {
		case unicode.IsLetter(codeRunes[i]):
			s = string(codeRunes[i])
		case codeRunes[i] == '\\':
			s = string(codeRunes[i+1])
			i++
		default:
			return "", fmt.Errorf("Symbol N%v must be a valid letter, \"%v\" was givven.", i, codeRunes[i])
		}
		n := 1

		if i + 1 < len(codeRunes) {
			if unicode.IsDigit(codeRunes[i + 1]) {
				n = int(codeRunes[i + 1] - '0')
				i++
			}
		}
		b.WriteString(strings.Repeat(s, n))
	}
	return b.String(), nil
}