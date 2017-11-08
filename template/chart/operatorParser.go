package chart

import (
	"github.com/ionous/iffy/template"
	"sort"
)

// OperatorParser reads a binary operator.
type OperatorParser struct {
	next []rune
	ofs  int
}

// GetOperator representing the result of parsing.
func (p OperatorParser) GetOperator() (ret template.Operator, okay bool) {
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
	Op   template.Operator
	Text string
}

var list = []_Match{
	{template.REM, "%"},
	{template.MUL, "*"},
	{template.ADD, "+"},
	{template.SUB, "-"},
	{template.QUO, "/"},
	{template.LSS, "<"},
	{template.LEQ, "<="},
	{template.NEQ, "<>"},
	{template.EQL, "="},
	{template.GTR, ">"},
	{template.GEQ, ">="},
	{template.LAND, "and"},
	{template.LOR, "or"}, // if this was || we'd have to make special provisions in the expression parser to handle the difference between a pipe (|) and an or (||)
}
