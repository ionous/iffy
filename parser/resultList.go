package parser

import (
	"github.com/ionous/iffy/ident"
)

// ResultList contains multiple results. Its methods help tease out its contents.
type ResultList struct {
	list  []Result
	count int
}

// WordsMatched returns the number of words matched.
func (rs *ResultList) WordsMatched() int {
	return rs.count
}

func (rs *ResultList) Results() []Result {
	return rs.list
}

// AddResult to the list, updating the number of words matched.
func (rs *ResultList) AddResult(r Result) {
	if rl, ok := r.(*ResultList); ok {
		rs.list = append(rs.list, rl.list...)
		rs.count += rl.count
	} else {
		rs.list = append(rs.list, r)
		rs.count += r.WordsMatched()
	}
	return
}

// Last result in the list, true if the list was not empty. Generally, when the parser succeeds, this is an Action.
func (rs *ResultList) Last() (ret Result, okay bool) {
	if cnt := len(rs.list); cnt > 0 {
		ret, okay = rs.list[cnt-1], true
	}
	return
}

// Objects used by this result. Idenfified via Noun.GetId()
func (rs *ResultList) Objects() (ret []ident.Id) {
	for _, r := range rs.list {
		switch k := r.(type) {
		case ResolvedObject:
			n := k.NounVisitor
			ret = append(ret, n.GetId())
		case ResolvedMulti:
			for _, n := range k.Nouns {
				ret = append(ret, n.GetId())
			}
		}
	}
	return
}
