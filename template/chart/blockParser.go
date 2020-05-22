package chart

import (
	"github.com/ionous/iffy/template/postfix"
	"github.com/ionous/iffy/template/types"
)

// BlockParser reads alternating text and directives.
// It is used for tests of Left and RightParser.
type BlockParser struct {
	out     []Directive
	err     error
	factory ExpressionStateFactory // for testing.
}

func (p *BlockParser) StateName() string {
	return "block parser"
}

// GetDirectives or error
func (p *BlockParser) GetDirectives() ([]Directive, error) {
	return p.out, p.err
}

// NewRune starts with the first character of a string.
func (p *BlockParser) NewRune(r rune) State {
	var left LeftParser
	return ParseChain(r, &left, Statement("after lhs", func(r rune) State {
		if text := left.GetText(); len(text) > 0 {
			d := Directive{Expression: quote(text)}
			p.append(d)
		}
		return ParseChain(r, spaces, Statement("lhs spacing", func(r rune) (ret State) {
			if r != eof {
				right := RightParser{factory: p.factory}
				ret = ParseChain(r, &right, Statement("after rhs", func(r rune) (ret State) {
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
