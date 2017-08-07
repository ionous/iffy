package parser

type Context struct {
	Scope
	// Soft bool // when Soft is enabled, we can try many matches
	// i almost wonder, could Span just simply take over with a new context implementation?
	// return its own Next's Next token transparently
	Results
}

type Results struct {
	Actor  string
	Action string
	// Prep   []Action -- take because prefererably held
	Matches   []Ranking
	NeedsNoun bool
}

func (r *Results) Complete() (ret bool) {
	// there will be more --
	if !r.NeedsNoun {
		i, cnt := 0, len(r.Matches)
		for ; i < cnt; i++ {
			if m := r.Matches[i]; len(m.Nouns) > 1 {
				break // we are not complete if we matched more than one noun
			}
		}
		ret = i == cnt
	}
	return
}

type Ranking struct {
	Rank  int
	Nouns []string
}

// in should be split on fields -- currently just space.
func Parse(scope Scope, match Scanner, in []string) (ret Results, okay bool) {
	ctx := Context{
		Scope: scope,
	}
	pos := Cursor{
		Words: in,
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
