package parser

import (
	"strings"

	"github.com/ionous/errutil"
)

// // Series matches any set of words.
// type Series struct {
// 	Next Rule
// }
// // Number matches any series of digits
// // FUTURE: commas? and words.
// type Number struct {
// }

// // Reverse queues the next match till later.
// type Reverse struct {
// }

// // Actor handles addressing actors by name. For example: "name, "
// type Actor struct {
// }

// Word matches one word.
type Word struct {
	Word string
}

func (w *Word) Scan(ctx Context, bounds Bounds, cs Cursor) (ret Result, err error) {
	if word := cs.CurrentWord(); len(word) == 0 {
		err = Underflow{Depth(cs.Pos)}
	} else if !strings.EqualFold(word, w.Word) {
		err = MismatchedWord{w.Word, word, Depth(cs.Pos)}
	} else {
		ret = ResolvedWord{word}
	}
	return
}

// AnyOf matches any one of the passed scanners; whichever first matches.
type AnyOf struct {
	Match []Scanner
}

// Scan implements Scanner.
func (m *AnyOf) Scan(ctx Context, bounds Bounds, cs Cursor) (ret Result, err error) {
	if cnt := len(m.Match); cnt == 0 {
		err = errutil.New("no rules specified for any of")
	} else {
		i, errorDepth := 0, -1 // keep the most informative error
		for ; i < cnt; i++ {
			if res, e := m.Match[i].Scan(ctx, bounds, cs); e == nil {
				ret, err = res, nil
				break
			} else if d := DepthOf(e); d > errorDepth {
				err, errorDepth = e, d
				// keep looking for success
			}
		}
	}
	return
}

// AllOf matches the passed matchers in order.
type AllOf struct {
	Match []Scanner
}

func (m *AllOf) Scan(ctx Context, bounds Bounds, cs Cursor) (Result, error) {
	return m.scan(ctx, bounds, cs)
}

func (m *AllOf) scan(ctx Context, bounds Bounds, cs Cursor) (ret *ResultList, err error) {
	var rl ResultList
	if cnt := len(m.Match); cnt == 0 {
		err = errutil.New("no rules specified for all of")
	} else {
		var i int
		for ; i < cnt; i++ {
			if res, e := m.Match[i].Scan(ctx, bounds, cs); e != nil {
				err = e
				break
			} else {
				rl.AddResult(res)
				cs = cs.Skip(res.WordsMatched())
			}
		}
		if i == cnt {
			ret = &rl
		}
	}
	return
}
