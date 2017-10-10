package parse

import (
	"fmt"
	"github.com/ionous/iffy/template/item"
	"strings"
)

// Lexer scans a window for tokens.
type Lexer struct {
	Window
	state state
	items []item.Data
}

// ScanDirective starts a new lexer as if it started inside a directive.
func ScanDirective(in Window) *Lexer {
	return &Lexer{in, lexDirective{}, nil}
}

// ScanText starts a new lexer as if it started outside a directive.
func ScanText(in Window) *Lexer {
	return &Lexer{in, lexText{}, nil}
}

// Next returns a single parsed token.
// It returns false if the lexer can't generate any tokens.
func (l *Lexer) Next() (ret item.Data, okay bool) {
	for l.state != nil && len(l.items) == 0 {
		l.state = l.state.Lex(l)
	}
	if cnt := len(l.items); cnt > 0 {
		ret = l.items[cnt-1]
		l.items = l.items[0 : cnt-1]
		okay = true
	}
	return
}

// Drain accumulates tokens until done, or until the passed number state iterations have happened.
// If iter < 0, it runs till completely done.
// Returns the queue.
func (l *Lexer) Drain(iter int) (ret []item.Data) {
	for i := 0; l.state != nil && (iter < 0 || i < iter); i++ {
		l.state = l.state.Lex(l)
	}
	return l.items
}

// ignore skips over the pending input before this point.
func (l *Lexer) ignore(ofs item.Pos) {
	l.currPos += ofs
	l.startPos = l.currPos
}

// helper to skip left-side whitespace and read any initial letters.
// returns the number of letters read ( if any )
func (l *Lexer) seedIdentifier() (ret item.Pos) {
	// trimLength := leftTrimLength(l.InputText())
	text := l.InputText()
	n := strings.TrimLeftFunc(text, isSpace)
	trimLength := item.Width(text) - item.Width(n)
	l.ignore(trimLength)
	for {
		if r, w := l.NextRune(); !isLetter(r) {
			l.PrevRune(w)
			break
		}
		ret++
	}
	return
}

// accept consumes the next rune if it's from the valid set.
// func (l*Lexer) accept(valid string) (okay bool) {
// 	if r, w := l.NextRune(); strings.ContainsRune(valid, r) {
// 		okay = true
// 	} else {
// 		l.PrevRune(w)
// 	}
// 	return
// }

// // acceptRun consumes a run of runes from the valid set.
// func (l*Lexer) acceptRun(valid string) {
// 	for l.accept(valid) {
// 	}
// }

// emit an item of the passed type to the client.
// the item contains the current start, ItemText(), and the total line count.
func (l *Lexer) emit(t item.Type) {
	text := l.ItemText()
	data := item.Data{t, l.startPos, text, l.line}
	l.items = append(l.items, data)
	// Some items contain text internally. If so, count their newlines.
	switch t {
	case /*ItemText, itemRawString,*/
		item.LeftBracket, item.RightBracket:
		l.line += strings.Count(text, "\n")
	}
	l.startPos = l.currPos
}

// helper for lexExpression to remove right-side whitespace.
func (l *Lexer) emitTrimmed(t item.Type) {
	trimLength := rightTrimLength(l.ItemText())
	l.currPos -= trimLength
	if len(l.ItemText()) > 0 {
		l.emit(t)
	}
	l.ignore(trimLength) // skip the whitespace
}

// emitError returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.NextRuneItem.
func (l *Lexer) emitError(format string, args ...interface{}) state {
	msg := fmt.Sprintf(format, args...)
	data := item.Data{item.Error, l.startPos, msg, l.line}
	l.items = append(l.items, data)
	return nil
}

// returns the length of the right side markers, or 0 if not at a right side delimiter.
// for example: "~}" returns 2
func (l *Lexer) rightBracketWidth() (ret item.Pos) {
	if in := l.InputText(); rightBracket.IsPrefix(in) {
		ret = rightBracket.Width()
	} else if rightTrim.IsPrefix(in) &&
		rightBracket.IsPrefix(in[rightTrim.Width():]) {
		ret = rightBracket.Width() + rightTrim.Width()
	}
	return
}
