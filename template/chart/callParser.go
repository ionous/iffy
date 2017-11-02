package chart

import (
	"github.com/ionous/iffy/template/postfix"
)

// CallParser reads a single function call and its arguments.
type CallParser struct {
	argFactory ExpressionStateFactory
	out        postfix.Expression
	err        error
}

func MakeCallParser(f ExpressionStateFactory) CallParser {
	return CallParser{argFactory: f}
}

func (p CallParser) GetExpression() (postfix.Expression, error) {
	return p.out, p.err
}

// NewRune starts with the first character past the bar
func (p *CallParser) NewRune(r rune) State {
	var id IdentParser
	return parseChain(r, spaces, makeChain(&id, Statement(func(r rune) (ret State) {
		// read an identifier, which ends with any unknown character.
		if n := id.GetName(); len(n) > 0 && isSeparator(r) {
			args := MakeArgParser(p.argFactory)
			// use makeChain to skip the separator itself
			ret = makeChain(spaces, makeChain(&args, stateExit(func() {
				if args, arity, e := args.GetArgs(); e != nil {
					p.err = e
				} else {
					cmd := Command{n, arity}
					p.out = append(p.out, args...)
					p.out = append(p.out, cmd)
				}
			})))
		}
		return
	})))
}
