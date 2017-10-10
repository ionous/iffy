package parse

import (
	"github.com/ionous/iffy/template/item"
	"strings"
)

// marker represents any common template symbol.
// text/template allows configuration of the delimiters and it uses trim markers with dashes.
// by keeping the delims constant, and knowing they arent letters, digits, whitespace, and dots --
// we are able to simplify the code.
// additionally, since tilde is not a golang symbol it doesnt conflict with unary expressions.
type marker string

func (m marker) IsPrefix(l string) bool {
	return strings.HasPrefix(l, string(m))
}

// Width returns the marker length as a post
func (m marker) Width() item.Pos {
	return item.Width(string(m))
}

// Scan returns the index of this marker in the passed string
func (m marker) Scan(l string) item.Pos {
	i := strings.Index(l, string(m))
	return item.Pos(i)
}

const (
	leftTrim     marker = "~" // If a directive begins "{~" then all preceeding whitespace is removed
	rightTrim    marker = "~" // if it ends "~}" the leading spaces are trimmed.
	leftBracket  marker = "{"
	rightBracket marker = "}"
	leftComment  marker = "/*"
	rightComment marker = "*/"
)
