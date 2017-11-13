package chart

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/template/postfix"
	"github.com/ionous/iffy/template/types"
)

// BlockParser reads alternating text and directives.
type BlockParser struct {
	out     []Directive
	err     error
	factory ExpressionStateFactory // for testing.
}

// LeftParser reads text, ending just after the opening of a directive.
// It deals with leading trim: "...{~".
type LeftParser struct {
	text, spaces []rune
	out          string
}

// RightParser reads the content of a directive, ending just after its end.
// It deals with trailing trim: : "~}...".
type RightParser struct {
	out     Directive
	err     error
	factory ExpressionStateFactory // for testing.
}

// GetDirectives or error
func (p *BlockParser) GetDirectives() ([]Directive, error) {
	return p.out, p.err
}
func (p *LeftParser) GetText() string {
	return p.out
}
func (p *RightParser) GetDirective() (Directive, error) {
	return p.out, p.err
}

// NewRune starts with the first character of a string.
func (p *BlockParser) NewRune(r rune) State {
	var left LeftParser
	return ParseChain(r, &left, Statement(func(r rune) State {
		if text := left.GetText(); len(text) > 0 {
			d := Directive{Expression: quote(text)}
			p.append(d)
		}
		return ParseChain(r, spaces, Statement(func(r rune) (ret State) {
			if r != eof {
				right := RightParser{factory: p.factory}
				ret = ParseChain(r, &right, Statement(func(r rune) (ret State) {
					if v, e := right.GetDirective(); e != nil {
						p.err = e
					} else {
						p.append(v)
						ret = p.NewRune(r) // loop, back to left half.
					}
					return
				}))
			}
			return
		}))
	}))
}

func (p *BlockParser) append(d Directive) {
	p.out = append(p.out, d)
}

func quote(t string) (ret postfix.Expression) {
	if len(t) > 0 {
		ret = []postfix.Function{types.Quote(t)}
	}
	return
}

// NewRune starts with the first character of a string, ends just after the opening bracket of a directive.
func (p *LeftParser) NewRune(r rune) (ret State) {
	switch {
	case isOpenBracket(r):
		ret = Statement(func(r rune) (ret State) {
			if !isTrim(r) {
				p.acceptSpaces()
			} else {
				ret = Terminal //end after eating this trim char
			}
			p.out = string(p.text)
			return
		})
	case r == eof:
		p.acceptSpaces()
		p.out = string(p.text)
	case isSpace(r):
		p.spaces = append(p.spaces, r)
		ret = p // loop...
	default:
		p.acceptSpaces()
		p.text = append(p.text, r)
		ret = p // loop...
	}
	return
}

func (p *LeftParser) acceptSpaces() {
	p.text, p.spaces = append(p.text, p.spaces...), nil
}

// NewRune starts on the first rune of directive content, after the opening bracket, trim, and leading whitespace.
func (p *RightParser) NewRune(r rune) State {
	dir := DirectiveParser{factory: p.factory}
	return ParseChain(r, &dir, Statement(func(r rune) (ret State) {
		if v, e := dir.GetDirective(); e != nil {
			p.err = e
		} else {
			switch {
			case isCloseBracket(r):
				p.out, ret = v, Terminal // eat the trim bracket

			case isTrim(r):
				ret = Statement(func(r rune) (ret State) {
					if !isCloseBracket(r) {
						p.err = errutil.Fmt("unknown character following right trim %q", r)
					} else {
						p.out = v
						ret = spaces // done, eat the subsequent spaces.
					}
					return
				})
			default:
				p.err = errutil.Fmt("unclosed directive %q", r)
			}
			return
		}
		return
	}))
}
