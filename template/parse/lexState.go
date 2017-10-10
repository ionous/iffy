package parse

import (
	"github.com/ionous/iffy/template/item"
)

// state implementations scan input, emitting output and yielding new states.
// most but not all states generate some sort of output.
type state interface {
	Lex(l *Lexer) state
}

type lexEnd struct{}
type lexDirective struct{}
type lexPrelude struct{}
type lexTrailingSpace struct{}
type lexReference struct{ dot bool }
type lexTrailingFilter struct{}
type lexFunctionHeader struct{}
type lexParameter struct{}
type lexParameterBody struct{}
type lexExpression struct{}
type lexFilter struct{}
type lexLeftBracket struct{}
type lexRightBracket struct{}
type lexText struct{}
type lexComment struct{}

func (lexEnd) Lex(l *Lexer) (ret state) {
	if l.nestingDepth > 0 {
		ret = l.emitError("unclosed directive")
	} else {
		l.emit(item.End)
		ret = nil
	}
	return
}

// scan the elements inside directive brackets ( aka brackets ).
// We expect to see a keyword, reference, function, expression, or end-bracket.
func (lexDirective) Lex(l *Lexer) (ret state) {
	if l.seedIdentifier() > 0 {
		ret = lexPrelude{}
	} else {
		ret = lexExpression{}
	}
	return
}

// scan a directive which starts with an identifier.
// we expect to see a keyword, reference, function, expression, or end-bracket.
// an identifier, here, means the same as go:
// a sequence of one or more letters and digits; the first of which must be a letter.
// https://golang.org/ref/spec#Identifiers
func (lexPrelude) Lex(l *Lexer) (ret state) {
	// an identifier all on its own can only be a keyword or reference.
	// technically, keywords never have numbers, but this is a helpful simplification.
	if l.rightBracketWidth() > 0 || isSpace(l.PeekRune()) {
		if !l.keywords[l.ItemText()] {
			ret = lexReference{}
		} else {
			l.emit(item.Keyword)
			ret = lexTrailingSpace{}
		}
	} else {
		switch r, w := l.NextRune(); {
		case isLetter(r) || isDigit(r): // letters and digits tell us nothing new
			ret = lexPrelude{}

		case isDot(r) || isEndOfInput(r) || isEndOfLine(r): // a dot means we are a reference
			l.PrevRune(w)
			ret = lexReference{}

		case isSeparator(r): // a function:
			l.currPos -= w // don't send the separator as part of the function name
			l.emit(item.Function)
			l.ignore(w) // skip over the separator
			ret = lexParameter{}

		default: // any other random symbol
			ret = lexExpression{}
		}
	}
	return
}

// after a keyword, we expect to see nothing but net.
func (lexTrailingSpace) Lex(l *Lexer) (ret state) {
	if l.rightBracketWidth() > 0 {
		ret = lexRightBracket{}
	} else {
		switch r, _ := l.NextRune(); {
		case isEndOfInput(r):
			ret = lexEnd{}
		case isEndOfLine(r):
			ret = l.emitError("unexpected line break after keyword")
		case isSpace(r):
			ret = lexTrailingSpace{}
		default:
			ret = l.emitError("expected only whitespace", r)
		}
	}
	return
}

// our current text is part of an identifier.
// we may add more letters, digits, and dots to it, or we may emit the reference and move on.
// filters can follow references.
// FIX? text/template allows new lines after identifiers
func (x lexReference) Lex(l *Lexer) (ret state) {
	r, w := l.NextRune()
	if isLetter(r) || isDigit(r) {
		ret = lexReference{dot: false}
	} else if isDot(r) {
		if !x.dot {
			ret = lexReference{dot: true}
		} else {
			ret = l.emitError("invalid double dot")
		}
	} else {
		l.PrevRune(w)
		l.emit(item.Reference)
		ret = lexTrailingFilter{}
	}
	return
}

// filters may follow references, functions, and expressions.
func (lexTrailingFilter) Lex(l *Lexer) (ret state) {
	if l.rightBracketWidth() > 0 {
		ret = lexRightBracket{}
	} else {
		switch r, w := l.NextRune(); {
		case isEndOfInput(r):
			ret = lexEnd{}
		case isEndOfLine(r):
			ret = l.emitError("unexpected line break")
		case isFilter(r):
			l.PrevRune(w)
			ret = lexFilter{}
		case isSpace(r):
			ret = lexTrailingFilter{}
		default:
			ret = l.emitError("unknown character", r)
		}
	}
	return
}

// scan a filter already seeded with an identifier,
// expecting to eventually see a function.
func (lexFunctionHeader) Lex(l *Lexer) (ret state) {
	r, w := l.NextRune()
	if isLetter(r) || isDigit(r) {
		// letters and digits tell us nothing new
		ret = lexFunctionHeader{}
	} else if isSeparator(r) {
		l.currPos -= w        // don't send the separator as part of the function name
		l.emit(item.Function) //
		l.ignore(w)           // skip over the separator
		ret = lexParameter{}  // on to the body
	} else {
		ret = l.emitError("invalid function name")
	}
	return
}

// we expect to see a spaces leading a parameter,
// followed by a filter, or ending bracket.
func (lexParameter) Lex(l *Lexer) (ret state) {
	if l.seedIdentifier() > 0 {
		ret = lexParameterBody{}
	} else {
		ret = lexTrailingFilter{}
	}
	return
}

