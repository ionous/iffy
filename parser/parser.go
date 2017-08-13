package parser

import (
	"github.com/ionous/errutil"
)

type Scanner interface {
	// search the words in cursor, looking for the best results.
	// note: by design, cursor may be out of range when scan is called.
	Scan(Scope, Cursor) (Result, error)
}

func Parse(scope Scope, match Scanner, in []string) (ret *ResultList, err error) {
	if r, e := match.Scan(scope, Cursor{Words: in}); e != nil {
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

func (cs Cursor) CurrentWord() (ret string, okay bool) {
	if cs.Pos < len(cs.Words) {
		ret, okay = cs.Words[cs.Pos], true
	}
	return
}
