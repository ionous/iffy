package parse

import (
	"github.com/ionous/iffy/template/item"
	"strings"
	"unicode"
)

const (
	spaceChars = " \t\r\n" // space characters as defined by golang: see also isSpace; isEndOfLine.
	// note: unicode whitespace includes more: '\v', '\f', U+0085 (NEL), U+00A0 (NBSP).
	eof        = -1  // stand in for a rune indicating the end of input.
	filterChar = '|' // one rune
	dotChar    = '.'
)

func isFilter(r rune) bool {
	return r == filterChar
}

func isDot(r rune) bool {
	return r == dotChar
}

// isSeparator reports whether r indicates a function call.
func isSeparator(r rune) bool {
	return r == ':' || r == '?' || r == '!'
}

// isSpace reports whether r is whitespace but not newline.
func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

// isBreak returns whether the input/line ends unexpectedly
func isEndOfInput(r rune) bool {
	return r == eof
}

// isEndOfLine reports whether r is an end-of-line character.
func isEndOfLine(r rune) bool {
	return r == '\r' || r == '\n'
}

// isLetter reports whether r is a golang letter
// https://golang.org/ref/spec#letter
func isLetter(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}

func isDigit(r rune) bool {
	return unicode.IsDigit(r)
}

// rightTrimLength returns the length of the spaces at the end of the string.
func rightTrimLength(s string) (ret item.Pos) {
	x := strings.TrimRight(s, spaceChars)
	return item.Width(s) - item.Width(x)
}

// leftTrimLength returns the length of the spaces at the beginning of the string.
func leftTrimLength(s string) (ret item.Pos) {
	x := strings.TrimLeft(s, spaceChars)
	return item.Width(s) - item.Width(x)
}