// at the start of a new parameter, we might see:
// a right bracket: the end of the current function;
// a left bracket: a sub-directive;
// a reference: "object.something";
// or a filter;
// expressions are disallowed.
func (lexParameterBody) Lex(l *Lexer) (ret state) {
	if r, w := l.NextRune(); isLetter(r) || isDigit(r) {
		// letters and digits tell us nothing new
		ret = lexParameterBody{}
	} else {
		// everything else terminates the parameter:
		l.PrevRune(w)
		l.emit(item.Reference)

		switch {
		case leftBracket.IsPrefix(l.InputText()):
			ret = lexLeftBracket{}
		case isFilter(r):
			ret = lexFilter{}
		case isSpace(r):
			ret = lexParameter{}
		default: // example: right bracket, right trim, filter, bad expression
			ret = lexTrailingFilter{}
		}
	}
	return
}

// starting from somewhere inside an expression directive, search for the end of the expression, and emit it.
// moves to either lexFilter or lexRightBracket
func (lexExpression) Lex(l *Lexer) (ret state) {
	if l.rightBracketWidth() > 0 {
		l.emitTrimmed(item.Expression)
		ret = lexRightBracket{}
	} else {
		switch r, w := l.NextRune(); {
		case isEndOfInput(r):
			l.emitTrimmed(item.Expression)
			ret = lexEnd{}

		case isEndOfLine(r):
			l.emitTrimmed(item.Expression)
			ret = l.emitError("unexpected line break in expression")

		// check for filter, but watch for logical or "||"
		case isFilter(r) && !isFilter(l.PeekRune()):
			l.PrevRune(w)
			l.emitTrimmed(item.Expression)
			ret = lexFilter{}

		default:
			// we accept any old character in an expression
			ret = lexExpression{}
		}
	}
	return
}

// at a filter bracket, eats it, emits it;
// expects to see a function header next.
func (lexFilter) Lex(l *Lexer) (ret state) {
	l.NextRune() // filter is always one character
	l.emit(item.Filter)
	if l.seedIdentifier() > 0 {
		ret = lexFunctionHeader{}
	} else {
		ret = l.emitError("after a filter, expected a function")
	}
	return
}

// at the left bracket, emit it and any trim marker;
// ex. "{~"
// expects to see a comment or directive next
// any preceding trimmed spaces will have been eaten in lexText; we just need to eat the marker.
func (lexLeftBracket) Lex(l *Lexer) (ret state) {
	// move forward across the bracket
	l.currPos += leftBracket.Width()
	// determine if there's a trim marker to eat.
	var afterMarker item.Pos
	if leftTrim.IsPrefix(l.InputText()) {
		afterMarker = leftTrim.Width()
	}
	if l.nestingDepth > 0 && afterMarker > 0 {
		ret = l.emitError("meaningless left trim in sub-directive")
	} else {
		// test for a comment; move to a comment.
		if leftComment.IsPrefix(l.InputTextAt(afterMarker)) {
			// comments dont emit
			l.ignore(afterMarker)
			ret = lexComment{}
		} else {
			l.emit(item.LeftBracket)
			l.ignore(afterMarker)
			l.nestingDepth++
			ret = lexDirective{}
		}
	}
	return
}

// at a right bracket (or trim), emit it, and possibly trim leading spaces;
// // ex. "~}  "
// if this was the outer most bracketed directive, expects a text block to follow;
// if this was a sub-directive, some parameters may follow.
func (lexRightBracket) Lex(l *Lexer) (ret state) {
	trimAfter := rightTrim.IsPrefix(l.InputText())
	if l.nestingDepth > 1 && trimAfter {
		ret = l.emitError("meaningless right trim in sub-directive")
	} else {
		if trimAfter {
			l.ignore(rightTrim.Width())
		}
		l.currPos += rightBracket.Width()
		l.emit(item.RightBracket)
		if trimAfter {
			spaces := leftTrimLength(l.InputText())
			l.ignore(spaces)
		}
		//
		l.nestingDepth--
		if l.nestingDepth == 0 {
			ret = lexText{}
		} else {
			ret = lexParameter{}
		}
	}
	return
}

// at the start of a text block; emit text as a single item.
// expects that a directive or the end of the input follows.
func (lexText) Lex(l *Lexer) (ret state) {
	// scan forward for a left bracket
	if w := leftBracket.Scan(l.InputText()); !w.Valid() {
		// none found: jump to the very end
		l.currPos = item.Width(l.input)
		// and if we have text: emit it
		if text := l.ItemText(); len(text) > 0 {
			l.emit(item.Text)
		}
		ret = lexEnd{}
	} else {
		l.currPos += w
		// evaluate the bracket: if it includes a whitespace trim request
		// determine the length of that whitespace.
		// ex. for "...~{" trimLength will equal 3.
		var trimLength item.Pos
		if leftTrim.IsPrefix(l.InputTextAt(leftBracket.Width())) {
			trimLength = rightTrimLength(l.ItemText())
		}
		// emit text ( by pulling back our current position to before the whitespace )
		l.currPos -= trimLength
		if text := l.ItemText(); len(text) > 0 {
			l.emit(item.Text)
		}
		// ignore the whitespace.
		l.ignore(trimLength)
		ret = lexLeftBracket{}
	}
	return
}

// scan a comment from the start of its marker.
// ex. "{/* comment */}"
// to end a comment means to leave its directive,
// expects that a text block follows.
// errors if the comment doesn't end.
func (lexComment) Lex(l *Lexer) (ret state) {
	l.currPos += leftComment.Width()
	if end := rightComment.Scan(l.InputText()); !end.Valid() {
		ret = l.emitError("unclosed comment")
	} else {
		// skip past the comment body and ending comment marker.
		l.currPos += end + rightComment.Width()
		// we expect to be at a right bracket
		if w := l.rightBracketWidth(); w == 0 {
			ret = l.emitError("comment ends without closing bracket")
		} else {
			l.currPos += w
			if w > rightBracket.Width() {
				l.currPos += leftTrimLength(l.InputText())
			}
			l.startPos = l.currPos
			ret = lexText{}
		}
	}
	return
}
