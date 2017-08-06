package parser

type Context struct {
	Scope
	Match Matcher
	Words []string
	Word  int
	Soft  bool // when Soft is enabled, we can try many matches
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

func (ctx *Context) Advance() (ret bool) {
	for try, at := true, 0; try && at < len(ctx.Words); try, at = ctx.Soft, at+1 {
		cs := Cursor(at + ctx.Word)
		if advance, ok := ctx.Match.Scan(cs, ctx); ok {
			ctx.Word += advance
			ctx.Match = ctx.Match.GetNext()
			ret = true
		}
	}
	return
}

type Cursor int

func (cs Cursor) Next(i int) Cursor {
	return cs + Cursor(i)
}

func (cs Cursor) NextWord(ctx *Context) (ret string, okay bool) {
	if int(cs) < len(ctx.Words) {
		ret, okay = ctx.Words[cs], true
	}
	return
}
