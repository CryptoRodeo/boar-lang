package helper

import (
	"bytes"
	"strings"

	"github.com/TwiN/go-color"
)

func ApplyColorToText(str string) string {
	var out bytes.Buffer
	text := strings.Split(str, "")
	for _, val := range text {
		out.WriteString(decideColor(val[0]))
	}
	return out.String()
}

func decideColor(token byte) string {
	switch {
	case isDigit(token):
		return color.Ize(color.Red, string(token))
	case isDelimiter(token) || token == '"':
		return color.InCyan(string(token))
	case isBrace(token):
		return color.Ize(color.Yellow, string(token))
	default:
		return color.Ize(color.Blue, string(token))
	}
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isDelimiter(ch byte) bool {
	str := string(ch)
	return str == "," || str == ";" || str == ":" || str == "."
}

func isBrace(ch byte) bool {
	str := string(ch)
	return str == "{" || str == "}" || str == "[" || str == "]"
}
