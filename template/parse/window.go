package parse

import (
	"github.com/ionous/iffy/template/item"
	"unicode/utf8"
)

// Window highlights a portion of an template input.
type Window struct {
	input        string          // the string being scanned
	startPos     item.Pos        // window offset ( measured from input's start  )
	currPos      item.Pos        // offset within window ( from input start )
	nestingDepth int             // nesting depth of sub exprs
	line         int             // number of newlines seen
	keywords     map[string]bool // special identifiers used by iffy: ex. "if"
}

// MakeWindow creates a scanner for the input string.
func MakeWindow(input string, keywords map[string]bool) Window {
	return Window{input: input, keywords: keywords}
}

// ItemText returns the text of the current item up to the current position
func (l *Window) ItemText() string {
	return l.input[l.startPos:l.currPos]
}

// InputText returns the text from the current position to the end of input
func (l *Window) InputText() string {
	return l.input[l.currPos:]
}

// InputTextAt behaves like InputText, only offset by the passed value.
func (l *Window) InputTextAt(ofs item.Pos) string {
	return l.input[l.currPos+ofs:]
}

// returns the next rune in the input, or the constant 'eof'
// it increments the internal line count for newline characters
func (l *Window) NextRune() (ret rune, width item.Pos) {
	if int(l.currPos) >= len(l.input) {
		ret = eof
	} else {
		r, w := utf8.DecodeRuneInString(l.InputText())
		if r == '\n' {
			l.line++
		}
		width = item.Pos(w)
		l.currPos += width
		ret = r
	}
	return
}

// returns, but does not consume, the next rune in the input.
func (l *Window) PeekRune() rune {
	r, w := l.NextRune()
	l.PrevRune(w)
	return r
}

// steps back one rune of width w.
func (l *Window) PrevRune(w item.Pos) {
	l.currPos -= w
	// Correct newline count.
	if w == 1 && l.input[l.currPos] == '\n' {
		l.line--
	}
}
