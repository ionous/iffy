package chart

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/template/postfix"
)

// SubdirParser reads a single operand or sub bracket pipeline.
// "a" -or- "{a!|...}"
type SubdirParser struct {
	exp postfix.Expression
	err error
}

func (p *SubdirParser) StateName() string {
	return "subdir"
}

func (p *SubdirParser) NewRune(r rune) (ret State) {
	switch {
	case isOpenBracket(r):
		var pipe PipeParser
		ret = MakeChain(&pipe, Statement("after pipe", func(r rune) (ret State) {
			if !isCloseBracket(r) {
				p.err = errutil.New("unclosed inner directive")
			} else {
				p.exp, p.err = pipe.GetExpression()
				ret = Terminal // eat the closing bracket.
			}
			return
		}))

	default:
		var op OperandParser
		ret = ParseChain(r, &op, StateExit("subdir", func() {
			p.exp, p.err = op.GetExpression()
		}))
	}
	return
}

func (p *SubdirParser) GetExpression() (postfix.Expression, error) {
	return p.exp, p.err
}
