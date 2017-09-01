package parser

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
