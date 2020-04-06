package scope

import (
	"testing"

	"github.com/ionous/iffy/assign"
)

func TestLoop(t *testing.T) {
	caps := [][]struct{ first, last bool }{
		{{true, true}},
		{{true, false}, {false, true}},
		{{true, false}, {false, false}, {false, true}},
	}
	for c, cap := range caps {
		var lf LoopFactory
		for i, cnt := 0, len(cap); i < cnt; i++ {
			cap := cap[i]
			count := i + 1
			atEnd := count == cnt
			s := lf.NextScope(EmptyScope{}, !atEnd)

			if p, e := s.GetVariable("index"); e != nil {
				t.Fatal("loop", i, e)
			} else if fidx, e := assign.ToFloat(p); e != nil || fidx != float64(count) {
				t.Fatal("index error", fidx, "at", c, i, e)
			} else if fidx != float64(count) {
				t.Fatal("loop", i, fidx, "!=", count)
			} else if p, e := s.GetVariable("first"); e != nil {
				t.Fatal(e)
			} else if first, e := assign.ToBool(p); e != nil || first != cap.first {
				t.Fatal("first error", first, "at", c, i, e)
			} else if p, e := s.GetVariable("last"); e != nil {
				t.Fatal(e)
			} else if last, e := assign.ToBool(p); e != nil || last != cap.last {
				t.Fatal("last error", last, "at", c, i, e)
			} else {
				_, e := s.GetVariable("nothing")
				if _, ok := e.(UnknownVariable); !ok {
					t.Fatal("expected loop error")
				} else {
					t.Log("loop", i, "of", cnt, fidx, first, last)
				}
			}
		}
	}
}
