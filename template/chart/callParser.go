package chart

import (
	"github.com/ionous/iffy/template/postfix"
	"github.com/ionous/iffy/template/types"
)

// CallParser reads a single function call and its arguments.
type CallParser struct {
	arity      int
	argFactory ExpressionStateFactory
	out        postfix.Expression
	err        error
}

func MakeCallParser(a int, f ExpressionStateFactory) CallParser {
	return CallParser{arity: a, argFactory: f}
}

func (p *CallParser) StateName() string {
	return "call parser"
}

func (p *CallParser) GetExpression() (ret postfix.Expression, err error) {
	return p.out, p.err
}

// NewRune starts with the first character past the bar.
func (p *CallParser) NewRune(r rune) State {
	var id IdentParser
	return ParseChain(r, spaces, MakeChain(&id, Statement("after call", func(r rune) (ret State) {
		// read a function: an identifier which ends with a separator.
		// We dont fail if we dont end with a separator --
		// the assumption is that it might be an id
		if n := id.Identifier(); len(n) > 0 && isSeparator(r) {
			args := ArgParser{factory: p.argFactory}
			// use MakeChain to skip the separator itself
			ret = MakeChain(spaces, MakeChain(&args, StateExit("call", func() {
				if args, arity, e := args.GetArguments(); e != nil {
					p.err = e
				} else {
					// this follows the spirit of postfix.Pipe without using the actual algorithm.
					var prev postfix.Expression
					prev, p.out = p.out, nil
					if len(args) > 0 {
						p.out = append(p.out, args...)
					}
					if len(prev) > 0 {
						p.out = append(p.out, prev...)
					}
					cmd := types.Command{n, arity + p.arity}
					p.out = append(p.out, cmd)
				}
			})))
		}
		return
	})))
}
