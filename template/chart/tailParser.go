package chart

import (
	"github.com/ionous/errutil"
)

// tailParser reads the optional "expression" following a directive's header.
type tailParser struct {
	spaces  []rune
	runes   []rune
	err     error
	fini    rune
	canTrim bool
}

// GetTail returns the expression text and the control rune, or an error.
// The control rune indicates *how* the tail was ended:
// an ending bracket, a trim character, or a filter.
func (p tailParser) GetTail() (ret string, fini rune, err error) {
	if e := p.err; e != nil {
		err = e
	} else {
		ret, fini = string(p.runes), p.fini
	}
	return
}

// NewRune is the first character after a directive's header
// The last unhandled rune in a well-formed directive is the terminating rune,
// which -- for trim -- is not the control rune returned from trim char.
func (p *tailParser) NewRune(r rune) (ret State) {
	switch {
	case isEndBracket(r):
		bracket := r
		p.fini = bracket // done.

	case p.canTrim && isTrim(r):
		trim := r
		ret = Statement(func(r rune) State {
			if isEndBracket(r) {
				p.fini = trim // done.
			} else {
				p.err = errutil.Fmt("unknown character following right trim %q", r) // done.
			}
			return nil // we are a state exit action.
		})

	case isSpace(r):
		p.spaces = append(p.spaces, r)
		ret = p // loop...

	case isFilter(r):
		filter := r
		ret = Statement(func(r rune) (ret State) {
			if !isFilter(r) {
				p.fini = filter // done.
			} else {
				p.addRunes(r, r)
				ret = p // loop...
			}
			return ret
		})

	case r == eof:
		p.err = errutil.New("unclosed directive") // done.

	default:
		p.addRunes(r)
		ret = p // loop....
	}
	return
}

// add the passed runes to the expression text, flushing any accumulated whitespace if needed.
func (p *tailParser) addRunes(runes ...rune) {
	if len(p.spaces) > 0 {
		runes = append(runes, p.spaces...)
		p.spaces = nil
	}
	p.runes = append(p.runes, runes...)
}
