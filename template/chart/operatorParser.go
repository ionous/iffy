package chart

import (
	"github.com/ionous/iffy/template/types"
	"sort"
)

// OperatorParser reads a binary operator.
type OperatorParser struct {
	next []rune
	ofs  int
}

// GetOperator representing the result of parsing.
func (p OperatorParser) GetOperator() (ret types.Operator, okay bool) {
	if len(p.next) > 0 {
		ret, okay = list[p.ofs].Op, true
	}
	return
}

// NewRune starts on the first character of the operator; unless we find an exact match, we are done with the state.
func (p *OperatorParser) NewRune(r rune) (ret State) {
	next := string(append(p.next, r))
	i := p.ofs + sort.Search(len(list)-p.ofs, func(i int) bool {
		return list[i+p.ofs].Text >= next
	})
	if i >= p.ofs && i < len(list) {
		el := list[i].Text
		if cnt := len(next); cnt <= len(el) && el[:cnt] == next {
			p.ofs = i
			p.next = []rune(next)
			ret = p
		}
	}
	return
}

type _Match struct {
	Op   types.Operator
	Text string
}

var list = []_Match{
	{types.REM, "%"},
	{types.MUL, "*"},
	{types.ADD, "+"},
	{types.SUB, "-"},
	{types.QUO, "/"},
	{types.LSS, "<"},
	{types.LEQ, "<="},
	{types.NEQ, "<>"},
	{types.EQL, "="},
	{types.GTR, ">"},
	{types.GEQ, ">="},
	{types.LAND, "and"},
	{types.LOR, "or"}, // if this was || we'd have to make special provisions in the expression parser to handle the difference between a pipe (|) and an or (||)
}
