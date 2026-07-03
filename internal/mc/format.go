package mc

import (
	"regexp"
	"strings"
)

const section = '\u00a7' // §

var ansiCodes = map[rune]string{
	'0': "30", '1': "34", '2': "32", '3': "36",
	'4': "31", '5': "35", '6': "33", '7': "37",
	'8': "90", '9': "94", 'a': "92", 'b': "96",
	'c': "91", 'd': "95", 'e': "93", 'f': "97",
	'l': "1", 'm': "9", 'n': "4", 'o': "3",
	'k': "", 'r': "0",
}

func lower(r rune) rune {
	if r >= 'A' && r <= 'Z' {
		return r + 32
	}
	return r
}

func ToANSI(s string) string {
	var b strings.Builder
	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		if runes[i] == section && i+1 < len(runes) {
			if ansi, ok := ansiCodes[lower(runes[i+1])]; ok {
				if ansi != "" {
					b.WriteString("\x1b[" + ansi + "m")
				}
				i++
				continue
			}
		}
		b.WriteRune(runes[i])
	}
	b.WriteString("\x1b[0m")
	return b.String()
}

func Strip(s string) string {
	var b strings.Builder
	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		if runes[i] == section && i+1 < len(runes) {
			i++
			continue
		}
		b.WriteRune(runes[i])
	}
	return b.String()
}

var ampRe = regexp.MustCompile(`&([0-9a-fk-orA-FK-OR])`)

func Amp(s string) string {
	return ampRe.ReplaceAllString(s, "\u00a7$1")
}
