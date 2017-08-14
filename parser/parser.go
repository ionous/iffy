package parser

import (
	"github.com/ionous/errutil"
)

// Scanner searches words looking for good results.
// ( perhaps its truly a tokenzer and the results, tokens )
type Scanner interface {
	// Scan for results.
	// note: by design, cursor may be out of range when scan is called.
	Scan(Context, Scope, Cursor) (Result, error)
}

func Parse(ctx Context, match Scanner, in []string) (ret *ResultList, err error) {
	if r, e := match.Scan(ctx, ctx.GetPlayerScope(""), Cursor{Words: in}); e != nil {
		err = e
	} else if rs, ok := r.(*ResultList); !ok {
		err = errutil.Fmt("expected result list, got %T", r)
	} else {
		ret = rs
	}
	return
}

type Cursor struct {
	Pos   int
	Words []string
}

func (cs Cursor) Skip(i int) Cursor {
	return Cursor{cs.Pos + i, cs.Words}
}

func (cs Cursor) CurrentWord() (ret string) {
	if cs.Pos < len(cs.Words) {
		ret = cs.Words[cs.Pos]
	}
	return
}
