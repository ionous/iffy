package chart

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/template/postfix"
)

// ExpParser reads either a call or a sequence.
type ExpParser struct {
	argFactory ExpressionStateFactory
	out        postfix.Expression
	err        error
}

func MakeExpParser(f ExpressionStateFactory) ExpParser {
	return ExpParser{argFactory: f}
}

func (p ExpParser) GetExpression() (postfix.Expression, error) {
	return p.out, p.err
}

// NewRune starts with the first character past the bar
func (p *ExpParser) NewRune(r rune) State {
	call := MakeCallParser(p.argFactory)
	seq := SequenceParser{}
	par := parallel(
		makeChain(&call, Statement(func(r rune) (ret State) {
			if x, e := call.GetExpression(); e != nil {
				p.err = errutil.Append(p.err, e)
			} else if len(x) > 0 {
				p.out = x // longest match wins
			}
			return
		})),
		makeChain(&seq, Statement(func(r rune) (ret State) {
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
