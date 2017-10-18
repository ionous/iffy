package chart

import (
	"github.com/ionous/errutil"
)

// epilogueParser reads the optional "expression" following a directive's prelude.
type epilogueParser struct {
	spaces  []rune
	runes   []rune
	err     error
	fini    rune
	canTrim bool
}

func newEpilogueParser(canTrim bool) *epilogueParser {
	return &epilogueParser{canTrim: canTrim}
}

// GetResult returns the expression text and the control rune, or an error.
// The control rune indicates the kind of epilogue:
// a closing bracket, a trimmed bracket, or a filter.
func (p epilogueParser) GetResult() (ret string, fini rune, err error) {
	if e := p.err; e != nil {
		err = e
	} else {
		ret, fini = string(p.runes), p.fini
	}
	return
}

// NewRune starts on the first character after a directive's prelude;
// the last unhandled rune in a well-formed directive is the terminating rune.
// Note, for trim the terminating rune is not the same as the resulting control rune.
func (p *epilogueParser) NewRune(r rune) (ret State) {
	switch {
	// skip quoted text ( so that brackets and trim in quotes dont end the epilogue )
	case isQuote(r):
		var quote quoteParser
		ret = parseChain(r, &quote, Statement(func(r rune) (ret State) {
			if q, e := quote.GetString(); e != nil {
				p.err = e
			} else {
				p.addRunes([]rune(q)...)
				ret = p.NewRune(r) // quote returns the char after the quote.
			}
			return
		}))

	case isCloseBracket(r):
		p.fini = r
		ret = terminal // done, eat the bracket.

	case isTrim(r):
		trim := r
		ret = Statement(func(r rune) (ret State) {
			if !p.canTrim {
				p.err = errutil.New("unexpected trim")
			} else if !isCloseBracket(r) {
				p.err = errutil.Fmt("unknown character following right trim %q", r)
			} else {
				p.fini = trim
				ret = terminal // done, eat the bracket.
			}
			return
		})

	case isSpace(r):
		p.spaces = append(p.spaces, r)
		ret = p // loop...

	case isFilter(r):
		filter := r
		ret = Statement(func(r rune) (ret State) {
			if r == filter {
				p.addRunes(filter, filter)
				ret = p // loop...
			} else {
				p.fini = filter
				ret = terminal // done, eat the filter
			}
			return
		})

	case r == eof:
		p.err = errutil.New("unclosed directive")

	default:
		p.addRunes(r)
		ret = p // loop....
	}
	return
}

// add the passed runes to the expression text, flushing any accumulated whitespace if needed.
func (p *epilogueParser) addRunes(runes ...rune) {
	if len(p.spaces) > 0 {
		p.runes = append(p.runes, p.spaces...)
		p.spaces = nil
	}
	p.runes = append(p.runes, runes...)
}
