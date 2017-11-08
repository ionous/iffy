package chart

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/template"
	"github.com/ionous/iffy/template/postfix"
)

// BlockParser reads alternating text and directives.
type BlockParser struct {
	out     []template.Directive
	err     error
	text    []rune
	spaces  []rune
	factory ExpressionStateFactory
}

type ExpressionStateFactory interface {
	NewExpressionState() ExpressionState
}
type ExpressionState interface {
	State
	GetExpression() (postfix.Expression, error)
}

// MakeBlockParser returns a new parser that generates directives via the passed factory.
func MakeBlockParser(f ExpressionStateFactory) (ret BlockParser) {
	ret.factory = f
	return
}

// GetDirectives or error
func (p *BlockParser) GetDirectives() (ret []template.Directive, err error) {
	if e := p.err; e != nil {
		err = e
	} else {
		p.flushText(false)
		ret = p.out
	}
	return
}

// NewRune starts with the first character of a string.
func (p *BlockParser) NewRune(r rune) (ret State) {
	switch {
	case isOpenBracket(r):
		ret = Statement(func(r rune) (ret State) {
			trim := isTrim(r) // write any pending text
			p.flushText(trim)
			if trim {
				// eat the trim character, and any content space
				ret = MakeChain(spaces, Statement(p.afterOpen))
			} else {
				// not trim, pass non-space content along
				ret = ParseChain(r, spaces, Statement(p.afterOpen))
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

// rune at the start of a directive's content.
func (p *BlockParser) afterOpen(r rune) State {
	keyp := KeyParser{exp: p.newExpressionParser()}
	expp := p.newExpressionParser()
	//
	var dir template.Directive
	var err error
	para := MakeParallel(
		MakeChain(expp, StateExit(func() {
			if x, e := expp.GetExpression(); e != nil {
				err = errutil.Append(err, e)
			} else if len(x) > 0 {
				dir = template.Directive{Expression: x} // last match wins
			}
		})),
		MakeChain(&keyp, StateExit(func() {
			if d, e := keyp.GetDirective(); e != nil {
				err = errutil.Append(err, e)
			} else {
				dir = d // last match wins; key wins equal matches
			}
		})),
	)
	return ParseChain(r, para, Statement(func(r rune) (ret State) {
		if err != nil {
			p.err = err
		} else {
			p.out = append(p.out, dir)
			ret = p.afterContent(r)
		}
		return
	}))
}

func (p *BlockParser) newExpressionParser() (ret ExpressionState) {
	if p.factory != nil {
		ret = p.factory.NewExpressionState()
	} else {
		ret = new(PipeParser)
	}
	return
}

// rune after the content of a directive: spaces, trim, closing bracket, etc.
func (p *BlockParser) afterContent(r rune) State {
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

// write any queued text as a block
// if trim is true, we skip trailing spaces, otherwise we write those too.
func (p *BlockParser) flushText(trim bool) {
	text, spaces := p.text, p.spaces
	p.text, p.spaces = nil, nil
	if !trim {
		text = append(text, spaces...)
	}
	if len(text) > 0 {
		q := template.Quote(text)
		p.out = append(p.out, template.Directive{Expression: []postfix.Function{q}})
	}
}
