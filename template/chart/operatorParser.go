package chart

import "sort"

type Match struct {
	Op   Operator
	Text string
}

var list = []Match{
	{REM, "%"},
	{MUL, "*"},
	{ADD, "+"},
	{SUB, "-"},
	{QUO, "/"},
	{LSS, "<"},
	{LEQ, "<="},
	{NEQ, "<>"},
	{EQL, "="},
	{GTR, ">"},
	{GEQ, ">="},
	{LAND, "and"},
	{LOR, "or"}, // if this was || we'd have to make special provisions in the expression parser to handle the difference between a pipe (|) and an or (||)
}

type OperatorParser struct {
	next []rune
	ofs  int
}

func (p OperatorParser) GetOperator() (ret Operator, okay bool) {
	if len(p.next) > 0 {
		ret, okay = list[p.ofs].Op, true
	}
	return
}

// unless we find an exact match, we are done with the state.
func (p *OperatorParser) NewRune(r rune) (ret State) {
	next := string(append(p.next, r))
	i := sort.Search(len(list)-p.ofs, func(i int) bool {
		return list[i].Text >= next
	})
	if i >= 0 && i < len(list) {
		el := list[i].Text
		if cnt := len(next); cnt <= len(el) && el[:cnt] == next {
			p.ofs = i
			p.next = []rune(next)
			ret = p
		}
	}
	return
}
