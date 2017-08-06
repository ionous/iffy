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
		cs := Cursor{at + ctx.Word, ctx.Words}
		if r := ctx.Match.Scan(ctx, &cs); r > 0 {
			ctx.Match, ret = ctx.Match.GetNext(), true
		}
	}
	return
}

type Cursor struct {
	ofs   int
	words []string
}

func (cs *Cursor) NextWord() (ret string, okay bool) {
	if cs.ofs < len(cs.words) {
		ret, okay = cs.words[cs.ofs], true
		cs.ofs++
	}
	return
}
