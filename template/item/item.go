package item

import (
	"fmt"
	"strconv"
)

// Data represents a token or text string returned from the scanner.
type Data struct {
	Type Type   // type of this data.
	Pos  Pos    // offset, in bytes, of this data in the input string.
	Val  string // value of this data.
	Line int    // line number at the start of this data.
}

// Pos is a byte offset in the original input text.
type Pos int

func (p Pos) Valid() bool    { return p >= 0 }
func (p Pos) String() string { return strconv.Itoa(int(p)) }

func Width(l string) Pos {
	n := len(l)
	return Pos(n)
}

// Type of the lexer item data.
type Type int

//go:generate stringer -type=Type
const (
	Error Type = iota // error occurred; value is text of error
	End               // end of stream

	Expression // pretty much anything not understood
	Function   // identifer with trailing separator; or, leading a filter
	Reference  // a non-keyword, possibly with dots.
	Keyword    // one of the set of keywords provided to the lexer

	LeftBracket  // left directive delimiter
	Filter       // pipe symbol, aka filter.
	RightBracket // right directive delimiter
	Text         // plain text -- ie. all things not between delims
)

// String representation of the item value.
// Type has its own stringer implementation.
func (i Data) String() (ret string) {
	switch {
	case i.Type == End:
		ret = "End"
	case i.Type == Error:
		ret = i.Val
	case len(i.Val) > 10:
		ret = fmt.Sprintf("%.10q...", i.Val)
	default:
		ret = fmt.Sprintf("%q", i.Val)
	}
	return
}
