package chart

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/template"
	"github.com/ionous/iffy/template/postfix"
)

// OldParser was the simplified BlockParser before it was split into left and right.
type OldParser struct {
	out     []Directive
	err     error
	text    []rune
	spaces  []rune
	factory ExpressionStateFactory // for testing.
}

// GetDirectives or error
func (p *OldParser) GetDirectives() (ret []Directive, err error) {
	if e := p.err; e != nil {
		err = e
	} else {
		p.flushText(false)
		ret = p.out
	}
	return
}

// NewRune starts with the first character of a string.
func (p *OldParser) NewRune(r rune) (ret State) {
	switch {
	case isOpenBracket(r):
		ret = Statement(func(r rune) (ret State) {
			// *prepare* to read the directive.
			dir := DirectiveParser{factory: p.factory}
			series := MakeChain(spaces, MakeChain(&dir, Statement(func(r rune) (ret State) {
				if v, e := dir.GetDirective(); e != nil {
					p.err = e
				} else {
					p.out = append(p.out, v)
					ret = ParseChain(r, spaces, Statement(p.afterContent))
				}
				return
			})))
			// if the first rune after the open is a trim char:
			// we want to eat that rune, and flush any whitespace weve accumulated.
			// if its not trim we want to pass that first rune along.
			if isTrim(r) {
				p.flushText(true) // eat the trim
				ret = series
			} else {
				p.flushText(false) // not trim, pass all content along
				ret = series.NewRune(r)
			}
			return
		})

	case isSpace(r):
		p.spaces = append(p.spaces, r)
		ret = p // loop...

	case r != eof:
		p.text = append(p.text, p.spaces...)
		p.text = append(p.text, r)
		p.spaces = nil
		ret = p // loop...
	}
	return
}

// rune after the content of a directive: spaces, trim, closing bracket, etc.
// we loop back to the block parser assuming there's no error.
func (p *OldParser) afterContent(r rune) State {
	return ParseChain(r, spaces, Statement(func(r rune) (ret State) {
		switch {
		case isCloseBracket(r):
			ret = p // loop...

		case isTrim(r):
			ret = MakeChain(Statement(func(r rune) (ret State) {
				if !isCloseBracket(r) {
					p.err = errutil.Fmt("unknown character following right trim %q", r)
				} else {
					ret = spaces // done, eat the closing bracket and subsequent spaces.
				}
				return
			}), p) // after trimming, loop ...

		default:
			p.err = errutil.Fmt("unclosed directive %q", r)
		}
		return
	}))
}

// write any queued text as a block;
// if trim is true, we skip trailing spaces, otherwise we write those too.
func (p *OldParser) flushText(trim bool) {
	text, spaces := p.text, p.spaces
	p.text, p.spaces = nil, nil
	if !trim {
		text = append(text, spaces...)
	}
	if len(text) > 0 {
		q := template.Quote(text)
		p.out = append(p.out, Directive{Expression: []postfix.Function{q}})
	}
}
