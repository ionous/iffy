package parser

import (
	"github.com/ionous/errutil"
)

func Parse(ctx Context, match Scanner, in []string) (ret *ResultList, err error) {
	if scope, e := ctx.GetPlayerScope(""); e != nil {
		err = e
	} else if r, e := match.Scan(ctx, scope, Cursor{Words: in}); e != nil {
		err = e
	} else if rs, ok := r.(*ResultList); !ok {
		err = errutil.Fmt("expected result list, got %T", r)
	} else {
		ret = rs
	}
	return
}

// Scanner searches words looking for good results.
// ( perhaps its truly a tokenzer and the results, tokens )
type Scanner interface {
	// Scan for results.
	// note: by design, cursor may be out of range when scan is called.
	Scan(Context, Scope, Cursor) (Result, error)
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
