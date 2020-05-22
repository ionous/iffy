package chart

import (
	"sort"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/template/types"
)

// OperatorParser reads a binary operator.
type OperatorParser struct {
	next []rune // accumulated runes
	ofs  int    // search offset with our operator list
}

func (p *OperatorParser) StateName() string {
	return "operators"
}

// GetOperator representing the result of parsing, can be -1 if nothing was matched.
func (p *OperatorParser) GetOperator() (ret types.Operator, err error) {
	if len(p.next) == 0 {
		ret = types.Operator(-1)
	} else {
		next := string(p.next)
		if match := list[p.ofs]; next == match.Text {
			ret = match.Op
		} else {
			err = errutil.New("not an operator; provisionally matched", next)
		}
	}
	return
}

// NewRune starts on the first character of the operator; unless we find an exact match, we are done with the state.
func (p *OperatorParser) NewRune(r rune) (ret State) {
	// grow string we're trying to match
	next := string(append(p.next, r))
	// find where the in the list of operators this match would live
	i := p.ofs + sort.Search(len(list)-p.ofs, func(i int) bool {
		return list[i+p.ofs].Text >= next
	})
	// a negative result from search is possible:
	// if so, it would shift the offset backwards.
	if i >= p.ofs && i < len(list) {
		el := list[i].Text
		// we *might* match el in full or in part.
		if cnt := len(next); cnt <= len(el) && el[:cnt] == next {
			// record our partial match
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

// must stay sorted. see TestOpListIsSorted.
var list = []_Match{
	{types.NEQ, "!="},
	{types.REM, "%"},
	{types.MUL, "*"},
	{types.ADD, "+"},
	{types.SUB, "-"},
	{types.QUO, "/"},
	{types.LSS, "<"},
	{types.LEQ, "<="},
	{types.EQL, "="},
	{types.GTR, ">"},
	{types.GEQ, ">="},
	{types.LAND, "and"},
	{types.LOR, "or"}, // if this were || we'd have to make special provisions in the expression parser to handle the difference between a pipe (|) and an or (||)
}
