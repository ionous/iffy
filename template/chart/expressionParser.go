package chart

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/template/postfix"
)

// ExpressionParser reads either a call or a sequence.
type ExpressionParser struct {
	argFactory ExpressionStateFactory
	out        postfix.Expression
	err        error
}

func MakeExpressionParser(f ExpressionStateFactory) ExpressionParser {
	return ExpressionParser{argFactory: f}
}

func (p ExpressionParser) GetExpression() (postfix.Expression, error) {
	return p.out, p.err
}

// NewRune starts with the first character past the bar
func (p *ExpressionParser) NewRune(r rune) State {
	call := MakeCallParser(0, p.argFactory)
	seq := SequenceParser{}
	par := MakeParallel(
		MakeChain(&call, Statement(func(r rune) (ret State) {
			if x, e := call.GetExpression(); e != nil {
				p.err = errutil.Append(p.err, e)
			} else if len(x) > 0 {
				p.out = x // longest match wins
			}
			return
		})),
		MakeChain(&seq, Statement(func(r rune) (ret State) {
			if x, e := seq.GetExpression(); e != nil {
				p.err = errutil.Append(p.err, e)
			} else if len(x) > 0 {
				p.out = x // longest match wins
			}
			return
		})),
	)
	return par.NewRune(r)
}
