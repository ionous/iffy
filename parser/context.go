package parser

import (
	"strings"
)

type Context struct {
	Scope
	Soft bool // when Soft is enabled, we can try many matches
	// i almost wonder, could Span just simply take over with a new context implementation?
	// return its own Next's Next token transparently
	Results
}

type Results struct {
	Actor  string
	Action string
	// Prep   []Action -- take because prefererably held
	Matches []Ranking
}

type Ranking struct {
	Rank  int
	Nouns []string
}

func Parse(scope Scope, match Scanner, in string) (ret Results, okay bool) {
	ctx := Context{
		Scope: scope,
	}
	pos := Cursor{
		Words: strings.Fields(in),
	}
	if r, ok := match.Scan(pos, &ctx); ok {
		if len(pos.Words) == r {
			ret, okay = ctx.Results, true
		}
	}
	return
}

// func (ctx *Context) Advance(cs Cursor, m Matcher) (step Cursor, ret Matcher, okay bool) {
// 	// note: when we are soft-scanning, we dont care if we dont match unless we are out of words.
// 	for try := true; try && cs.HasNext(); try, cs = ctx.Soft, cs.Step(1) {
// 		if advance, ok := m.Scan(cs, ctx); ok {
// 			ctx.Soft = false
// 			ret, step, okay = m.GetNext(), cs.Step(advance), true
// 			break
// 		}
// 	}
// 	return
// }

type Cursor struct {
	Word  int
	Words []string
}

func (cs Cursor) Step(i int) Cursor {
	return Cursor{cs.Word + i, cs.Words}
}

func (cs Cursor) HasNext() bool {
	return cs.Word < len(cs.Words)
}

func (cs Cursor) GetNext() (ret string, okay bool) {
	if cs.Word < len(cs.Words) {
		ret, okay = cs.Words[cs.Word], true
	}
	return
}
