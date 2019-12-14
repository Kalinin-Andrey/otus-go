package hw02

import (
	"fmt"
	"strings"
)

func SimpleUnpackString( code string ) (string, error) {
	var res string
	codeRunes := []rune(code)

	for i := 0; i < len(codeRunes); i++ {
		var s string

		if ('a' <= codeRunes[i] && codeRunes[i] <= 'z') || ('A' <= codeRunes[i] && codeRunes[i] <= 'Z') {
			s = string(codeRunes[i])
		} else if codeRunes[i] == '\\' {
			s = string(codeRunes[i+1])
			i++
		} else {
			return "", fmt.Errorf("Symbol N%v must be a valid letter, \"%v\" was givven.", i, codeRunes[i])
		}
		n := 1
		if i + 1 < len(codeRunes) {
			if '0' <= codeRunes[i + 1] && codeRunes[i + 1] <= '9' {
				n = int(codeRunes[i + 1] - '0')
				i++
			}
		}
		res += strings.Repeat(s, n)
	}
	return res, nil
}