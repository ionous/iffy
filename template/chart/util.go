package chart

import (
	"unicode"
)

const eof = rune(-1)

// isLetter reports whether r is a golang letter.
// https://golang.org/ref/spec#letter
func isLetter(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}

// https://golang.org/ref/spec#decimal_digit
func isNumber(r rune) bool {
	return r >= '0' && r <= '9'
}

// we use qualifiers within identifier names for custom counters, etc.
func isQualifier(r rune) bool {
	return r == '#' || r == '-'
}

// https://golang.org/ref/spec#hex_digit
func isHex(r rune) bool {
	return (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F') || isNumber(r)
}

// isSpace reports whether r is whitespace but not newline.
func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

// isSeparator reports whether r indicates a function call.
func isSeparator(r rune) bool {
	return r == ':' || r == '?' || r == '!'
}

// isQuote reports whether r is a quote marker
// currently, they all act like https://golang.org/ref/spec#interpreted_string_lit.
func isQuote(r rune) bool {
	return r == '\'' || r == '"'
}

// for consistancy with isCloseBracket
func isOpenBracket(r rune) bool {
	return r == '{'
}

// the rune messes up sublime go switch bracket matching.
func isCloseBracket(r rune) bool {
	return r == '}'
}

// for consistancy with isCloseBracket
func isOpenParen(r rune) bool {
	return r == '('
}

// the rune messes up sublime go paren matching.
func isCloseParen(r rune) bool {
	return r == ')'
}

func isPipe(r rune) bool {
	return r == '|'
}

func isTrim(r rune) bool {
	return r == '~'
}

func isDot(r rune) bool {
	return r == '.'
}
