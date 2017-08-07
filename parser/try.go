package parser

import (
	"strings"
)

// Word matches one word.
type Word struct {
	Word string
}

func (w *Word) Scan(scan Cursor, ctx *Context) (ret int, okay bool) {
	if n, ok := scan.GetNext(); ok {
		if strings.EqualFold(n, w.Word) {
			ret, okay = 1, true
		}
	}
	return
}

// AnyOf matches any one of the passed scanners; whichever first matches.
type AnyOf struct {
	Match []Scanner
}

func (m *AnyOf) Scan(scan Cursor, ctx *Context) (ret int, okay bool) {
	for _, s := range m.Match {
		if advance, ok := s.Scan(scan, ctx); !ok {
			ctx.Results = Results{} // reset results on failure.
		} else {
			ret, okay = advance, true
			break
		}
	}
	return
}

// AllOf matches the passed matchers in order.
type AllOf struct {
	Match []Scanner
}

func (m *AllOf) Scan(scan Cursor, ctx *Context) (ret int, okay bool) {
	if i, cnt, ofs := 0, len(m.Match), 0; cnt > 0 {
		for ; i < cnt; i++ {
			s := m.Match[i]
			if advance, ok := s.Scan(scan.Step(ofs), ctx); !ok {
				break
			} else {
				ofs += advance
			}
		}
		if i == cnt {
			ret, okay = ofs, true
		}
	}
	return
}
