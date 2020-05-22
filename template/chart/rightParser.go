package chart

import "github.com/ionous/errutil"

// RightParser reads the content of a directive, ending just after the closing brace and any trimed whitespace. "~}...".
type RightParser struct {
	out     Directive
	err     error
	factory ExpressionStateFactory // for testing.
}

func (p *RightParser) StateName() string {
	return "rhs"
}
func (p *RightParser) GetDirective() (Directive, error) {
	return p.out, p.err
}

// NewRune starts on the first rune of directive content, after the opening bracket, trim, and leading whitespace.
func (p *RightParser) NewRune(r rune) State {
	dir := DirectiveParser{factory: p.factory}
	return ParseChain(r, &dir, Statement("after rhs directive", func(r rune) (ret State) {
		if v, e := dir.GetDirective(); e != nil {
			p.err = e
		} else {
			switch {
			case isCloseBracket(r):
				p.out, ret = v, Terminal // eat the trim bracket

			case isTrim(r):
				ret = Statement("closing", func(r rune) (ret State) {
					if !isCloseBracket(r) {
						p.err = errutil.Fmt("unknown character following right trim %q", r)
					} else {
						p.out = v
						ret = spaces // done, eat the subsequent spaces.
					}
					return
				})

			default:
				p.err = errutil.Fmt("unexpected end of directive after '%s:%s'", v.Key, v.Expression.String())
			}
			return
		}
		return
	}))
}
